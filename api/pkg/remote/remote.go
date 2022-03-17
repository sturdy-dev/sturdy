package remote

type Remote struct {
	ID                string `db:"id"`
	CodebaseID        string `db:"codebase_id"`
	Name              string `db:"name"`
	URL               string `db:"url"`
	BasicAuthUsername string `db:"basic_username"`
	BasicAuthPassword string `db:"basic_password"`
	TrackedBranch     string `db:"tracked_branch"`
}
