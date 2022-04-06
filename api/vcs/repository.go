package vcs

import (
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"
)

// RepoGitReader only need access to read .git on the filesystem.
type RepoGitReader interface {
	CodebaseID() codebases.ID
	IsTrunk() bool
	ViewID() *string
	IsRebasing() bool

	CurrentDiffNoIndex(opts ...DiffOption) (*git.Diff, error)
	DiffCommits(firstCommitID, secondCommitID string) (*git.Diff, error)
	DiffCommitToRoot(commitID string) (*git.Diff, error)
	DiffRootToCommit(commitID string) (*git.Diff, error)

	RemoteBranchCommit(remoteName, branchName string) (*git.Commit, error)

	GetDefaultBranch() (string, error)

	HeadBranch() (string, error)
	HeadCommit() (*git.Commit, error)

	BranchCommitID(branchName string) (string, error)

	GetCommitParents(commitID string) ([]string, error)
	CommitMessage(id string) (author *git.Signature, message string, err error)
	ShowCommit(id string) (diffs []string, entry *LogEntry, err error)
	GetCommitDetails(id string) (*CommitDetails, error)
	BranchHasCommit(branchName, commitID string) (bool, error)
	CommonAncestor(commitA, commitB string) (string, error)

	FileContentsAtCommit(commitID, filePath string) ([]byte, error)
	FileBlobAtCommit(commitID, filePath string) (*git.Blob, error)
	DirectoryChildrenAtCommit(commitID, directoryPath string) ([]string, error)

	LogHead(limit int) ([]*LogEntry, error)
	LogBranch(branchName string, limit int) ([]*LogEntry, error)

	OpenRebase() (*SturdyRebase, error)
}

// RepoGitWriter can read and write to .git
type RepoGitWriter interface {
	RepoGitReader

	CreateRootCommit() error
	CommitIndexTree(treeID *git.Oid, message string, signature git.Signature) (string, error)
	CommitIndexTreeWithReference(treeID *git.Oid, message string, signature git.Signature, ref string) (string, error)

	CreateBranchTrackingUpstream(branchName string) error
	DeleteBranch(name string) error
	CreateNewBranchOnHEAD(name string) error
	CreateNewBranchAt(name string, targetSha string) error
	CreateNewCommitBasedOnCommit(newBranchName string, existingCommitID string, signature git.Signature, message string) (string, error)

	CleanStaged() error
	Push(logger *zap.Logger, branchName string) error
	ForcePush(logger *zap.Logger, branchName string) error
	FetchBranch(branches ...string) error

	PushNamedRemoteWithRefspec(remoteName string, creds transport.AuthMethod, refspecs []config.RefSpec) (userError string, err error)
	PushRemoteUrlWithRefspec(remoteUrl string, creds transport.AuthMethod, refspecs []config.RefSpec) (userError string, err error)
	FetchNamedRemoteWithCreds(remoteName string, creds transport.AuthMethod, refspecs []config.RefSpec) error
	FetchUrlRemoteWithCreds(remoteUrl string, creds transport.AuthMethod, refspecs []config.RefSpec) error

	SetDefaultBranch(targetBranch string) error
	CreateAndSetDefaultBranch(headBranchName string) error

	CreateCommitWithFiles(files []FileContents, newBranchName string) (string, error)

	ResetMixed(commitID string) error

	GitGC() error
	GitReflogExpire() error
	GitRemotePrune(remoteName string) error

	MergeBranches(ourBranchName, theirBranchName string) (*git.Index, error)
	MergeBranchInto(branchName, mergeIntoBranchName string) (mergeCommitId string, err error)

	ApplyPatchesToIndex(patches [][]byte) (*git.Oid, error)
}

type RepoReaderGitWriter interface {
	RepoReader
	RepoGitWriter
}

// RepoReader needs to read repository files on the filesystem.
type RepoReader interface {
	RepoGitReader

	Path() string

	Diffs(...DiffOption) (*git.Diff, error)
	CurrentDiff(opts ...DiffOption) (*git.Diff, error)

	AddFilesToIndex(files []string) (*git.Oid, error)
	AddAndCommit(msg string) (string, error)

	LargeFilesClean(codebaseID codebases.ID, paths []string) ([][]byte, error)

	CanApplyPatch(patch []byte) (bool, error)
}

// RepoWriter might modify files on the filesystem.
type RepoWriter interface {
	RepoReader
	RepoGitWriter

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

	AddNamedRemote(name, url string) error
	CreateRef(name, commitSha string) error
}
