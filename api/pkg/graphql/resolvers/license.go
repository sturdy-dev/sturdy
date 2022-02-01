package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type LicenseRootResolver interface {
	// Internal
	InternalByKey(ctx context.Context, key string) (LicenseResolver, error)
	InternalListForOrganizationID(ctx context.Context, id string) ([]LicenseResolver, error)
}

type LicenseResolver interface {
	ID() graphql.ID
	Key() string
	ExpiresAt() (int32, error)
	CreatedAt() (int32, error)
	Seats() int32

	Status(context.Context) (LicenseStatus, error)
	Messages(context.Context) ([]LicenseMessageResolver, error)
}

type LicenseStatus string

const (
	LicenseStatusUnknown LicenseStatus = ""
	LicenseStatusValid   LicenseStatus = "Valid"
	LicenseStatusInvalid LicenseStatus = "Invalid"
)

type LicenseMessageType string

const (
	LicenseMessageTypeUnknown      LicenseMessageType = ""
	LicenseMessageTypeNotification LicenseMessageType = "Notification"
	LicenseMessageTypeBanner       LicenseMessageType = "Banner"
	LicenseMessageTypeFullscreen   LicenseMessageType = "Fullscreen"
)

type LicenseMessageLevel string

const (
	LicenseMessageLevelUnknown LicenseMessageLevel = ""
	LicenseMessageLevelInfo    LicenseMessageLevel = "Info"
	LicenseMessageLevelWarning LicenseMessageLevel = "Warning"
	LicenseMessageLevelError   LicenseMessageLevel = "Error"
)

type LicenseMessageResolver interface {
	Text() string
	Type() LicenseMessageType
	Level() LicenseMessageLevel
}
