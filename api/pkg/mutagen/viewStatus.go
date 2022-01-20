package mutagen

import "time"

type ViewStatus struct {
	ID                    string          `db:"id"`
	State                 ViewStatusState `db:"state"`
	StagingStatusPath     *string         `db:"staging_status_path"`
	StagingStatusReceived *int            `db:"staging_status_received"`
	StagingStatusTotal    *int            `db:"staging_status_total"`
	SturdyVersion         string          `db:"sturdy_version"`
	LastError             *string         `db:"last_error"`
	UpdatedAt             time.Time       `db:"updated_at"`
}

// ViewStatusState is the Mutagen State
//
// https://github.com/sturdy-dev/mutagen/blob/09424246c56dd317a946e3c7b9e99bdd1f8879ed/pkg/synchronization/state.proto#L12
type ViewStatusState string

const (
	ViewStatusStateDisconnected           = "Disconnected"
	ViewStatusStateHaltedOnRootEmptied    = "HaltedOnRootEmptied"
	ViewStatusStateHaltedOnRootDeletion   = "HaltedOnRootDeletion"
	ViewStatusStateHaltedOnRootTypeChange = "HaltedOnRootTypeChange"
	ViewStatusStateConnectingAlpha        = "ConnectingAlpha"
	ViewStatusStateConnectingBeta         = "ConnectingBeta"
	ViewStatusStateWatching               = "Watching"
	ViewStatusStateScanning               = "Scanning"
	ViewStatusStateWaitingForRescan       = "WaitingForRescan"
	ViewStatusStateReconciling            = "Reconciling"
	ViewStatusStateStagingAlpha           = "StagingAlpha" // Downloading
	ViewStatusStateStagingBeta            = "StagingBeta"  // Uploading
	ViewStatusStateTransitioning          = "Transitioning"
	ViewStatusStateSaving                 = "Saving"
)
