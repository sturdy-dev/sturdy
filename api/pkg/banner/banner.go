package banner

import (
	_ "embed"

	"github.com/fatih/color"

	"getsturdy.com/api/pkg/version"
)

//go:embed sturdy.txt
var sturdyBanner string

func PrintBanner() {
	c := color.New(color.FgBlack, color.FgHiYellow)

	c.Println()
	c.Println()
	c.Println(sturdyBanner)
	c.Printf("Sturdy %s %s\n", version.Type.String(), version.Version)
	c.Println("The server is ready, open the Sturdy App to get started! 🐣")
	c.Printf("Hostname: %s", "http://localhost:30080/") // TODO: Actually make this dynamic
	c.Println()
	c.Println()
}
