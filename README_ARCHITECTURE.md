# Sturdy: Architecture

## Terminology

* **Codebase** - Contains the history of the project, has many workspaces, members, and integrations
* **Workspace** - A logical separation of work-in-progress (think: git-branch). When the workspace is not actively used by a view, it's state is saved in a "snapshot" (tracked in the database).
* **Suggesting workspaces** - A "fork" of a workspace, created from a snapshot of a workspace. Creates suggestions towards the original workspace.
* **View** - A checkout of a workspace on disk.
* **Change** - Landed changes from a workspace on top of trunk
* **Organizations** - Holds one or more codebases, has members. In enterprise and cloud builds of Sturdy the organization is responsible for billing.
