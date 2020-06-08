package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

const (
	tokenLength = 32
)

func SliceUniqInt(s []int) []int {
	if len(s) == 0 {
		return s
	}
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func SliceUniqStr(s []string) []string {
	if len(s) == 0 {
		return s
	}
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func HashFileMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func CheckTags(tags string) string {
	limit := 5
	seen := map[string]struct{}{}
	tmpTags := strings.Replace(tags, " ", ",", -1)
	tmpTags = strings.Replace(tmpTags, "，", ",", -1)
	tagList := make([]string, limit)
	j := 0
	for _, tag := range strings.Split(tmpTags, ",") {
		if len(tag) > 1 && len(tag) < 25 {
			if _, ok := seen[tag]; ok {
				continue
			}
			seen[tag] = struct{}{}
			tagList[j] = tag
			j++
			if j == limit {
				break
			}
		}
	}
	return strings.Join(tagList[:j], ",")
}

// A token is generated by returning tokenLength bytes
// from crypto/rand
func generateToken() []byte {
	bytes := make([]byte, tokenLength)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic(err)
	}

	return bytes
}

func b64encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func b64decode(data string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return decoded
}

// GetNewToken 获取新的csrf token
func GetNewToken() string {
	// golang会将参数中的+替换为空格, 这里生成token时就直接替换
	// 	http://weakyon.com/2017/05/04/something-of-golang-url-encoding.html
	return strings.Replace(b64encode(generateToken()), "+", "", -1)
}

// VerifyToken verifies the sent token equals the real one
// and returns a bool value indicating if tokens are equal.
// Supports masked tokens. realToken comes from Token(r) and
// sentToken is token sent unusual way.
func VerifyToken(realToken, sentToken string) bool {
	r := b64decode(realToken)
	s := b64decode(sentToken)
	return subtle.ConstantTimeCompare(r, s) == 1
}
