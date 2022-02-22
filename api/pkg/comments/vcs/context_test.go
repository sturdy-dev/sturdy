package vcs

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"io/ioutil"

	events2 "getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/internal/inmemory"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/view"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestContext(t *testing.T) {
	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)

	repoProvider := provider.New(tmpBase, "")

	codebaseID := uuid.NewString()
	workspaceID := uuid.NewString()
	viewID := uuid.NewString()

	trunkPath := repoProvider.TrunkPath(codebaseID)
	viewPath := repoProvider.ViewPath(codebaseID, viewID)

	_, err = vcs.CreateBareRepoWithRootCommit(trunkPath)
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

			snapshotRepo := inmemory.NewInMemorySnapshotRepo()
			workspaceRepo := db_workspaces.NewMemory()
			viewRepo := inmemory.NewInMemoryViewRepo()
			events := events2.NewInMemory()
			logger := zap.NewNop()
			executorProvider := executor.NewProvider(logger, repoProvider)
			codebaseUserRepo := inmemory.NewInMemoryCodebaseUserRepo()
			eventsSender := events2.NewSender(codebaseUserRepo, workspaceRepo, events)
			gitSnapshotter := snapshotter.NewGitSnapshotter(snapshotRepo, workspaceRepo, workspaceRepo, viewRepo, eventsSender, executorProvider, logger)

			err = workspaceRepo.Create(*ws)
			assert.NoError(t, err)
			err = viewRepo.Create(*vw)
			assert.NoError(t, err)

			if tc.useSnapshot {
				// make a snapshot
				snapshot, err := gitSnapshotter.Snapshot(codebaseID, workspaceID, snapshots.ActionViewSync, snapshotter.WithOnView(viewID))
				assert.NoError(t, err)
				ws.LatestSnapshotID = &snapshot.ID
			}

			res, contextStartsAt, err := GetContext(tc.lineNumber, tc.lineIsNew, "file.txt", ws, executorProvider, snapshotRepo)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, res)
			assert.Equal(t, tc.expectedContextStartsAt, contextStartsAt)
		})
	}
}
