package controller

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"
	"github.com/corvofeng/go-flarum/util"
	"github.com/op/go-logging"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewS3Client(cfg model.S3ConfigConf) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
}

func S3PushObject(logger *logging.Logger, user model.User, cfg model.S3ConfigConf, fileHeader *multipart.FileHeader) (model.UserFiles, error) {
	client, err := NewS3Client(cfg)
	uf := model.UserFiles{}
	if err != nil {
		return uf, fmt.Errorf("failed to create S3 client: %w", err)
	}
	uuid := util.GetUUID()
	now := time.Now()
	dateString := now.Format("2006/01/02") // Formats as YYYY/MM/DD

	// For example:
	// https://rawforcorvofeng.cn/flarum/2025/06/20/f6d9dbaa-34f0-455f-8f8f-ab2682907af5-image.png
	objectName := path.Join(
		cfg.Path,
		dateString,
		fmt.Sprintf("%s-%s", uuid, fileHeader.Filename),
	)
	file, err := fileHeader.Open()
	if err != nil {
		return uf, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	uploadInfo, err := client.PutObject(
		context.Background(),
		cfg.Bucket,
		objectName,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return uf, fmt.Errorf("failed to upload file to S3: %w", err)
	}
	filePath, err := url.JoinPath(cfg.S3BaseURL, objectName)
	if err != nil {
		return uf, fmt.Errorf("failed to construct file URL: %w", err)
	}

	logger.Debugf("File uploaded successfully: %s, %+v", filePath, uploadInfo)

	userfile := model.UserFiles{
		UserID:     user.ID,
		UUID:       uuid,
		FileName:   fileHeader.Filename,
		FilePath:   filePath,
		FileSize:   fileHeader.Size,
		FileType:   fileHeader.Header.Get("Content-Type"),
		Visibility: "public", // 可根据需要设置为 "private"
	}
	return userfile, nil
}

func LocalUploadObject(logger *logging.Logger, user model.User, uploadDir string, fileHeader *multipart.FileHeader) (model.UserFiles, error) {
	file, err := fileHeader.Open()
	uf := model.UserFiles{}
	if err != nil {
		return uf, fmt.Errorf("file open error: %w", err)
	}
	defer file.Close()

	// 检查文件类型（可选）
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return uf, fmt.Errorf("invalid file extension")
	}
	uuid := util.GetUUID()
	now := time.Now()
	dateString := now.Format("2006/01/02") // Formats as YYYY/MM/DD

	// 生成保存路径
	objectName := path.Join(
		uploadDir,
		dateString,
		fmt.Sprintf("%s-%s", uuid, fileHeader.Filename),
	)
	// 确保目录存在
	err = os.MkdirAll(filepath.Dir(objectName), os.ModePerm)
	if err != nil {
		return uf, fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(objectName)
	if err != nil {
		return uf, fmt.Errorf("file create error: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return uf, fmt.Errorf("file copy error: %w", err)
	}

	logger.Debug("Local uploaded file:", fileHeader.Filename)
	userfile := model.UserFiles{
		UserID:     user.ID,
		UUID:       uuid,
		FileName:   fileHeader.Filename,
		FilePath:   "/" + objectName,
		FileSize:   fileHeader.Size,
		FileType:   fileHeader.Header.Get("Content-Type"),
		Visibility: "public", // 可根据需要设置为 "private"
	}

	return userfile, nil
}

// 上传接口
func FlarumUpload(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	logger := ctx.GetLogger()

	// 解析表单
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Parse form error: "+err.Error()))
		return
	}

	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		h.flarumErrorJsonify(w, createSimpleFlarumError("No file uploaded"))
		return
	}

	uploadDir := h.App.Cf.Main.UploadDir
	s3cfg := h.App.Cf.Main.S3Config
	currentUser := ctx.currentUser
	maxUploadSize := h.App.Cf.Site.UploadMaxSizeByte
	cnt, err := model.SQLGetUserDailyUploads(h.App.GormDB, currentUser.ID)
	logger.Debugf("User %d has uploaded %d files today", currentUser.ID, cnt)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Database error: "+err.Error()))
		return
	}
	if cnt >= maxUploadSize {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Daily upload limit reached"))
		return
	}

	var uploadedFiles []model.UserFiles
	for _, fileHeader := range files {
		var uf model.UserFiles

		if s3cfg.Endpoint == "" {
			uf, err = LocalUploadObject(logger, *currentUser, uploadDir, fileHeader)
		} else {
			// 如果没有配置 S3，则使用本地上传
			uf, err = S3PushObject(logger, *currentUser, s3cfg, fileHeader)
		}
		if err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Upload error: "+err.Error()))
			return
		}
		logger.Debugf("File uploaded successfully: %s", uf.FilePath)
		err = model.SQLSaveUserFile(h.App.GormDB, &uf)
		if err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Database save error: "+err.Error()))
			return
		}

		uploadedFiles = append(uploadedFiles, uf)
	}

	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针

	var res []flarum.Resource
	for _, uf := range uploadedFiles {
		fofUF := model.FlarumCreateFoFUploadFiles(uf, *ctx.currentUser)
		res = append(res, fofUF)
	}
	apiDoc.SetData(res)
	h.jsonify(w, coreData.APIDocument)
}

// http://127.0.0.1:8082/api/v1/flarum/fof/uploads?filter%5Buser%5D=2&page%5Boffset%5D=0
// FlarumUploads 获取用户上传的文件
func FlarumUploadsAll(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	logger := ctx.GetLogger()
	userfiles, err := model.SQLGetUserFiles(h.App.GormDB, ctx.currentUser.ID)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Database error: "+err.Error()))
		return
	}
	logger.Debugf("User %d has uploaded %d files", ctx.currentUser.ID, len(userfiles))
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument // 注意, 获取到的是指针
	var res []flarum.Resource

	for _, uf := range userfiles {
		fofUF := model.FlarumCreateFoFUploadFiles(uf, *ctx.currentUser)
		res = append(res, fofUF)
	}
	apiDoc.SetData(res)
	h.jsonify(w, coreData.APIDocument)
}
