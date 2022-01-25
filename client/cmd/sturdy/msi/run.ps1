$VERSION=$args[0]

echo "version=$VERSION"

echo "Cleanup"
Remove-Item 'sturdy-windows' -Recurse -ErrorAction SilentlyContinue
Remove-Item 'sturdy-windows.zip' -ErrorAction SilentlyContinue

echo "Downloading"
Invoke-WebRequest -Uri "https://getsturdy.com/getsturdy.com/client/sturdy-$VERSION-windows-amd64.zip" -OutFile "sturdy-windows.zip"

echo "Copying"
Expand-Archive -LiteralPath 'sturdy-windows.zip' -DestinationPath sturdy-windows
Copy-Item "sturdy-windows\sturdy.exe" -Destination "cmd\sturdy\msi\contents"
Copy-Item "sturdy-windows\sturdy-sync.exe" -Destination "cmd\sturdy\msi\contents"

Push-Location cmd\sturdy\msi

echo "Building msi.exe"
go build -v

echo "Creating sturdy.msi"
$pwd=Get-Location
.\msi.exe --root "$pwd\contents" --version $VERSION

Pop-Location