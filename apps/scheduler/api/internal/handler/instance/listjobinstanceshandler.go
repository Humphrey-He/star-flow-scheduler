package instance

import (
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/logic/instance"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListJobInstancesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListJobInstancesRequest
		if err := httpx.Parse(r, &req); err != nil {
			errx.WriteError(w, r, errx.InvalidParam(err.Error()))
			return
		}

		l := instance.NewListJobInstancesLogic(r.Context(), svcCtx)
		resp, err := l.ListJobInstances(&req)
		if err != nil {
			errx.WriteError(w, r, err)
		} else {
			errx.WriteOK(w, r, resp)
		}
	}
}
