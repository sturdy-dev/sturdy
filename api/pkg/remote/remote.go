package remote

import (
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/crypto"
)

type Remote struct {
	ID                string            `db:"id"`
	CodebaseID        codebases.ID      `db:"codebase_id"`
	Name              string            `db:"name"`
	URL               string            `db:"url"`
	BasicAuthUsername *string           `db:"basic_username"`
	BasicAuthPassword *string           `db:"basic_password"`
	KeyPairID         *crypto.KeyPairID `db:"keypair_id"`
	TrackedBranch     string            `db:"tracked_branch"`
	BrowserLinkRepo   string            `db:"browser_link_repo"`
	BrowserLinkBranch string            `db:"browser_link_branch"`
}
