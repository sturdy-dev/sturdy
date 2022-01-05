package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/kolide/launcher/pkg/packagekit"

	"github.com/google/uuid"
	"github.com/kolide/launcher/pkg/packagekit/authenticode"
	"github.com/kolide/launcher/pkg/packagekit/wix"
	"go.opencensus.io/trace"
)

func main() {
	flagRoot := flag.String("root", "C:\\Users\\KirilVidelov\\sturdy\\cmd\\sturdy\\msi\\contents", "")
	flagVersion := flag.String("version", "", "")
	flag.Parse()

	// MSI versions can only contain numbers and periods
	versionWithoutV := strings.Trim(*flagVersion, "v")

	opts := &packagekit.PackageOptions{
		Name:       "Sturdy",
		Identifier: "sturdy",
		Root:       *flagRoot,
		Version:    versionWithoutV,
		WixPath:    "C:\\Program Files (x86)\\WiX Toolset v3.11\\bin",
		WixUI:      true,
	}

	var output bytes.Buffer
	err := PackageWixMSI(context.Background(), &output, opts, false)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("sturdy.msi", output.Bytes(), 0o644)
	if err != nil {
		panic(err)
	}
}

// We need to use variables to stub various parts of the wix
// xml. While we could use wix's internal variable system, it's a
// little more debugable to do it with go's. This way, we can
// inspect the intermediate xml file.
//
// This might all be cleaner moved from a template to a marshalled
// struct. But enumerating the wix options looks very ugly
//
//go:embed assets/main.wxs
var wixTemplateBytes []byte

// This is used for icons and splash screens and the like. It would be
// better in pkg/packaging, and passed into packagekit, but that's a
// deeper refactor.
//
//go:embed assets/*
var assets embed.FS

const (
	signtoolPath = `C:\Program Files (x86)\Windows Kits\10\bin\10.0.18362.0\x64\signtool.exe`
)

func PackageWixMSI(ctx context.Context, w io.Writer, po *packagekit.PackageOptions, includeService bool) error {
	ctx, span := trace.StartSpan(ctx, "packagekit.PackageWixMSI")
	defer span.End()

	if err := isDirectory(po.Root); err != nil {
		return err
	}

	// We include a random nonce as part of the ProductCode
	// guid. This is so that any MSI rebuild triggers the Major
	// Upgrade flow, and not the "Another version of this product
	// is already installed" error. The Minor Upgrade Flow might
	// be more appropriate, but requires substantial reworking of
	// how versions and builds are calculated. See
	// https://www.firegiant.com/wix/tutorial/upgrades-and-modularization/
	// for opinionated background
	guidNonce, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generating uuid as guid nonce: %w", err)

	}
	extraGuidIdentifiers := []string{
		po.Version,
		runtime.GOARCH,
		guidNonce.String(),
	}

	var templateData = struct {
		Opts        *packagekit.PackageOptions
		UpgradeCode string
		ProductCode string
	}{
		Opts:        po,
		UpgradeCode: generateMicrosoftProductCode("launcher" + po.Identifier),
		ProductCode: generateMicrosoftProductCode("launcher"+po.Identifier, extraGuidIdentifiers...),
	}

	wixTemplate, err := template.New("WixTemplate").Parse(string(wixTemplateBytes))
	if err != nil {
		return fmt.Errorf("not able to parse main.wxs template: %w", err)
	}

	mainWxsContent := new(bytes.Buffer)
	if err := wixTemplate.ExecuteTemplate(mainWxsContent, "WixTemplate", templateData); err != nil {
		return fmt.Errorf("executing WixTemplate: %w", err)
	}

	wixArgs := []wix.WixOpt{}

	if po.WixSkipCleanup {
		wixArgs = append(wixArgs, wix.SkipCleanup())
	}

	if po.WixPath != "" {
		wixArgs = append(wixArgs, wix.WithWix(po.WixPath))
	}

	{
		// Regardless of whether or not there's a UI in the MSI, we
		// still want the icon file to be included.
		assetFiles := []string{"sturdy.ico"}

		if po.WixUI {
			assetFiles = append(assetFiles, "msi_banner.bmp", "msi_splash.bmp")
			wixArgs = append(wixArgs, wix.WithUI())
		}

		for _, f := range assetFiles {
			fileBytes, err := assets.ReadFile("assets/" + f)
			if err != nil {
				return fmt.Errorf("getting asset %s: %w", f, err)
			}

			wixArgs = append(wixArgs, wix.WithFile(f, fileBytes))
		}
	}

	if includeService {
		launcherService := wix.NewService("launcher.exe",
			wix.WithDelayedStart(),
			wix.ServiceName(fmt.Sprintf("Launcher%sSvc", strings.Title(po.Identifier))),
			wix.ServiceArgs([]string{"svc", "-config", po.FlagFile}),
			wix.ServiceDescription(fmt.Sprintf("The Kolide Launcher (%s)", po.Identifier)),
		)

		if po.DisableService {
			wix.WithDisabledService()(launcherService)
		}

		wixArgs = append(wixArgs, wix.WithService(launcherService))
	}

	wixTool, err := wix.New(po.Root, mainWxsContent.Bytes(), wixArgs...)
	if err != nil {
		return fmt.Errorf("making wixTool: %w", err)
	}
	defer wixTool.Cleanup()

	// Use wix to compile into an MSI
	msiFile, err := wixTool.Package(ctx)
	if err != nil {
		return fmt.Errorf("wix packaging: %w", err)
	}

	// Sign?
	if po.WindowsUseSigntool {
		if err := authenticode.Sign(
			ctx, msiFile,
			authenticode.WithExtraArgs(po.WindowsSigntoolArgs),
			authenticode.WithSigntoolPath(signtoolPath),
		); err != nil {
			return fmt.Errorf("authenticode signing: %w", err)
		}
	}

	// Copy MSI into our filehandle
	msiFH, err := os.Open(msiFile)
	if err != nil {
		return fmt.Errorf("opening msi output file: %w", err)
	}
	defer msiFH.Close()

	if _, err := io.Copy(w, msiFH); err != nil {
		return fmt.Errorf("copying output: %w", err)
	}

	setInContext(ctx, ContextLauncherVersionKey, po.Version)

	return nil
}

// generateMicrosoftProductCode create a stable guid from a set of
// inputs. This is used to identify the product / sub product /
// package / version, and whatnot. We need to either store them, or
// generate them in a predictable fasion based on a set of inputs. See
// doc.go, or
// https://docs.microsoft.com/en-us/windows/desktop/Msi/productcode
//
// It is equivlent to uuid.NewSHA1(kolideUuidSpace,
// []byte(launcherkolide-app0.7.0amd64)) but provided here so we have
// a clear point to test stability against.
func generateMicrosoftProductCode(ident1 string, identN ...string) string {
	// Define a Kolide uuid space. This could also have used uuid.NameSpaceDNS
	uuidSpace := uuid.NewSHA1(uuid.Nil, []byte("Sturdy"))

	data := strings.Join(append([]string{ident1}, identN...), "")

	guid := uuid.NewSHA1(uuidSpace, []byte(data))

	return strings.ToUpper(guid.String())
}
