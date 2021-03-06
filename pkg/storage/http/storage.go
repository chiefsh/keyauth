package http

import (
	"net/http"

	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/keyauth/pkg"
	"github.com/infraboard/keyauth/pkg/storage"
)

func (h *handler) UploadGEOIPDBFile(w http.ResponseWriter, r *http.Request) {
	rctx := context.GetContext(r)
	tk, err := pkg.GetTokenFromContext(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	req := storage.NewUploadFileRequestFromHTTP(r)
	req.BucketName = rctx.PS.ByName("name")
	req.WithToken(tk)

	err = h.service.UploadFile(req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, "ok")
	return
}
