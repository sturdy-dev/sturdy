<p align="center"><img src="https://getsturdy.com/assets/Web/Logo/DuckAndName.png" height="128"></p>

# Welcome to Sturdy! 📣🐣

**[Sturdy](https://getsturdy.com/) is the easiest way to collaborate on code.**

## Features

* Interact with version control at a higher abstraction level (e.g. no need for pushing, pulling, stashing etc.)
* Discover work-in-progress code within your team in real-time
* Try your teammate's code locally with a single click
* Suggest code changes / ideas to a colleague by simply typing in your IDE
* Cloud or self-hosted!
* Enhance your existing GitHub setup, or _break free_ and use standalone Sturdy

## Versions
* [Sturdy Cloud](https://getsturdy.com/) - Let's you use all Sturdy features, fully managed by the team behind Sturdy. Ship code to your projects, review, and ship code. Using 100% Sturdy, or Sturdy on top of GitHub. Get started for **free**.
* [Sturdy Enterprise](#) - Run Sturdy in your own environment.
* [Sturdy OSS](#) - The fully Open Source version of Sturdy! Provides all the core functionality for free, and completely Open Source.

## Get Started

Want to run Sturdy on your machine?

```bash
docker run --interactive --tty --publish 30080:80 --publish 30022:22 \
    --volume "$HOME/.sturdydata:/var/data" \
    getsturdy/server
```

## Learn more

See the [Sturdy Docs](https://getsturdy.com/v2/docs) to learn more about Sturdy and why it's cool!

## Development

See [README_DEVELOPMENT.md](README_DEVELOPMENT.md) for instructions of how to build and develop Sturdy.

## We're hiring!

Come and help make Sturdy even better! We're growing and are [hiring for multiple positions](https://getsturdy.com/careers).

## License

This repository contains both OSS-licensed and non-OSS-licensed files.

All files under any directory named `enterprise` fall under [LICENSE.enterprise](LICENSE.enterprise).

The remaining files fall under [TODO]. 

<!-- Test: 13 -->
<!-- 2021-11-23 - Hello from Electron/Windows! -->
