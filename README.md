# Welcome to Sturdy! 📣🐢

## Links

Here's a collection to tools, services, and so on that we're using. Everyone should have access to all of these.

* [Kitemaker](https://toil.kitemaker.co/YtReAM-Sturdy/Xub625-Sturdy/boards/current) - Planning
* [AWS](https://sturdy.awsapps.com/start#/)
* * [Grafana](https://g-475af40809.grafana-workspace.eu-west-1.amazonaws.com/login) - Monitoring and alerting
* * [Logs: API](https://eu-north-1.console.aws.amazon.com/cloudwatch/home?region=eu-north-1#logsV2:log-groups/log-group/driva/log-events)
* * [Logs: SSH-server](https://eu-north-1.console.aws.amazon.com/cloudwatch/home?region=eu-north-1#logsV2:log-groups/log-group/mutagen-ssh/log-events)
* * [Guide: Filter CloudWatch Logs](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html)
* [PostHog](https://app.posthog.com/dashboard/2592) - Analytics
* [Slack](https://join.slack.com/t/getsturdy/signup) - Chatting!
* [Discord](https://discord.gg/5HnSdzMqtA) - Community Server
* [Umami](https://umami.getsturdy.com/realtime) - Real time web analytics
* [Retool](https://sturdy.retool.com/) - Backoffice

## Development

```bash
# Run unit tests
docker compose -f ci/docker-compose.yaml -f ci/unit/docker-compose.yaml up --build --exit-code-from runner

# Run unit + E2E tests
docker compose -f ci/docker-compose.yaml -f ci/e2e/docker-compose.yaml up --build --exit-code-from runner
```

It's possible to run the services without Docker.

* Run PostreSQL, LFS, and the SSH servers: `./up`
* Build and run the API server: `go build -tags enterprise,cloud -v -o mash mash/cmd/api && ./mash --http-listen-addr 127.0.0.1:3000 --unauthenticated-graphql-introspection`
* Run the web application: `cd /web && yarn install && yarn run dev`.
* Invoke the `sturdy` CLI with `--config` set, to override the default configuration (example below). Use like so: `sturdy auth --config ~/.sturdy-local`

### Example `sturdy` (the CLI) configuration
```json
{
  "remote": "127.0.0.1:3001",
  "insecure-remote": true,
  "api-remote": "http://127.0.0.1:3000",
  "sync-remote": "127.0.0.1:2222",
  "git-remote": "127.0.0.1:3002"
}
```

## GraphQL API

For authentication and CORS, set the following headers:

Get the auth cookie from a session in your browser (Developer Tools (Option + Cmd + I) -> Application -> Cookies -> `auth`)

```
Cookie: auth=YOUR_JWT_GOES_HERE
Origin: http://localhost:8080
```

[Altair](https://altair.sirmuel.design/) is a great desktop client to explore/test the graph.

😎

<!-- Test: 12 -->
<!-- 2021-11-23 - Hello from Electron/Windows! -->
