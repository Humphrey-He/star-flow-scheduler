package errx

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteOK(w http.ResponseWriter, r *http.Request, data interface{}) {
	httpx.WriteJsonCtx(r.Context(), w, http.StatusOK, Response{
		Code:    CodeOK,
		Message: "OK",
		Data:    data,
	})
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	be := FromError(err)
	status := http.StatusInternalServerError
	switch be.Code {
	case CodeInvalidParam:
		status = http.StatusBadRequest
	case CodeNotFound:
		status = http.StatusNotFound
	case CodeStatusConflict, CodeIdempotentConflict:
		status = http.StatusConflict
	}

	httpx.WriteJsonCtx(r.Context(), w, status, Response{
		Code:    be.Code,
		Message: be.Message,
	})
}
