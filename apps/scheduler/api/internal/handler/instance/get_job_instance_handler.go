package instance

import (
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/logic/instance"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetJobInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetJobInstanceRequest
		if err := httpx.Parse(r, &req); err != nil {
			errx.WriteError(w, r, errx.InvalidParam(err.Error()))
			return
		}

		l := instance.NewGetJobInstanceLogic(r.Context(), svcCtx)
		resp, err := l.GetJobInstance(&req)
		if err != nil {
			errx.WriteError(w, r, err)
		} else {
			errx.WriteOK(w, r, resp)
		}
	}
}
