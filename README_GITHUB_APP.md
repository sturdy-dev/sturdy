# Sturdy GitHub App locally

> TODO: This needs to be updated.

Step 0, install localtunnel (or similar):

```
npm install -g localtunnel
lt --port 3000
```

1. Create a new app - https://github.com/settings/apps/new
2. Set the name to `sturdy-$YOUR_NAME [localhost]`
3. Set the homepage to `http://localhost:8080`
4. Callback URL: `http://localhost:8080/setup-github`
5. **Untick** `Expire user authorization tokens`
6. **Tick** `Request user authorization (OAuth) during installation` 
7. **Tick** `Redirect on update`
8. Make sure that webhooks are active
9. Set the webhook URL to `${YOUR_LOCALTUNNEL_HOSTNAME}/v3/github/webhook`

### Permissions

* Checks - Read & Write
* Contents - Read & Write
* Metadata - Read only
* Pull Requests - Read & Write
* Workflows - Read & Write

### Events

* Pull Request
* Push
* Pull request review

10. After creating the app, generate a new "Client secret"
11. Start the Sturdy backend with the following flags set (replace values to fit your app)

```
--github-app-id 97700
--github-app-client-id Iv1.814e400000000
--github-app-secret eecc220debbbxxxxxxxxxxxxxxxxxxxxx
--github-app-private-key-path ~/Downloads/sturdy-localhost.2021-07-06.private-key.pem
```

Aaaaand that's it! Happy hacking!