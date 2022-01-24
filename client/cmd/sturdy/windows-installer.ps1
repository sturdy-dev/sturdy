$APPDATA=[Environment]::GetFolderPath('ApplicationData')
New-Item -Path "$APPDATA" -Name "sturdy" -ItemType "directory" -Force | Out-Null

$VERSION="v0.8.1-beta"

# Stop sturdy if running
[bool] $STURDY_ALREADY_INSTALLED = 0
if (Get-Command sturdy -errorAction SilentlyContinue)
{
    $STURDY_ALREADY_INSTALLED=1
    echo "Stopping Sturdy"
    sturdy stop
    # Stop-Process -Name "sturdy-sync" -Force
    # rm .\.sturdy-sync\daemon\daemon.lock
}

# Delete an existing bin directory if one exists
Remove-Item "$APPDATA\sturdy\bin" -Recurse -ErrorAction Ignore | Out-Null
Remove-Item "$APPDATA\sturdy\bin" -Recurse -ErrorAction Ignore | Out-Null

echo "Downloading Sturdy $VERSION"
Invoke-WebRequest -Uri https://getsturdy.com/client/sturdy-$VERSION-windows-amd64.zip -OutFile "$APPDATA\sturdy\sturdy-$VERSION-windows-amd64.zip"

echo "Installing Sturdy $VERSION"
Add-Type -AssemblyName System.IO.Compression.FileSystem ; [System.IO.Compression.ZipFile]::ExtractToDirectory("$APPDATA\sturdy\sturdy-$VERSION-windows-amd64.zip", "$APPDATA\sturdy\bin")

$STURDY_PATH="$APPDATA\sturdy\bin"

if ( ! $env:PATH.Contains("$STURDY_PATH") ) {
    echo "Updating PATH (for this user only)"

    # Add to permanent environment user $PATH
    # Does not need to run as an administrator
    [Environment]::SetEnvironmentVariable(
        "Path",
        [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) + ";$STURDY_PATH",
        [EnvironmentVariableTarget]::User)

    # Add to this shells PATH
    $env:Path += ";$STURDY_PATH"
} else {
    echo "Your PATH is already set, skipping."
}

echo "Sturdy has been installed to: $STURDY_PATH"
echo "You're now ready to use Sturdy!"

sturdy version

if ( $STURDY_ALREADY_INSTALLED ) {
    echo "Starting Sturdy..."
    sturdy start
} else {
    echo "Please restart the terminal for the changes to the PATH to take effect"
}
