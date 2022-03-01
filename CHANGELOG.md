# Changelog

> **From 🥚 to 🐣.**

This changelog contains the changelog for self-hosted Sturdy (OSS and Enterprise).  

Sturdy in the Cloud is continuously deployed, and will contain newer features not yet available in a release.  

Releases are pushed to [Docker Hub](https://hub.docker.com/r/getsturdy/server/).

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
