package vcs

import (
	"errors"
	"fmt"
	"strings"

	git "github.com/libgit2/git2go/v33"
)

var ErrFileNotFound = errors.New("file not found in tree")

func (r *repository) FileContentsAtCommit(commitID, filePath string) ([]byte, error) {
	blob, err := r.FileBlobAtCommit(commitID, filePath)
	if err != nil {
		return nil, err
	}
	return blob.Contents(), nil
}

func (r *repository) FileBlobAtCommit(commitID, filePath string) (*git.Blob, error) {
	defer getMeterFunc("FileBlobAtCommit")()
	oid, err := git.NewOid(commitID)
	if err != nil {
		return nil, err
	}

	commit, err := r.r.LookupCommit(oid)
	if err != nil {
		return nil, err
	}
	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	defer tree.Free()

	entry, err := tree.EntryByPath(filePath)
	if err != nil {
		return nil, ErrFileNotFound
	}

	blob, err := r.r.LookupBlob(entry.Id)
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func (r *repository) DirectoryChildrenAtCommit(commitID, directoryPath string) ([]string, error) {
	defer getMeterFunc("DirectoryChildrenAtCommit")()
	oid, err := git.NewOid(commitID)
	if err != nil {
		return nil, err
	}

	commit, err := r.r.LookupCommit(oid)
	if err != nil {
		return nil, err
	}
	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	defer tree.Free()

	trimmedPath := strings.Trim(directoryPath, "/")
	var prefixPath string
	var subTree *git.Tree

	if trimmedPath == "" {
		subTree = tree
	} else {
		prefixPath = trimmedPath + "/"
		entry, err := tree.EntryByPath(trimmedPath)
		if err != nil {
			return nil, err
		}

		if entry.Type != git.ObjectTree {
			return nil, fmt.Errorf("path doesn't represent a directory: %v", directoryPath)
		}

		subTree, err = r.r.LookupTree(entry.Id)
		if err != nil {
			return nil, err
		}
		defer subTree.Free()
	}

	var entries []string
	for i := uint64(0); i < subTree.EntryCount(); i++ {
		entries = append(entries, prefixPath+subTree.EntryByIndex(i).Name)
	}

	return entries, nil
}
