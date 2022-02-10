package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"getsturdy.com/api/pkg/installations/statistics"
	"getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	"getsturdy.com/api/pkg/ip"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func Create(logger *zap.Logger, service *service.Service) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		statistic := &statistics.Statistic{}
		if err := json.NewDecoder(r.Body).Decode(statistic); err != nil {
			http.Error(w, "invalid format", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(statistic); err != nil {
			http.Error(w, fmt.Sprintf("invalid input: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if ip, found := ip.FromContext(r.Context()); found {
			ips := ip.String()
			statistic.IP = &ips
		}
		statistic.ReceivedAt = time.Now()

		if err := service.Accept(r.Context(), statistic); err != nil {
			logger.Error("failed to accept statistic", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
