# Changelog

> **From 🥚 to 🐣.**

This changelog contains the changelog for self-hosted Sturdy (OSS and Enterprise).  

Sturdy in the Cloud is continuously deployed, and will contain newer features not yet available in a release.  

Releases are pushed to [Docker Hub](https://hub.docker.com/r/getsturdy/server/).

# Server v1.2.0 (2022-02-18)

* Improved how Sturdy imports changes from GitHub – Merge commits are now correctly identified and converted to `changes`.
* Fix invite-links for self hosted installations.
* Enabled garbage collection of unused objects – This significantly improves the performance of installations with many snapshots.
* Fixed an issue in OSS builds where new users would not be added to the servers organization.
* ... and many more internal changes!

# App v0.5.0 (2022-02-17)

* Support connecting to _any_ Sturdy server – Access the settings with `Cmd+,` / `Ctrl+,` or from the Sturdy icon in the menubar or system tray. Requires Sturdy server 1.1.0 or newer.
* ... and many more internal changes!

# Server v1.0.0 (2022-02-08)

Sturdy is now Open Source, and can be self-hosted! 

* [Run Sturdy anywhere](https://getsturdy.com/docs/self-hosted), with self-hosted Sturdy!
* Licensed under Apache 2.0 and the Sturdy Enterprise License
