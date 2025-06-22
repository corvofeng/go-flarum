package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/corvofeng/go-flarum/model"
	"github.com/corvofeng/go-flarum/model/flarum"
)

func FlarumBlogMeta(w http.ResponseWriter, r *http.Request) {
	ctx := GetRetContext(r)
	h := ctx.h
	logger := ctx.GetLogger()
	_mid, _ := h.safeGetParm(r, "mid")
	// if _mid == "" {
	// 	h.flarumErrorJsonify(w, createSimpleFlarumError("mid is empty"))
	// 	return
	// }

	/*
		{
		  "data": {
		    "type": "blogMeta",
		    "attributes": {
		      "summary": "topic summary",
		      "featuredImage": "https://rawforcorvofeng.cn/blog/2025/06/15/1749955395119.png",
		      "isFeatured": false,
		      "isSized": false,
		      "isPendingReview": false,
		      "relationships": null
		    },
		    "id": "1"
		  }
		}
	*/
	type BlogMetaData struct {
		Data struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Summary         string `json:"summary"`
				FeaturedImage   string `json:"featuredImage"`
				IsSized         bool   `json:"IsSized"`
				IsFeatured      bool   `json:"isFeatured"`
				IsPendingReview bool   `json:"isPendingReview"`
			} `json:"attributes"`
			Relationships struct {
				Discussions struct {
					Data struct {
						Type string `json:"type"`
						ID   string `json:"id"`
					} `json:"data"`
				} `json:"discussion"`
			} `json:"relationships"`
		} `json:"data"`
	}

	bytedata, err := io.ReadAll(r.Body)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Read body error:"+err.Error()))
		return
	}
	diss := BlogMetaData{}
	err = json.Unmarshal(bytedata, &diss)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("json Decode err:"+err.Error()))
		return
	}
	logger.Debugf("Update blog meta with: %+v", diss)
	if _mid != "" && diss.Data.ID != "" && diss.Data.ID != _mid {
		h.flarumErrorJsonify(w, createSimpleFlarumError("mid not match: "+_mid+" != "+fmt.Sprintf("%d", diss.Data.ID)))
		return
	}
	bmID := uint64(0)
	if diss.Data.ID != "" {
		bmID, err = strconv.ParseUint(diss.Data.ID, 10, 64)
		if err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Parse blog meta ID error: "+err.Error()))
			return
		}
	}

	bm := model.BlogMeta{
		ID:              bmID,
		Summary:         diss.Data.Attributes.Summary,
		FeaturedImage:   diss.Data.Attributes.FeaturedImage,
		IsSized:         diss.Data.Attributes.IsSized,
		IsFeatured:      diss.Data.Attributes.IsFeatured,
		IsPendingReview: diss.Data.Attributes.IsPendingReview,
	}

	if diss.Data.ID == "" {
		tid, err := strconv.ParseUint(diss.Data.Relationships.Discussions.Data.ID, 10, 64)
		if err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Parse discussion ID error: "+err.Error()))
			return
		}
		bm.TopicID = tid
	} else {
		b, err := model.SQLGetBlogMetaByID(h.App.GormDB, bm.ID)
		if err != nil {
			h.flarumErrorJsonify(w, createSimpleFlarumError("Get blog meta by ID error: "+err.Error()))
			return
		}
		bm.TopicID = b.TopicID
	}

	err = bm.CreateOrUpdate(h.App.GormDB)
	if err != nil {
		h.flarumErrorJsonify(w, createSimpleFlarumError("Create or update blog meta error: "+err.Error()))
		return
	}
	coreData := flarum.NewCoreData()
	apiDoc := &coreData.APIDocument
	apiDoc.SetData(model.FlarumCreateBlogMeta(bm, ctx.currentUser))

	// h.jsonify(w, model.FlarumCreateBlogMeta(bm, ctx.currentUser))
	// coreData, err := createFlarumPostAPIDoc(ctx, h.App.GormDB, redisDB, *h.App.Cf, rf, scf.TimeZone)
	// if err != nil {
	// 	h.flarumErrorJsonify(w, createSimpleFlarumError("Get api doc error"+err.Error()))
	// 	return
	// }
	h.jsonify(w, coreData.APIDocument)

}
