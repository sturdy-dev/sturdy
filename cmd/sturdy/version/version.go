package version

import "fmt"

var Version = "v0.0.0-development"

func VersionCMD() {
	fmt.Printf("Sturdy %s\n", Version)
}
