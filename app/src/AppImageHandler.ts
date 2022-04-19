import * as fs from 'fs'
import { app } from 'electron'
import * as child_process from 'child_process'

export const CreateAppImageDesktopFile = (path: string) => {
  console.log('Creating .desktop file')

  const contents = `[Desktop Entry]
Type=Application
Name=Sturdy
Exec=${path} %u
StartupNotify=false
MimeType=text/html;x-scheme-handler/sturdy;`

  fs.writeFileSync(`${app.getPath('home')}/.local/share/applications/sturdy.desktop`, contents)

  // update desktop database
  child_process.exec('update-desktop-database ~/.local/share/applications')
}
