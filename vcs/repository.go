package vcs

import (
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

// Repo only need access to .git on the filesystem.
type Repo interface {
	Free()
	CodebaseID() string
	IsTrunk() bool
	ViewID() *string
	IsRebasing() bool

	CurrentDiffNoIndex() (*git.Diff, error)
	DiffCommits(firstCommitID, secondCommitID string) (*git.Diff, error)

	CreateRootCommit() error
	CommitIndexTree(treeID *git.Oid, message string, signature git.Signature) (string, error)
	CommitIndexTreeWithReference(treeID *git.Oid, message string, signature git.Signature, ref string) (string, error)
	RemoteBranchCommit(remoteName, branchName string) (*git.Commit, error)

	CreateBranchTrackingUpstream(branchName string) error
	DeleteBranch(name string) error
	CreateNewBranchOnHEAD(name string) error
	CreateNewBranchAt(name string, targetSha string) error
	CreateNewCommitBasedOnCommit(newBranchName string, existingCommitID string, signature git.Signature, message string) (string, error)

	CleanStaged() error

	Push(logger *zap.Logger, branchName string) error
	ForcePush(logger *zap.Logger, branchName string) error
	PushNamedRemoteWithRefspec(logger *zap.Logger, remoteName string, creds git.CredentialsCallback, refspecs []string) (userError string, err error)

	RemoteFetchWithCreds(remoteName string, creds git.CredentialsCallback, refspecs []string) error
	FetchBranch(branches ...string) error

	SetDefaultBranch(targetBranch string) error
	GetDefaultBranch() (string, error)
	CreateAndSetDefaultBranch(headBranchName string) error

	HeadBranch() (string, error)
	HeadCommit() (*git.Commit, error)

	BranchCommitID(branchName string) (string, error)
	BranchFirstNonMergeCommit(branchName string) (string, error)

	GetCommitParents(commitID string) ([]string, error)
	CommitMessage(id string) (author *git.Signature, message string, err error)
	ShowCommit(id string) (diffs []string, entry *LogEntry, err error)
	BranchHasCommit(branchName, commitID string) (bool, error)

	CreateCommitWithFiles(files []FileContents, newBranchName string) (string, error)

	FileContentsAtCommit(commitID, filePath string) ([]byte, error)
	FileBlobAtCommit(commitID, filePath string) (*git.Blob, error)
	DirectoryChildrenAtCommit(commitID, directoryPath string) ([]string, error)

	GitGC() error
	GitReflogExpire() error

	ResetMixed(commitID string) error

	LogHead(limit int) ([]*LogEntry, error)
	LogBranch(branchName string, limit int) ([]*LogEntry, error)

	MergeBranches(ourBranchName, theirBranchName string) (*git.Index, error)
	MergeBranchInto(branchName, mergeIntoBranchName string) error

	ApplyPatchesToIndex(patches [][]byte) (*git.Oid, error)

	RevertOnBranch(revertCommitID, branchName string) (string, error)

	OpenRebase() (*SturdyRebase, error)
}

// RepoReader needs to read repository files on the filesystem.
type RepoReader interface {
	Repo

	Path() string

	Diffs(...DiffOption) (*git.Diff, error)
	CurrentDiff(opts ...DiffOption) (*git.Diff, error)

	AddFilesToIndex(files []string) (*git.Oid, error)
	AddAndCommit(msg string) (string, error)

	LargeFilesClean(codebaseID string, paths []string) ([][]byte, error)

	CanApplyPatch(patch []byte) (bool, error)
}

// RepoWriter might modify files on the filesystem.
type RepoWriter interface {
	RepoReader

	CheckoutFile(fileName string) error
	DeleteFile(fileName string) error

	CheckoutBranchWithForce(branchName string) error
	CheckoutBranchSafely(branchName string) error
	CreateAndCheckoutBranchAtCommit(commitID, newBranchName string) error

	MoveBranchToCommit(branchName, targetCommitSha string) error
	MoveBranch(branchName, targetBranchName string) error
	MoveBranchToHEAD(branchName string) error

	CherryPickOnto(commitID, onto string) (newCommitID string, conflicted bool, conflictingFiles []string, err error)

	InitRebaseRaw(head, onto string) (*SturdyRebase, []RebasedCommit, error)

	LargeFilesPull() error

	ApplyPatchesToWorkdir(patches [][]byte) error

	ResetHard(commitID string) error
}
