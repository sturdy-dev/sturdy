package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"getsturdy.com/api/pkg/installations/statistics"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	"getsturdy.com/api/pkg/ip"

	"go.uber.org/zap"
)

func Create(logger *zap.Logger, service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		statistic := &statistics.Statistic{}
		if err := json.NewDecoder(r.Body).Decode(statistic); err != nil {
			http.Error(w, "invalid format", http.StatusBadRequest)
			return
		}

		if ip, found := ip.FromContext(r.Context()); found {
			statistic.IP = ip
		}
		statistic.ReceivedAt = time.Now()

		if err := service.Accept(r.Context(), statistic); err != nil {
			logger.Error("failed to accept statistic", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
