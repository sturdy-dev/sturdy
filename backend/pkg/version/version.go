package version

import (
	"time"
)

var Version = "development"
var BootedAt time.Time

func init() {
	BootedAt = time.Now()
}
