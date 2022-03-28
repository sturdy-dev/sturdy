# App

This directory is the root directory for the Sturdy desctop app which is using Electron.

## Versioning

| Mask          | Meaning                          |
| ------------- | -------------------------------- |
| x.x.x         | stable releases                  |
| x.x.x-beta.x  | beta (stable, but needs testing) |
| x.x.x-alpha.x | alpha (unstable)                 |

## How to publish a new version

1. Update verion in [package.json](./package.json)
2. `$ ./build-electron-builder.sh --upload`
