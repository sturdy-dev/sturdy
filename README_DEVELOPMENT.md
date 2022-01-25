# Developing Sturdy

> TODO: Adopt this page for open-source Sturdy

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

It's possible to run the services without Docker.

* Run PostgreSQL, LFS, and the SSH servers: `./up --build`
* Build and run the API server: `cd api && go build getsturdy.com/api/cmd/api && ./api --http-listen-addr 127.0.0.1:3000 --analytics.enabled=false`
* Build and run the web frontend: `cd web && yarn && yarn codegen && yarn dev`
* Build and run the Electron app: `cd app && yarn && yarn dev`

## Testing

```bash
# Run unit tests
cd api && go test mash/...

# Run unit tests in Docker
docker compose -f ci/docker-compose.yaml -f ci/unit/docker-compose.yaml up --build --exit-code-from runner

# Run unit + E2E tests
docker compose -f ci/docker-compose.yaml -f ci/e2e/docker-compose.yaml up --build --exit-code-from runner
```

## GraphQL API

For authentication and CORS, set the following headers:

Get the auth cookie from a session in your browser (Developer Tools (Option + Cmd + I) -> Application -> Cookies -> `auth`)

```
Cookie: auth=YOUR_JWT_GOES_HERE
Origin: http://localhost:8080
```

[Altair](https://altair.sirmuel.design/) is a great desktop client to explore/test the graph.
