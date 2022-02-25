# Developing Sturdy

## Overview

Sturdy has 4 main components:

| Name          | Description                                                  | Links                                                                                 |
|:--------------|:-------------------------------------------------------------|:--------------------------------------------------------------------------------------|
| `api`         | The backend. Does all of the version control. Built with Go  | [Sources](./api)                                                                      |
| `web`         | The website. Built with TypesScript, Vue 3, urql, and more.  | [Sources](./web)                                                                      |
| `sturdy-sync` | The file synchronizer. A fork of mutagen.                    | [Sources](./ssh), [The Sturdy fork of Mutagen](https://github.com/sturdy-dev/mutagen) |
| `app`         | The electron app. Runs sturdy-sync, and renders the web app. | [Sources](./app)                                                                      |


```
# This is a simplified diagram of how all components connect to each other
# Protocols are annotated [like this].

 Electron ──[HTTPS]───► Web ───[GraphQL]────► API
    │                                          ▲
    │                                          │
    │                                      [REST/JSON]
    │                                          │  
    ▼                                          │
Local-Syncer ─────────[SSH]────────────► Remote-Syncer
```

## Easy development

Sturdy has a Docker container called the "oneliner", which contains all components of Sturdy in a single easy-to-run container.
This is the easiest way to get a full development environment for all components except for the Electron App (most of the time however, connecting the production build of the Sturdy app is good enough). 

```bash
./scripts/run-oneliner.sh
```

## Development

To support a full development environment, with hot reloading and fast restarts. 
* Ensure libgit2 is installed - `https://libgit2.org/`
* Run PostgreSQL, LFS, and the SSH servers in Docker: `./up --build`
* Build and run the API
  server: `cd api && go build getsturdy.com/api/cmd/api && ./api --http-listen-addr 127.0.0.1:3000 --analytics.enabled=false`
* Build and run the web frontend: `cd web && yarn && yarn codegen && yarn dev`
* Build and run the Electron app: `cd app && yarn && yarn dev`

## Use the GraphQL API

The easiest way to develop against the API-server is to test it through the GraphQL API.

1. Get a auth token from `http://localhost:8080/api` (when logged in).
2. Use a GraphQL client like Insomnia or Altair.
3. Connect to `http://localhost:3000/graphql`
4. Authenticate with `Authentication: bearer $YOUR_TOKEN`
5. You might have to set the Origin header to `http://localhost:8080`

## Testing...

### the API Server

```bash
# Run unit tests
cd api && go test getsturdy.com/api/...

# Run unit tests in Docker
docker compose -f ci/docker-compose.yaml -f ci/unit/docker-compose.yaml up --build --exit-code-from runner

# Run unit + E2E tests
docker compose -f ci/docker-compose.yaml -f ci/e2e/docker-compose.yaml up --build --exit-code-from runner
```

### the website

```bash
yarn test && yarn lint
```

## _[Footnote]_

For Sturdy employees: [How to build and run Sturdy Cloud](https://docs.google.com/document/d/1GFk2liBUL8xqbEacVpX7mPkuJKMXh1RHYqRRp9qrQ5Q/edit)
