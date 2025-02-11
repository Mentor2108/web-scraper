package rest

import (
	"backend-service/defn"
	"backend-service/service"
	"backend-service/util"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func GetFileById(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	ctx := r.Context()
	log := util.GetGlobalLogger(ctx)

	fileId := param.ByName("id")
	if strings.EqualFold(fileId, "") {
		cerr := util.NewCustomError(ctx, "empty-file-id", errors.New("no file id provided"))
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	fileContent, fileInfo, cerr := service.GetFile(ctx, fileId)
	if cerr != nil {
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	var contentType string
	switch fileInfo.FileType {
	case ".html":
		contentType = defn.ContentTypeHTMLText
	case ".txt":
		contentType = defn.ContentTypePlainText
	case ".md":
		contentType = defn.ContentTypeMarkdownText
	default:
		contentType = defn.ContentTypeOctetStream
	}
	w.Header().Add(defn.HTTPHeaderContentType, contentType)
	w.Write(fileContent)
}
