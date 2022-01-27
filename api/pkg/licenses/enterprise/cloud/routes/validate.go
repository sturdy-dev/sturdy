package routes

import (
	"net/http"

	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
	service_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Validate(
	logger *zap.Logger,
	serviceLicenses *service_licenses.Service,
	serviceValidations *service_validations.Service,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		licenseKey := c.Param("key")

		license, err := serviceLicenses.ValidateByKey(c.Request.Context(), licenseKey)
		if err != nil {
			logger.Error("failed to get license", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		if _, err := serviceValidations.Create(c.Request.Context(), license.ID, license.Status); err != nil {
			logger.Error("failed to create validation", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, license)
	}
}
