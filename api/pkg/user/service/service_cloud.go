//go:build cloud
// +build cloud

package service

import (
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/emails/transactional"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_onetime "getsturdy.com/api/pkg/onetime/service"
	db_user "getsturdy.com/api/pkg/user/db"
)

type Service struct {
	*commonService
}

func New(
	logger *zap.Logger,
	userRepo db_user.Repository,
	jwtService *service_jwt.Service,
	onetimeService *service_onetime.Service,
	transactionalEmailSender transactional.EmailSender,
	analyticsClient analytics.Client,
) *Service {
	return &Service{
		commonService: &commonService{
			logger:                   logger,
			userRepo:                 userRepo,
			jwtService:               jwtService,
			onetimeService:           onetimeService,
			transactionalEmailSender: transactionalEmailSender,
			analyticsClient:          analyticsClient,
		},
	}
}
