package vcs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"

	git "github.com/libgit2/git2go/v33"
)

type repository struct {
	r           *git.Repository
	path        string
	lfsHostname string
}

func (repo *repository) CodebaseID() codebases.ID {
	return codebases.ID(filepath.Base(filepath.Dir(repo.path)))
}

func (repo *repository) IsTrunk() bool {
	return filepath.Base(repo.path) == "sturdytrunk"
}

func (repo *repository) ViewID() *string {
	if repo.IsTrunk() {
		return nil
	}
	viewID := filepath.Base(repo.path)
	return &viewID
}

var (
	meteredMethod = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "sturdy_git_call_millis",
			Buckets: []float64{
				.01, .025, .05,
				.1, .25, .5,
				1, 2.5, 5,
				10, 25, 50,
				100, 250, 500,
				1000, 2500, 5000,
				10000, 25000, 50000,
				100000, 250000, 500000,
			},
		}, []string{"method"})
)

func getMeterFunc(method string) func() {
	t0 := time.Now()
	return func() {
		meteredMethod.With(prometheus.Labels{"method": method}).Observe(float64(time.Since(t0).Milliseconds()))
	}
}

func CreateBareRepoWithRootCommit(path string) (*repository, error) {
	defer getMeterFunc("CreateBareRepoWithRootCommit")()

	if _, err := git.InitRepository(path, true); err != nil {
		return nil, err
	}
	r, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	if err := r.CreateRootCommit(); err != nil {
		return nil, fmt.Errorf("failed to create root commit in trunk: %w", err)
	}
	if err := r.CreateAndSetDefaultBranch("sturdytrunk"); err != nil {
		return nil, fmt.Errorf("failed to create new default branch: %w", err)
	}
	return r, nil
}

func CreateNonBareRepoWithRootCommit(path, headBranchName string) (*repository, error) {
	defer getMeterFunc("CreateNonBareRepoWithRootCommit")()
	_, err := git.InitRepository(path, false)
	if err != nil {
		return nil, err
	}
	r, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	if err := r.CreateRootCommit(); err != nil {
		return nil, fmt.Errorf("failed to create root commit in trunk: %w", err)
	}
	if err := r.CreateAndSetDefaultBranch(headBranchName); err != nil {
		return nil, fmt.Errorf("failed to create new default branch: %w", err)
	}
	return r, nil
}

func CreateEmptyBareRepo(path string) (*repository, error) {
	_, err := git.InitRepository(path, true)
	if err != nil {
		return nil, err
	}
	r, err := OpenRepo(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func CloneRepo(source, target string) (*repository, error) {
	cmd := exec.Command("git", "clone", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		return nil, err
	}
	return OpenRepo(target)
}

func CloneRepoBare(source, target string) (*repository, error) {
	cmd := exec.Command("git", "clone", "--bare", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		return nil, err
	}
	return OpenRepo(target)
}

func RemoteCloneWithCreds(url, target string, creds git.CredentialsCallback, bare bool) (*repository, error) {
	opts := git.CloneOptions{
		FetchOptions: git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				CredentialsCallback:      creds,
				CertificateCheckCallback: func(cert *git.Certificate, valid bool, hostname string) error { return nil },
			},
		},
		Bare: bare,
	}
	clonedRepo, err := git.Clone(url, target, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed cloning with creds: %w", err)
	}
	clonedRepo.Free()

	r, err := OpenRepo(target)
	if err != nil {
		return nil, fmt.Errorf("failed opening newly created repo at %s: %w", target, err)
	}
	return r, nil
}

// OpenRepo
// Deprecated: use OpenRepoWithLFS instead
func OpenRepo(path string) (*repository, error) {
	defer getMeterFunc("OpenRepo")()

	r, err := git.OpenRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo at path %s: %w", path, err)
	}
	return &repository{
		r:    r,
		path: path,
	}, nil
}

func OpenRepoWithLFS(path, lfsHostname string) (*repository, error) {
	defer getMeterFunc("OpenRepoWithLFS")()

	r, err := git.OpenRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo at path %s: %w", path, err)
	}
	return &repository{
		r:           r,
		path:        path,
		lfsHostname: lfsHostname,
	}, nil
}

func (repo *repository) Diffs(options ...DiffOption) (*git.Diff, error) {
	if repo.IsRebasing() {
		return repo.CurrentDiffNoIndex()
	}
	return repo.CurrentDiff(options...)
}

func (repo *repository) IsRebasing() bool {
	return repo.r.State() == git.RepositoryStateRebase || repo.r.State() == git.RepositoryStateRebaseMerge
}

// CommitIndexTree commits whatever is currently in the index from the treeID provided
func (repo *repository) CommitIndexTree(treeID *git.Oid, message string, signature git.Signature) (string, error) {
	defer getMeterFunc("CommitIndexTreeWithReference")()

	return repo.CommitIndexTreeWithReference(treeID, message, signature, "HEAD")
}

func (repo *repository) CommitIndexTreeWithReference(treeID *git.Oid, message string, signature git.Signature, ref string) (string, error) {
	defer getMeterFunc("CommitIndexTreeWithReference")()

	tree, err := repo.r.LookupTree(treeID)
	if err != nil {
		return "", fmt.Errorf("lookup tree failed: %w", err)
	}
	defer tree.Free()

	branch, err := repo.r.Head()
	if err != nil {
		return "", fmt.Errorf("get head failed: %w", err)
	}
	defer branch.Free()

	commitTarget, err := repo.r.LookupCommit(branch.Target())
	if err != nil {
		return "", fmt.Errorf("lookup commit failed: %w", err)
	}
	defer commitTarget.Free()

	oid, err := repo.r.CreateCommit(ref, &signature, &signature, message, tree, commitTarget)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}
	return oid.String(), nil
}

// AddFilesToIndex adds the given list of file paths to the index. The index is written.
// Returns the treeID or an error. The treeID can be passed to CommitIndexTree.
func (repo *repository) AddFilesToIndex(files []string) (*git.Oid, error) {
	defer getMeterFunc("AddFilesToIndex")()

	index, err := repo.r.Index()
	if err != nil {
		return nil, fmt.Errorf("failed to access vcs index: %w", err)
	}
	defer index.Free()

	err = index.AddAll(files, git.IndexAddDefault, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to add %d files: %w", len(files), err)
	}
	err = index.Write()
	if err != nil {
		return nil, err
	}
	oid, err := index.WriteTree()
	if err != nil {
		return nil, fmt.Errorf("write tree failed: %w", err)
	}
	return oid, nil
}

func (repo *repository) CheckoutFile(fileName string) error {
	defer getMeterFunc("CheckoutFile")()

	cmd := exec.Command("git", "checkout", "--", fileName)
	cmd.Dir = repo.path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		return err
	}
	return nil
}

func (repo *repository) DeleteFile(fileName string) error {
	defer getMeterFunc("DeleteFile")()

	return os.Remove(path.Join(repo.path, fileName))
}

func (repo *repository) AddAndCommit(msg string) (string, error) {
	defer getMeterFunc("AddAndCommit")()

	oid, err := repo.AddFilesToIndex([]string{"."})
	if err != nil {
		return "", fmt.Errorf("failed to add changes to index: %w", err)
	}
	sig := git.Signature{
		Name:  "Sturdy",
		Email: "support@getsturdy.com",
		When:  time.Now(),
	}
	commitID, err := repo.CommitIndexTree(oid, msg, sig)
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}
	return commitID, nil
}

func (repo *repository) CreateRootCommit() error {
	defer getMeterFunc("CreateRootCommit")()

	sig := &git.Signature{
		Name:  "Sturdy",
		Email: "support@getsturdy.com",
		When:  time.Now(),
	}
	index, err := repo.r.Index()
	if err != nil {
		return fmt.Errorf("failed to access vcs index: %w", err)
	}
	defer index.Free()

	treeId, err := index.WriteTree()
	if err != nil {
		return err
	}
	tree, err := repo.r.LookupTree(treeId)
	if err != nil {
		return err
	}
	err = index.Write()
	if err != nil {
		return err
	}
	_, err = repo.r.CreateCommit("HEAD", sig, sig, "Root Commit", tree)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}
	return nil
}

func (repo *repository) getHead() (*git.Tree, error) {
	defer getMeterFunc("getHead")()

	head, err := repo.r.RevparseSingle("HEAD^{tree}")
	if err != nil {
		return nil, fmt.Errorf("failed access repo HEAD: %w", err)
	}
	defer head.Free()
	headTree, err := repo.r.LookupTree(head.Id())
	if err != nil {
		return nil, fmt.Errorf("failed access repo HEAD: %w", err)
	}
	return headTree, nil
}

func (r *repository) RemoteBranchCommit(remoteName, branchName string) (*git.Commit, error) {
	defer getMeterFunc("remoteBranchCommit")()

	ref, err := r.r.References.Lookup("refs/remotes/" + remoteName + "/" + branchName)
	if err != nil {
		return nil, fmt.Errorf("failed to look up reference %s: %w", branchName, err)
	}
	defer ref.Free()

	commit, err := r.r.LookupCommit(ref.Branch().Target())
	if err != nil {
		return nil, fmt.Errorf("failed to look up commit: %w", err)
	}

	return commit, nil
}

func (r *repository) CreateBranchTrackingUpstream(branchName string) error {
	defer getMeterFunc("CreateBranchTrackingUpstream")()

	commit, err := r.RemoteBranchCommit("origin", branchName)
	if err != nil {
		return fmt.Errorf("failed to look up reference %s: %w", branchName, err)
	}
	defer commit.Free()

	branch, err := r.r.CreateBranch(branchName, commit, true)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	defer branch.Free()

	return nil
}

func (r *repository) CheckoutBranchWithForce(branchName string) error {
	defer getMeterFunc("CheckoutBranchWithForce")()
	return r.checkoutBranch(branchName, true)
}

func (r *repository) CheckoutBranchSafely(branchName string) error {
	defer getMeterFunc("CheckoutBranchSafely")()
	return r.checkoutBranch(branchName, false)
}

func (r *repository) checkoutBranch(branchName string, forceFully bool) error {
	branch, err := r.r.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return fmt.Errorf("failed to look up branch %s: %w", branchName, err)
	}
	defer branch.Free()

	commit, err := r.r.LookupCommit(branch.Target())
	if err != nil {
		return fmt.Errorf("failed to look up commit: %w", err)
	}
	defer commit.Free()

	tree, err := r.r.LookupTree(commit.TreeId())
	if err != nil {
		return fmt.Errorf("failed to look up tree: %w", err)
	}
	defer tree.Free()

	// TODO: CheckoutRemoveUntracked should be used even when CheckoutSafe but not currently possible
	// because during syncing a checkout is performed before commits are created
	strat := git.CheckoutSafe
	if forceFully {
		strat = git.CheckoutForce | git.CheckoutRemoveUntracked
	}

	err = r.r.CheckoutTree(tree, &git.CheckoutOpts{
		Strategy: strat,
	})
	if err != nil {
		return fmt.Errorf("failed to check out tree: %w", err)
	}
	newHead := "refs/heads/" + branchName
	err = r.r.SetHead(newHead)
	if err != nil {
		return fmt.Errorf("failed to set head to %s: %w", newHead, err)
	}
	return nil
}

func (r *repository) CreateAndCheckoutBranchAtCommit(commitID, newBranchName string) error {
	defer getMeterFunc("CreateAndCheckoutBranchAtCommit")()

	checkoutCmd := exec.Command("git", "checkout", "-f", commitID)
	checkoutCmd.Dir = r.path
	output, err := checkoutCmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		return fmt.Errorf("checkout failed: %w", err)
	}

	createBranchCmd := exec.Command("git", "checkout", "-b", newBranchName)
	createBranchCmd.Dir = r.path
	err = createBranchCmd.Run()
	if err != nil {
		return fmt.Errorf("checkout -b failed: %w", err)
	}

	return nil
}

func isGitNotFound(err error) bool {
	var gitError *git.GitError
	if errors.As(err, &gitError) {
		return gitError.Code == git.ErrorCodeNotFound
	}
	return false
}

func (r *repository) DeleteBranch(name string) error {
	defer getMeterFunc("DeleteBranch")()

	branch, err := r.r.LookupBranch(name, git.BranchAll)
	switch {
	case err == nil:
	case isGitNotFound(err):
		return nil
	default:
		return fmt.Errorf("failed to find branch %s: %w", name, err)
	}
	defer branch.Free()

	if err = branch.Delete(); err != nil {
		return fmt.Errorf("failed to delete branch %s: %w", name, err)
	}
	return nil
}

func (r *repository) CreateNewBranchOnHEAD(name string) error {
	defer getMeterFunc("CreateNewBranchOnHEAD")()

	branch, err := r.r.Head()
	if err != nil {
		return err
	}
	defer branch.Free()

	c, err := r.r.LookupCommit(branch.Target())
	if err != nil {
		return err
	}
	defer c.Free()

	newBranch, err := r.r.CreateBranch(name, c, true)
	if err != nil {
		return err
	}
	defer newBranch.Free()

	return nil
}

func (r *repository) CreateNewBranchAt(name string, targetSha string) error {
	defer getMeterFunc("CreateNewBranchAt")()

	id, err := git.NewOid(targetSha)
	if err != nil {
		return err
	}

	c, err := r.r.LookupCommit(id)
	if err != nil {
		return err
	}
	defer c.Free()

	newBranch, err := r.r.CreateBranch(name, c, true)
	if err != nil {
		return err
	}
	defer newBranch.Free()

	return nil
}

// CreateNewCommitBasedOnCommit creates a new commit using the tree and parent of existingCommitID.
func (r *repository) CreateNewCommitBasedOnCommit(newBranchName string, existingCommitID string, signature git.Signature, message string) (string, error) {
	defer getMeterFunc("CreateNewCommitBasedOnCommit")()

	// delete branch if already exists
	_ = r.DeleteBranch(newBranchName)

	id, err := git.NewOid(existingCommitID)
	if err != nil {
		return "", fmt.Errorf("failed to get existing commit: %w", err)
	}

	c, err := r.r.LookupCommit(id)
	if err != nil {
		return "", fmt.Errorf("failed to lookup existing commit: %w", err)
	}

	tree, err := c.Tree()
	if err != nil {
		return "", fmt.Errorf("failed to get tree: %w", err)
	}

	newCommit, err := r.r.CreateCommit("refs/heads/"+newBranchName, &signature, &signature, message, tree, c.Parent(0))
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}

	return newCommit.String(), nil
}

func (r *repository) Push(logger *zap.Logger, branchName string) error {
	defer getMeterFunc("Push")()

	_, err := r.pushNamedRemote(logger, branchName, "origin", branchName, nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) ForcePush(logger *zap.Logger, branchName string) error {
	defer getMeterFunc("ForcePush")()

	_, err := r.pushNamedRemote(logger, branchName, "origin", branchName, nil, true)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) pushNamedRemote(logger *zap.Logger, branchName, remoteName, remoteBranchName string, creds git.CredentialsCallback, force bool) (userError string, err error) {
	// https://git-scm.com/book/it/v2/Git-Internals-The-Refspec
	// "+" means to allow force pushes
	var maybeForcePrefix string
	if force {
		maybeForcePrefix = "+"
	}

	rs := maybeForcePrefix + "refs/heads/" + branchName + ":refs/heads/" + remoteBranchName

	remote, err := r.r.Remotes.Lookup(remoteName)
	if err != nil {
		return fmt.Sprintf("Push failed: %s", err.Error()), err
	}
	defer remote.Free()

	return r.pushRemoteWithRefSpec(logger, remote, creds, []string{rs})
}

func (r *repository) pushRemoteWithRefSpec(logger *zap.Logger, remote *git.Remote, creds git.CredentialsCallback, refspecs []string) (userError string, err error) {
	var githubProtectedBranchDeclined bool
	var githubProtectedBranchDeclinedRefName string

	opts := &git.PushOptions{}
	if creds != nil {
		opts.RemoteCallbacks = git.RemoteCallbacks{
			CredentialsCallback:      creds,
			CertificateCheckCallback: func(cert *git.Certificate, valid bool, hostname string) error { return nil },
		}
	}
	opts.RemoteCallbacks.PushUpdateReferenceCallback = func(refname, status string) error {
		if strings.Contains(status, "protected branch hook declined") {
			githubProtectedBranchDeclinedRefName = refname
			githubProtectedBranchDeclined = true
		}

		if status != "" {
			logger.Error("pushing to github failed", zap.String("status", status), zap.String("refname", refname))
		}

		return nil
	}
	err = remote.Push(refspecs, opts)
	if err != nil {
		return fmt.Sprintf("Push failed: %s", err.Error()),
			fmt.Errorf("failed to push refspecs=%v: %w", refspecs, err)
	}

	if githubProtectedBranchDeclined {
		// This error is propagated to the user
		return fmt.Sprintf("GitHub rejected the push to %s as it's protected by branch protection rules.", githubProtectedBranchDeclinedRefName), fmt.Errorf("failed to push to github: stopped by branch protection rules")
	}

	return "", nil
}

func (r *repository) annotatedCommitFromBranchName(name string) (*git.AnnotatedCommit, error) {
	branch, err := r.r.LookupBranch(name, git.BranchAll)
	if err != nil {
		return nil, err
	}
	defer branch.Free()

	branchRef, err := branch.Resolve()
	if err != nil {
		return nil, err
	}
	defer branchRef.Free()

	annotatedCommit, err := r.r.AnnotatedCommitFromRef(branchRef)
	if err != nil {
		return nil, err
	}
	return annotatedCommit, nil
}

func (r *repository) commitFromBranchName(name string) (*git.Commit, error) {
	branch, err := r.r.LookupBranch(name, git.BranchAll)
	if err != nil {
		return nil, err
	}
	defer branch.Free()

	branchRef, err := branch.Resolve()
	if err != nil {
		return nil, err
	}
	defer branchRef.Free()

	commit, err := r.r.LookupCommit(branchRef.Target())
	if err != nil {
		return nil, err
	}
	return commit, nil
}

// MoveBranch moves branchName to point to targetBranchName
func (r *repository) MoveBranch(branchName, targetBranchName string) error {
	defer getMeterFunc("MoveBranch")()

	targetCommit, err := r.commitFromBranchName(targetBranchName)
	if err != nil {
		return err
	}
	defer targetCommit.Free()

	branch, err := r.r.CreateBranch(branchName, targetCommit, true)
	if err != nil {
		return err
	}
	defer branch.Free()

	return nil
}

func (r *repository) MoveBranchToCommit(branchName, targetCommitSha string) error {
	defer getMeterFunc("MoveBranchToCommit")()

	commitOid, err := git.NewOid(targetCommitSha)
	if err != nil {
		return err
	}
	targetCommit, err := r.r.LookupCommit(commitOid)
	if err != nil {
		return err
	}
	defer targetCommit.Free()

	branch, err := r.r.CreateBranch(branchName, targetCommit, true)
	if err != nil {
		return err
	}
	defer branch.Free()

	return nil
}

func (r *repository) MoveBranchToHEAD(branchName string) error {
	defer getMeterFunc("MoveBranchToHEAD")()

	targetCommit, err := r.HeadCommit()
	if err != nil {
		return err
	}
	defer targetCommit.Free()

	branch, err := r.r.CreateBranch(branchName, targetCommit, true)
	if err != nil {
		return err
	}
	defer branch.Free()

	return nil
}

func (r *repository) FetchBranch(branches ...string) error {
	defer getMeterFunc("FetchBranch")()
	return r.fetch(branches...)
}

func (r *repository) fetch(branches ...string) error {
	remotes, err := r.r.Remotes.List()
	if err != nil {
		return err
	}

	for _, remoteName := range remotes {
		remote, err := r.r.Remotes.Lookup(remoteName)
		if err != nil {
			return err
		}
		defer remote.Free()

		var refspecs []string
		for _, branch := range branches {
			refspecs = append(refspecs, fmt.Sprintf("+refs/heads/%s:refs/remotes/%s/%s", branch, remoteName, branch))
		}

		err = remote.Fetch(refspecs, &git.FetchOptions{
			DownloadTags:    git.DownloadTagsNone,
			UpdateFetchhead: true,
		}, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) listBranches() ([]string, error) {
	bi, err := r.r.NewBranchIterator(git.BranchAll)
	if err != nil {
		return nil, err
	}
	defer bi.Free()

	var branches []string
	err = bi.ForEach(func(branch *git.Branch, bt git.BranchType) error {
		name, err := branch.Name()
		if err != nil {
			return err
		}
		branches = append(branches, name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *repository) HeadBranch() (string, error) {
	defer getMeterFunc("HeadBranch")()
	head, err := r.r.Head()
	if err != nil {
		return "", err
	}
	defer head.Free()

	name, err := head.Branch().Name()
	if err != nil {
		return "", err
	}

	return name, nil
}

func (r *repository) SetDefaultBranch(targetBranch string) error {
	defer getMeterFunc("SetDefaultBranch")()
	ref, err := r.r.References.Lookup("HEAD")
	if err != nil {
		return err
	}
	defer ref.Free()

	// git symbolic-ref HEAD refs/heads/sturdytrunk
	newRef, err := ref.SetSymbolicTarget("refs/heads/"+targetBranch, "")
	if err != nil {
		return err
	}
	defer newRef.Free()

	return nil
}

func (r *repository) GetDefaultBranch() (string, error) {
	defer getMeterFunc("GetDefaultBranch")()
	ref, err := r.r.References.Lookup("HEAD")
	if err != nil {
		return "", err
	}
	defer ref.Free()

	return ref.SymbolicTarget(), nil
}

func (r *repository) CreateAndSetDefaultBranch(headBranchName string) error {
	defer getMeterFunc("CreateAndSetDefaultBranch")()

	if err := r.CreateNewBranchOnHEAD(headBranchName); err != nil {
		return fmt.Errorf("failed to create a new trunk: %w", err)
	}

	if err := r.SetDefaultBranch(headBranchName); err != nil {
		return fmt.Errorf("failed to set new trunk branch: %w", err)
	}

	return nil
}

var ErrNotFound = errors.New("not found")

func (r *repository) HeadCommit() (*git.Commit, error) {
	defer getMeterFunc("HeadCommit")()
	ref, err := r.r.Head()
	if err != nil {
		//nolint:errorlint
		if gErr, ok := err.(*git.GitError); ok && gErr.Code == git.ErrorCodeUnbornBranch {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer ref.Free()

	commit, err := r.r.LookupCommit(ref.Target())
	if err != nil {
		return nil, err
	}

	return commit, nil
}

// CleanStaged does the equivalent of "git restore --staged ."
// Which is the same as "git reset --mixed HEAD"
func (r *repository) CleanStaged() error {
	defer getMeterFunc("CleanStaged")()
	commit, err := r.HeadCommit()
	if err != nil {
		return err
	}
	defer commit.Free()

	if err := r.ResetMixed(commit.Id().String()); err != nil {
		return err
	}
	return nil
}
