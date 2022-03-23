# Changelog

> **From 🥚 to 🐣.**

This changelog contains the changelog for self-hosted Sturdy (OSS and Enterprise).  

Sturdy in the Cloud is continuously deployed, and will contain newer features not yet available in a release.  

Releases are pushed to [Docker Hub](https://hub.docker.com/r/getsturdy/server/).

# Server v1.5.0 (2022-03-23)

_Changes since v1.4.0_

* [Improvement] Added git-integration with support for _any_ git hosting provider (GitLab, BitBucket, Azure Dev Ops, etc...). Available from "Codebase Settings > Integrations > Git" (BETA).
* [Improvement] GitHub-app self-verification. Verifies that the app is setup correctly with the necessary permissions.
* [Improvement] Workspaces are now completely renamed to "Draft Changes"
* [Improvement] Much faster workflow when running Sturdy on GitHub
* [Improvement] Fixed issues with the built-in search on Drafts and Changes (open the search window with `Cmd+F` or `/`)
* [Improvement] Search for file names in the diff searcher
* [Improvement] The `getsturdy/server` Docker image now runs on both `linux/amd64` and `linux/arm64` (new).
* [Improvement] TLS support for self-hosted servers, with built in Let's Encrypt support
* [Performance] Improved performance when resizing the app-window.
* [Fix] Improved speed and reliability of updates of events in the sidebar
* [Fix] Fixed some rendering issues in the Syncer
* [Fix] Fixed issues with uploading avatars on self-hosted servers
* [Fix] Fixed a bug where a repository could not be imported twice from GitHub

# Server v1.4.0 (2022-03-08)

* [Improvement] Improvements to the workflow. When sharing a change to trunk, you're now redirected to the change that you created. A workspace (to be renamed to "Draft Changes" in an upcoming release) can now only be used once. After sharing a change Sturdy will automatically create a new workspace on trunk and connect it to your computer.
* [Improvement] Allow to rename a organization
* [Improvement] Notifications for comments on changes are now sent
* [Improvement] See the activity feed for changes
* [Fix] Snapshots are now more reliable
* [Fix] Improved reliability and fixed a panic in the Events Sender (powering GraphQL Subscriptions)
* [Fix] Real-time diff streaming is now more reliable

# Server v1.3.0 (2022-03-01)

* [Improvement] Improved reliability when starting the oneliner for the first time. Reduced number of timeouts related to setting up the bundled PostgreSQL server.
* [Improvement] Added a new changelog overview for codebases.
* [Fix] Fixed a bug in the app sidebar navigation where clicking a codebase would sometimes not navigate you to the codebase.

# Server v1.2.1 (2022-02-25)

* [Improvement] Make sure that a connected directory always is connected to a workspace. If a workspace connected to a directory is archived, a new workspace will be created and connected to that directory.
* [Improvement] Inactive and unused workspace by other users are now hidden in the sidebar
* [Fix] Fixed an issue where some changes imported from GitHub where not revertable
* [Fix] Fixed an issue where some workspaces did not have a "Based On" change tracked
* [Fix] Fixed an issue with changes that contained files that where (at the same time) renamed, edited, and had new file permissions.
* [Fix] Fixed an issue where it was not possible to make comments in a workspace (live comments) on deleted lines in files that have been moved.
* [Fix] Fixed an issue where it was not possible to upload custom user avatars
* [Performance] Fetching and loading the changelog is now faster

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
