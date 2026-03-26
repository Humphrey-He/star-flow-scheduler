package scheduler

import (
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/logic/scheduler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListJobInstancesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListJobInstancesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := scheduler.NewListJobInstancesLogic(r.Context(), svcCtx)
		resp, err := l.ListJobInstances(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
