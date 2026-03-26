package scheduler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/logic/scheduler"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/svc"
	"github.com/Humphrey-He/star-flow-scheduler/apps/scheduler/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetJobInstanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetJobInstanceRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := scheduler.NewGetJobInstanceLogic(r.Context(), svcCtx)
		resp, err := l.GetJobInstance(&req)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusNotFound, map[string]any{
					"error": map[string]any{
						"code":    "not_found",
						"message": "instance not found",
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
