# Oneliner

This directory contains [s6-overlay][] configuration to run all Sturdy components inside a single docker container.

## Components:

- `api` is the Sturdy api server
- `postgresql` is a Postgres database instance
- `repos` runs scripts to setup data directory for repositories
- `reproxy` is a simple http proxy server to serve both SPA frontend and API server on the same port
- `rudolfs` is a git-lfs backend
- `ssh` is the Sturdy ssh server for mutagen connections
- `sslmux` is a tcp multiplexer to serve ssh and http on the same port
- all `*-log` components are used only for the logging perposes. All they do is they append name of the service to the
  log message
- `*-prepare` and `repos` run scripts to prepare data directories

[s6-overlay]: https://github.com/just-containers/s6-overlay/
