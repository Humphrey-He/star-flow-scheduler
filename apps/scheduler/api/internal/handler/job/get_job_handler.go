package job

import (
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/errx"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/logic/job"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetJobHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetJobRequest
		if err := httpx.Parse(r, &req); err != nil {
			errx.WriteError(w, r, errx.InvalidParam(err.Error()))
			return
		}

		l := job.NewGetJobLogic(r.Context(), svcCtx)
		resp, err := l.GetJob(&req)
		if err != nil {
			errx.WriteError(w, r, err)
		} else {
			errx.WriteOK(w, r, resp)
		}
	}
}
