package version

import (
	"time"
)

var (
	Version            = "development"
	BootedAt time.Time = time.Now()
)

func IsDevelopment() bool {
	return Version == "development"
}

type DistributionType uint

const (
	DistributionTypeUndefined DistributionType = iota
	DistributionTypeOSS
	DistributionTypeEnterprise
	DistributionTypeCloud
)

func (t DistributionType) String() string {
	switch t {
	case DistributionTypeOSS:
		return "oss"
	case DistributionTypeEnterprise:
		return "enterprise"
	case DistributionTypeCloud:
		return "cloud"
	default:
		return "undefined"
	}
}
