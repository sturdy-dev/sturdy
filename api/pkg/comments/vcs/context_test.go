package vcs

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/internal/testmodule"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	type deps struct {
		di.In
		GitSnapshotter   snapshotter.Snapshotter
		WorkspaceRepo    db_workspaces.Repository
		SnapshotRepo     db_snapshots.Repository
		ViewRepo         db_view.Repository
		ExecutorProvider executor.Provider
		RepoProvider     provider.RepoProvider
	}
	var d deps
	if err := di.Init(&d, testmodule.TestModule); err != nil {
		t.Fatal(err)
	}

	codebaseID := uuid.NewString()
	workspaceID := uuid.NewString()
	viewID := uuid.NewString()

	trunkPath := d.RepoProvider.TrunkPath(codebaseID)
	viewPath := d.RepoProvider.ViewPath(codebaseID, viewID)

	_, err := vcs.CreateBareRepoWithRootCommit(trunkPath)
	if err != nil {
		panic(err)
	}
	viewGitRepo, err := vcs.CloneRepo(trunkPath, viewPath)
	if err != nil {
		panic(err)
	}

	err = viewGitRepo.CreateNewBranchOnHEAD(workspaceID)
	assert.NoError(t, err)
	err = viewGitRepo.CheckoutBranchWithForce(workspaceID)
	assert.NoError(t, err)

	ws := &workspaces.Workspace{
		ID:         workspaceID,
		CodebaseID: codebaseID,
		ViewID:     &viewID,
	}
	vw := &view.View{
		ID:          viewID,
		CodebaseID:  codebaseID,
		WorkspaceID: workspaceID,
	}

	err = d.WorkspaceRepo.Create(*ws)
	assert.NoError(t, err)
	err = d.ViewRepo.Create(*vw)
	assert.NoError(t, err)

	// write file with line numbers
	var content strings.Builder
	for i := 1; i < 15; i++ {
		content.WriteString(fmt.Sprintf("%d\n", i))
	}
	err = ioutil.WriteFile(viewPath+"/file.txt", []byte(content.String()), 0o666)
	assert.NoError(t, err)

	_, err = viewGitRepo.AddAndCommit("adding")
	assert.NoError(t, err)

	// replace with line numbers +100
	content = strings.Builder{}
	for i := 1; i < 15; i++ {
		content.WriteString(fmt.Sprintf("%d\n", i+100))
	}
	err = ioutil.WriteFile(viewPath+"/file.txt", []byte(content.String()), 0o666)
	assert.NoError(t, err)

	cases := []struct {
		lineNumber              int
		lineIsNew               bool
		useSnapshot             bool
		expected                string
		expectedContextStartsAt int
	}{
		{
			lineNumber:              8,
			lineIsNew:               false,
			expected:                "6\n7\n8\n9\n10\n",
			expectedContextStartsAt: 6,
		},
		{
			lineNumber:              1,
			lineIsNew:               false,
			expected:                "1\n2\n3\n",
			expectedContextStartsAt: 1,
		},
		{
			lineNumber:              14,
			lineIsNew:               false,
			expected:                "12\n13\n14\n",
			expectedContextStartsAt: 12,
		},
		{
			lineNumber:              50,
			lineIsNew:               false,
			expected:                "",
			expectedContextStartsAt: -1,
		},
		{
			lineNumber:              -400,
			lineIsNew:               false,
			expected:                "",
			expectedContextStartsAt: -1,
		},

		// new version
		{
			lineNumber:              8,
			lineIsNew:               true,
			expected:                "106\n107\n108\n109\n110\n",
			expectedContextStartsAt: 6,
		},
		{
			lineNumber:              1,
			lineIsNew:               true,
			expected:                "101\n102\n103\n",
			expectedContextStartsAt: 1,
		},
		{
			lineNumber:              14,
			lineIsNew:               true,
			expected:                "112\n113\n114\n",
			expectedContextStartsAt: 12,
		},
		{
			lineNumber:              50,
			lineIsNew:               true,
			expected:                "",
			expectedContextStartsAt: -1,
		},
		{
			lineNumber:              -400,
			lineIsNew:               true,
			expected:                "",
			expectedContextStartsAt: -1,
		},

		// test with snapshot
		{
			lineNumber:              8,
			lineIsNew:               false,
			useSnapshot:             true,
			expected:                "6\n7\n8\n9\n10\n",
			expectedContextStartsAt: 6,
		},
		{
			lineNumber:              14,
			lineIsNew:               true,
			useSnapshot:             true,
			expected:                "112\n113\n114\n",
			expectedContextStartsAt: 12,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%d-%v", tc.lineNumber, tc.lineIsNew), func(t *testing.T) {
			if tc.useSnapshot {
				// make a snapshot
				snapshot, err := d.GitSnapshotter.Snapshot(codebaseID, workspaceID, snapshots.ActionViewSync, snapshotter.WithOnView(viewID), snapshotter.WithNoThrottle())
				if assert.NoError(t, err) && assert.NotNil(t, snapshot) {
					ws.LatestSnapshotID = &snapshot.ID
				}
			} else {
				ws.LatestSnapshotID = nil
			}

			res, contextStartsAt, err := GetWorkspaceContext(tc.lineNumber, tc.lineIsNew, "file.txt", nil, ws, d.ExecutorProvider, d.SnapshotRepo)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.expectedContextStartsAt, contextStartsAt)
		})
	}
}
