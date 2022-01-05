# Sturdy: Architecture

* [Overview: Diagram](https://docs.google.com/drawings/d/1M0kAMxzDiFVH93as-n01Do3buBGvFwr1k3F_OHLNlIg/edit)

Internally, Sturdy is largely based on git.

## Terminology

* Codebase
* Workspace - Where a developer (represented by a git-branch). When the workspace is not actively used by a view, it's state is saved in a "snapshot" (tracked in the database).
* View - A checkout of a workspace (branch) on disk, has at most 1 mutagen connections. A workspace can have an unlimited number of views, but only one of them is "authoritative":
* Authoritative View - The owner of a workspace can only have one active view per workspace. This is the only view that changes can be created from, and suggestions can be accepted into.
* sturdytrunk - Name of the main branch
* Change - A git commit, if the change is rebased, the change can have many git commits
  
## Git File System

All repositories and views are currently located on a single disk (mounted at `/repos`), with the following structure:

```
/$CODEBASE_ID
-------------/trunk - A "bare" git checkout. Contains the sturdytrunk, all workspaces and snapshots
                    - The "origin" remote can be a connected GitHub repository
                    - The Sturdy API reads from this git repository
-------------/$VIEW_ID - A checkout of the repository, on a workspace branch.
                       - Mutagen writes to this directory.
                       - The Sturdy API performs git actions here, via libgit2 and the git cli
                       - The "origin" remote is "../trunk"
```

## Tools

* [Mutagen](https://mutagen.io/)