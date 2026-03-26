package scheduler

import (
	"errors"
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/logic/scheduler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/repo"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateJobHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateJobRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := scheduler.NewCreateJobLogic(r.Context(), svcCtx)
		resp, err := l.CreateJob(&req)
		if err != nil {
			if errors.Is(err, repo.ErrAlreadyExists) {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusConflict, map[string]any{
					"error": map[string]any{
						"code":    "already_exists",
						"message": "job_code already exists",
					},
				})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
