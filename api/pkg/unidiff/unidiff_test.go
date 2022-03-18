package unidiff

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInvert(t *testing.T) {
	testCases := []struct {
		fileName         string
		expectedFileName string
	}{
		{"sample_rename.diff", "sample_rename_inverted.diff"},
		{"sample_new.diff", "sample_new_inverted.diff"},
		{"sample_deleted.diff", "sample_deleted_inverted.diff"},
		{"sample_changed.diff", "sample_changed_inverted.diff"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			data, err := ioutil.ReadFile("testdata/" + tc.fileName)
			assert.NoError(t, err)

			expected, err := ioutil.ReadFile("testdata/" + tc.expectedFileName)
			assert.NoError(t, err)

			inverted, err := NewUnidiff(NewStringsPatchReader([]string{string(data)}), zap.NewNop()).WithInverter().Decorate()
			assert.NoError(t, err)

			assert.Equal(t, string(expected), inverted[0].Hunks[0].Patch)
		})
	}
}

func TestExpanded(t *testing.T) {
	testCases := []struct {
		fileName           string
		expectedsFileNames []string
	}{
		{"sample_two_hunks.diff", []string{"sample_two_hunks_expected0.diff", "sample_two_hunks_expected1.diff"}},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			data, err := ioutil.ReadFile("testdata/" + tc.fileName)
			assert.NoError(t, err)

			var expecteds []string
			for _, efn := range tc.expectedsFileNames {
				e, err := ioutil.ReadFile("testdata/" + efn)
				assert.NoError(t, err)
				expecteds = append(expecteds, string(e))
			}

			expanded, err := NewUnidiff(NewStringsPatchReader([]string{string(data)}), zap.NewNop()).WithExpandedHunks().Patches()
			assert.NoError(t, err)
			assert.Equal(t, expecteds, expanded)
		})
	}
}

func TestDecorate(t *testing.T) {
	testCases := []struct {
		inputFileName string
		ignoreBinary  bool
		withExpand    bool
		expected      []FileDiff
	}{
		{
			inputFileName: "sample_two_hunks.diff",
			expected: []FileDiff{
				{
					OrigName:      "one.txt",
					NewName:       "one.txt",
					PreferredName: "one.txt",
					Hunks: []Hunk{{
						ID:    "d27a3720c448525c4bdac6149b4344f3e09617f2bffe3bad6ddaaa0e86ebff5e",
						Patch: "diff --git \"a/one.txt\" \"b/one.txt\"\nindex 4fce4a5..fef85d8 100644\n--- \"a/one.txt\"\n+++ \"b/one.txt\"\n@@ -2,7 +2,6 @@ a\n b\n c\n d\n-e\n f\n g\n h\n@@ -16,7 +15,6 @@ o\n p\n q\n r\n-s\n t\n y\n v\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_two_hunks.diff",
			ignoreBinary:  true,
			expected: []FileDiff{
				{
					OrigName:      "one.txt",
					NewName:       "one.txt",
					PreferredName: "one.txt",
					Hunks: []Hunk{{
						ID:    "d27a3720c448525c4bdac6149b4344f3e09617f2bffe3bad6ddaaa0e86ebff5e",
						Patch: "diff --git \"a/one.txt\" \"b/one.txt\"\nindex 4fce4a5..fef85d8 100644\n--- \"a/one.txt\"\n+++ \"b/one.txt\"\n@@ -2,7 +2,6 @@ a\n b\n c\n d\n-e\n f\n g\n h\n@@ -16,7 +15,6 @@ o\n p\n q\n r\n-s\n t\n y\n v\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_deleted.diff",
			expected: []FileDiff{
				{
					OrigName:      "bar",
					NewName:       "/dev/null",
					PreferredName: "bar",
					IsDeleted:     true,
					Hunks: []Hunk{{
						ID:    "5b938598d85a4f0fbe6cf44dd864c82e9374a9ef1569eb013be2239393896005",
						Patch: "diff --git \"a/bar\" /dev/null\ndeleted file mode 100644\nindex a1f8944..0000000\n--- \"a/bar\"\n+++ /dev/null\n@@ -1,4 +0,0 @@\n-foo\n-foo\n-foo\n-foo\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_new.diff",
			expected: []FileDiff{
				{
					OrigName:      "/dev/null",
					NewName:       "README_XOXO.md",
					PreferredName: "README_XOXO.md",
					IsNew:         true,
					Hunks: []Hunk{{
						ID:    "be3e999e431c5f955a98f6d25e27ba8329d5c08f9d03a6f13d057c973d8e0d0a",
						Patch: "diff --git /dev/null \"b/README_XOXO.md\"\nnew file mode 100644\nindex 0000000..bc56c4d\n--- /dev/null\n+++ \"b/README_XOXO.md\"\n@@ -0,0 +1,1 @@\n+Foo\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_file_extended_empty_new.diff",
			expected: []FileDiff{
				{
					OrigName:      "/dev/null",
					NewName:       "vendor/go/build/testdata/empty/dummy",
					PreferredName: "vendor/go/build/testdata/empty/dummy",
					IsNew:         true,
					Hunks: []Hunk{{
						ID:    "5b8042925dd40be6b84660435a6953b4290290189040805d6b8de83b1d441d13",
						Patch: "diff --git /dev/null \"b/vendor/go/build/testdata/empty/dummy\"\nnew file mode 100644\nindex 0000000..e69de29\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_file_extended_empty_rename.diff",
			expected: []FileDiff{
				{
					OrigName:      "docs/integrations/Email_Notifications.md",
					NewName:       "docs/integrations/email-notifications.md",
					PreferredName: "docs/integrations/email-notifications.md",
					IsMoved:       true,
					Hunks: []Hunk{{
						ID:    "0aec6557aa614110301ac44c0b17ea3d44ff161f9fea0766999d7a63ee08e737",
						Patch: "diff --git \"a/docs/integrations/Email_Notifications.md\" \"b/docs/integrations/email-notifications.md\"\nsimilarity index 100%\nrename from \"docs/integrations/Email_Notifications.md\"\nrename to \"docs/integrations/email-notifications.md\"\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_file_extended_empty_deleted.diff",
			expected: []FileDiff{
				{
					OrigName:      "vendor/go/build/testdata/empty/dummy",
					NewName:       "/dev/null",
					PreferredName: "vendor/go/build/testdata/empty/dummy",
					IsDeleted:     true,
					Hunks: []Hunk{{
						ID:    "b58e43d2a21696bb6d7f5275809705e0e60ac965675a5352d528cf913d3d490a",
						Patch: "diff --git \"a/vendor/go/build/testdata/empty/dummy\" /dev/null\ndeleted file mode 100644\nindex e69de29..0000000\n",
					}},
				},
			},
		},

		{
			inputFileName: "sample_binary_new_space_in_name.diff",
			expected: []FileDiff{
				{
					OrigName:      "/dev/null",
					NewName:       "new binary with space.txt",
					PreferredName: "new binary with space.txt",
					IsNew:         true,
					Hunks: []Hunk{{
						ID:    "dcf21c66d2242dcfe253beb3870ab376fadeaaa5ad34fd4405987f12bbb305e4",
						Patch: "diff --git /dev/null \"b/new binary with space.txt\"\nnew file mode 100644\nindex 0000000..593f470\nBinary files /dev/null and b/new binary with space.txt differ\n",
					}},
				},
			},
		},

		{
			// this diff is extended, but not a binary change
			inputFileName: "sample_file_extended_empty_deleted.diff",
			ignoreBinary:  true,
			expected: []FileDiff{
				{
					OrigName:      "vendor/go/build/testdata/empty/dummy",
					NewName:       "/dev/null",
					PreferredName: "vendor/go/build/testdata/empty/dummy",
					IsDeleted:     true,
					Hunks: []Hunk{{
						ID:    "b58e43d2a21696bb6d7f5275809705e0e60ac965675a5352d528cf913d3d490a",
						Patch: "diff --git \"a/vendor/go/build/testdata/empty/dummy\" /dev/null\ndeleted file mode 100644\nindex e69de29..0000000\n",
					}},
				},
			},
		},

		{
			// binary diff is ignored
			inputFileName: "sample_binary_new_space_in_name.diff",
			ignoreBinary:  true,
			expected:      nil,
		},

		{
			// binary diff is ignored
			inputFileName: "sample_binary_differs.diff",
			ignoreBinary:  true,
			expected:      nil,
		},

		{
			inputFileName: "sample_binary_differs.diff",
			expected: []FileDiff{
				{
					OrigName:      "app/assets/bin/sturdy",
					NewName:       "app/assets/bin/sturdy",
					PreferredName: "app/assets/bin/sturdy",
					Hunks: []Hunk{{
						ID: "822db9ff4a8c48569f74dbc1c63b756539fb3ea40e671f7b555222cd02128c25",
						// Preserved full diff
						Patch: "diff --git \"a/app/assets/bin/sturdy\" \"b/app/assets/bin/sturdy\"\nold mode 100755\nnew mode 100644\nindex 16edd4f..9b8fb76\nBinary files a/app/assets/bin/sturdy and b/app/assets/bin/sturdy differ\n",
					}},
				},
			},
		},

		{
			// binary diff is ignored
			inputFileName: "sample_hunked_additions.diff",
			withExpand:    true,
			expected: []FileDiff{{
				OrigName:      "pre.txt",
				NewName:       "post.txt",
				PreferredName: "post.txt",
				IsDeleted:     false,
				IsNew:         false,
				IsMoved:       true,
				Hunks: []Hunk{
					{
						ID:         "f134b2f129004e22e9006334d33eed8417e97fe4731fb140da7b157af1a2e05b",
						Patch:      "diff --git \"a/pre.txt\" \"b/post.txt\"\nindex 7904388..0f424bb 100644\n--- \"a/pre.txt\"\n+++ \"b/post.txt\"\n@@ -8,6 +8,11 @@ b\n b\n b\n b\n+1\n+1\n+1\n+1\n+1\n c\n c\n c\n",
						IsOutdated: false,
						IsApplied:  false,
					},
					{
						ID:         "da6038761751242b290a69e5eb25f917c1ce010dbad086374d92f5da4900c71a",
						Patch:      "diff --git \"a/pre.txt\" \"b/post.txt\"\nindex 7904388..0f424bb 100644\n--- \"a/pre.txt\"\n+++ \"b/post.txt\"\n@@ -38,6 +43,14 @@ g\n g\n g\n g\n+2\n+2\n+2\n+2\n+2\n+2\n+2\n+2\n g\n g\n g\n",
						IsOutdated: false,
						IsApplied:  false,
					},
				},
			}},
		},
	}

	logger, _ := zap.NewDevelopment()

	for _, tc := range testCases {
		t.Run(tc.inputFileName, func(t *testing.T) {
			input, err := ioutil.ReadFile("testdata/" + tc.inputFileName)
			assert.NoError(t, err)

			d := NewUnidiff(NewBytesPatchReader([][]byte{input}), logger)

			if tc.ignoreBinary {
				d = d.WithIgnoreBinary()
			}
			if tc.withExpand {
				d = d.WithExpandedHunks()
			}

			res, err := d.Decorate()
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestFilterAndApply(t *testing.T) {
	input, err := ioutil.ReadFile("testdata/500.diff")
	assert.NoError(t, err)

	diffs, err := NewUnidiff(NewBytesPatchReader([][]byte{input}), zap.NewNop()).WithExpandedHunks().Decorate()
	assert.NoError(t, err)

	expectedDiffs := []FileDiff{
		{
			OrigName: "500.txt", NewName: "500.txt", PreferredName: "500.txt",
			Hunks: []Hunk{
				{ID: "5784fab89b347ef201e571311cd5c989d766cd3ed37d0b0f466f4ccf0831cd3e", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -10,6 +10,13 @@\n 10\n 11\n 12\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 13\n 14\n 15\n", IsOutdated: false, IsApplied: false},
				{ID: "31754de92b950e213380f3fd4b907d59c9dd7e9ff182176c319d9429167227ca", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -63,30 +70,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n", IsOutdated: false, IsApplied: false},
				{ID: "2c9fde653eef9c9dec3fca3571509b7826d7fc1aecfba9c9217608589bfb1005", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -218,26 +221,6 @@\n 218\n 219\n 220\n-221\n-222\n-223\n-224\n-225\n-226\n-227\n-228\n-229\n-230\n-231\n-232\n-233\n-234\n-235\n-236\n-237\n-238\n-239\n-240\n 241\n 242\n 243\n", IsOutdated: false, IsApplied: false},
				{ID: "e95665a4eff9d8aa42ca27a9c138ef84ff6aa1a48a0961d8c0784ecc6db1ef1d", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -309,6 +292,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false},
			},
		},
	}
	assert.Equal(t, expectedDiffs, diffs)

	tmpBase, err := ioutil.TempDir(os.TempDir(), "mash")
	assert.NoError(t, err)
	repoProvider := provider.New(tmpBase, "")
	codebaseID := codebases.ID(uuid.NewString())
	trunkPath := repoProvider.TrunkPath(codebaseID)
	_, err = vcs.CreateBareRepoWithRootCommit(trunkPath)
	assert.NoError(t, err)

	// Apply filters, and test that the patches can be applied after filtering
	permutTestCases := []struct {
		hunkIndexes []int
		withInvert  bool
		withJoiner  bool
		expected    []Hunk
	}{
		// Singles
		{hunkIndexes: []int{0}, expected: []Hunk{expectedDiffs[0].Hunks[0]}},
		{hunkIndexes: []int{1}, expected: []Hunk{{ID: "407fed518cd123cfebca67055f5b1b90b960efdf4305ca3680225f6aca83d614", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -63,30 +63,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{2}, expected: []Hunk{{ID: "28126836ce917eb057bd813ffcd3ef098e439cca1903566168648ebec1d188b1", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -218,26 +218,6 @@\n 218\n 219\n 220\n-221\n-222\n-223\n-224\n-225\n-226\n-227\n-228\n-229\n-230\n-231\n-232\n-233\n-234\n-235\n-236\n-237\n-238\n-239\n-240\n 241\n 242\n 243\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{3}, expected: []Hunk{{ID: "99cf82a826199fc69e0a4ca6944c7dffe0ff6d3905804a96177d62f8853eef32", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -309,6 +309,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false}}},

		// Doubles
		{hunkIndexes: []int{0, 1}, expected: []Hunk{{ID: "5784fab89b347ef201e571311cd5c989d766cd3ed37d0b0f466f4ccf0831cd3e", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -10,6 +10,13 @@\n 10\n 11\n 12\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 13\n 14\n 15\n", IsOutdated: false, IsApplied: false}, {ID: "31754de92b950e213380f3fd4b907d59c9dd7e9ff182176c319d9429167227ca", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -63,30 +70,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{0, 2}, expected: []Hunk{{ID: "5784fab89b347ef201e571311cd5c989d766cd3ed37d0b0f466f4ccf0831cd3e", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -10,6 +10,13 @@\n 10\n 11\n 12\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 13\n 14\n 15\n", IsOutdated: false, IsApplied: false}, {ID: "19665cc91aff1468b8e9a9d98515d3006b3f164e077b65e256c1527c6f7d396c", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -218,26 +225,6 @@\n 218\n 219\n 220\n-221\n-222\n-223\n-224\n-225\n-226\n-227\n-228\n-229\n-230\n-231\n-232\n-233\n-234\n-235\n-236\n-237\n-238\n-239\n-240\n 241\n 242\n 243\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{0, 3}, expected: []Hunk{{ID: "5784fab89b347ef201e571311cd5c989d766cd3ed37d0b0f466f4ccf0831cd3e", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -10,6 +10,13 @@\n 10\n 11\n 12\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 13\n 14\n 15\n", IsOutdated: false, IsApplied: false}, {ID: "8279caa2e83b72d51129e107d2758e8f58efeebdcba761cc294aceff726aab6e", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -309,6 +316,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{1, 2}, expected: []Hunk{{ID: "407fed518cd123cfebca67055f5b1b90b960efdf4305ca3680225f6aca83d614", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -63,30 +63,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n", IsOutdated: false, IsApplied: false}, {ID: "964e3864594a6a49305a1f9c017df5582f2bbeaa0ea53983cd5dd5930b7f6f50", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -218,26 +214,6 @@\n 218\n 219\n 220\n-221\n-222\n-223\n-224\n-225\n-226\n-227\n-228\n-229\n-230\n-231\n-232\n-233\n-234\n-235\n-236\n-237\n-238\n-239\n-240\n 241\n 242\n 243\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{1, 3}, expected: []Hunk{{ID: "407fed518cd123cfebca67055f5b1b90b960efdf4305ca3680225f6aca83d614", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -63,30 +63,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n", IsOutdated: false, IsApplied: false}, {ID: "66069faf8bc0e8ed0318bf35da2a405cc1a9c0ad003ee763723a1590bdd681c7", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -309,6 +305,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false}}},
		{hunkIndexes: []int{2, 3}, expected: []Hunk{{ID: "28126836ce917eb057bd813ffcd3ef098e439cca1903566168648ebec1d188b1", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -218,26 +218,6 @@\n 218\n 219\n 220\n-221\n-222\n-223\n-224\n-225\n-226\n-227\n-228\n-229\n-230\n-231\n-232\n-233\n-234\n-235\n-236\n-237\n-238\n-239\n-240\n 241\n 242\n 243\n", IsOutdated: false, IsApplied: false}, {ID: "78332236fc109ddc75b743a3d4844ae29e2e19a22cfa6bfb3d6acb940fd7cd8c", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -309,6 +289,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false}}},

		// Inverted single
		{hunkIndexes: []int{3}, withInvert: true, expected: []Hunk{{ID: "7b807fdc6c7219e4b4869b4ab6655914f4504a62ea56dc98be2bbb5f589ff9f8", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 23c60d6..3f1dcfc 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -292,13 +292,6 @@\n 309\n 310\n 311\n-added\n-added\n-added\n-added\n-added\n-added\n-added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false}}},

		// Inverted double
		{hunkIndexes: []int{1, 3}, withInvert: true, expected: []Hunk{
			{ID: "ce68869a05069a97c4cb60904b83035bea9fc00e10992fd82309a1e86cee791a", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 23c60d6..3f1dcfc 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -70,26 +70,30 @@\n 63\n 64\n 65\n+66\n+67\n+68\n+69\n+70\n+71\n+72\n+73\n+74\n-66modded\n-67modded\n-68modded\n-69modded\n-70modded\n-71modded\n-72modded\n-73modded\n-74modded\n 75\n 76\n 77\n-added\n-added\n-added\n-added\n-added\n 78\n 79\n 80\n+81\n+82\n+83\n+84\n+85\n+86\n+87\n+88\n+89\n 90\n 91\n 92\n", IsOutdated: false, IsApplied: false},
			{ID: "8b977697576ffe6202918f110c926a3d3f50547472aa7cd76c0466a820cea97b", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 23c60d6..3f1dcfc 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -292,13 +296,6 @@\n 309\n 310\n 311\n-added\n-added\n-added\n-added\n-added\n-added\n-added\n 312\n 313\n 314\n", IsOutdated: false, IsApplied: false},
		}},

		// Joined double
		{hunkIndexes: []int{1, 3}, withJoiner: true, expected: []Hunk{
			{ID: "244bc7e6899c6e21c781a9bf70a7b1fd45510b403d2436b0a05e882b95dba30f", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -63,30 +63,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n@@ -309,6 +305,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n"},
		}},

		// Joined quad
		{hunkIndexes: []int{0, 1, 2, 3}, withJoiner: true, expected: []Hunk{
			{ID: "6c3e1d4d1bf5de0f161d918178ccdeb05fedbaa9ef884051c50162583d6506e1", Patch: "diff --git \"a/500.txt\" \"b/500.txt\"\nindex 3f1dcfc..23c60d6 100644\n--- \"a/500.txt\"\n+++ \"b/500.txt\"\n@@ -10,6 +10,13 @@\n 10\n 11\n 12\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 13\n 14\n 15\n@@ -63,30 +70,26 @@\n 63\n 64\n 65\n-66\n-67\n-68\n-69\n-70\n-71\n-72\n-73\n-74\n+66modded\n+67modded\n+68modded\n+69modded\n+70modded\n+71modded\n+72modded\n+73modded\n+74modded\n 75\n 76\n 77\n+added\n+added\n+added\n+added\n+added\n 78\n 79\n 80\n-81\n-82\n-83\n-84\n-85\n-86\n-87\n-88\n-89\n 90\n 91\n 92\n@@ -218,26 +221,6 @@\n 218\n 219\n 220\n-221\n-222\n-223\n-224\n-225\n-226\n-227\n-228\n-229\n-230\n-231\n-232\n-233\n-234\n-235\n-236\n-237\n-238\n-239\n-240\n 241\n 242\n 243\n@@ -309,6 +292,13 @@\n 309\n 310\n 311\n+added\n+added\n+added\n+added\n+added\n+added\n+added\n 312\n 313\n 314\n"},
		}},
	}

	for _, tc := range permutTestCases {
		t.Run(fmt.Sprintf("%+v withInvert=%v", tc.hunkIndexes, tc.withInvert), func(t *testing.T) {
			var filter []string
			for _, i := range tc.hunkIndexes {
				filter = append(filter, diffs[0].Hunks[i].ID)
			}

			d := NewUnidiff(NewBytesPatchReader([][]byte{input}), zap.NewNop()).WithExpandedHunks().WithHunksFilter(filter...)
			if tc.withInvert {
				d = d.WithInverter()
			}
			if tc.withJoiner {
				d = d.WithJoiner()
			}

			filteredOne, err := d.Decorate()
			assert.NoError(t, err)
			filteredPatch := filteredOne[0]
			assert.Equal(t, tc.expected, filteredPatch.Hunks)

			// Test that the patch can be applied
			viewID := uuid.NewString()
			viewPath := repoProvider.ViewPath(codebaseID, viewID)
			t.Logf("viewPath=%s", viewPath)
			viewRepo, err := vcs.CloneRepo(trunkPath, viewPath)
			assert.NoError(t, err)

			// Reset the file to the "pre" (normal apply) or "post" (invert) state
			var preContents []byte
			if tc.withInvert {
				preContents, err = ioutil.ReadFile("testdata/500-post.txt")
			} else {
				preContents, err = ioutil.ReadFile("testdata/500-pre.txt")
			}
			assert.NoError(t, err)

			err = ioutil.WriteFile(path.Join(viewPath, "500.txt"), preContents, 0o666)
			assert.NoError(t, err)

			// Commit as is
			_, err = viewRepo.AddAndCommit("pre / post")
			assert.NoError(t, err)

			var allPatches [][]byte
			for _, h := range filteredPatch.Hunks {
				allPatches = append(allPatches, []byte(h.Patch))

				// The filtered hunk patch should apply
				canApply, err := viewRepo.CanApplyPatch([]byte(h.Patch))
				assert.True(t, canApply)
				assert.NoError(t, err)
			}

			// Test that all patches can be applied together
			_, err = viewRepo.ApplyPatchesToIndex(allPatches)
			assert.NoError(t, err)
		})
	}
}

func TestFilter(t *testing.T) {
	inputFiles := []string{
		"sample_changed.diff",
		"sample_binary_new_space_in_name.diff",
		"sample_file_extended_empty_new.diff",
		"500.diff",
	}

	var allPatches [][]byte
	for _, fileName := range inputFiles {
		contents, err := ioutil.ReadFile("testdata/" + fileName)
		assert.NoError(t, err)
		allPatches = append(allPatches, contents)
	}

	testCases := []struct {
		name         string
		withInverter bool
	}{
		{
			name: "normal",
		},
		{
			name:         "reverse",
			withInverter: true,
		},
	}

	unfilteredDiffs, err := NewUnidiff(NewBytesPatchReader(allPatches), zap.NewNop()).WithExpandedHunks().Decorate()
	assert.NoError(t, err)
	assert.Len(t, unfilteredDiffs, len(inputFiles))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, diff := range unfilteredDiffs {
				for _, hunk := range diff.Hunks {
					t.Run(diff.PreferredName+"@"+hunk.ID, func(t *testing.T) {
						// Parse again, and filter for this hunk
						d := NewUnidiff(NewBytesPatchReader(allPatches), zap.NewNop()).WithExpandedHunks().WithHunksFilter(hunk.ID)
						if tc.withInverter {
							d = d.WithInverter()
						}

						filteredDiff, err := d.Decorate()
						assert.NoError(t, err)
						if assert.Len(t, filteredDiff, 1) {
							assert.Len(t, filteredDiff[0].Hunks, 1)
						}
					})
				}
			}
		})
	}
}

func TestDecorateSeparateBinary(t *testing.T) {
	inputFiles := []string{
		"sample_changed.diff",
		"sample_binary_new_space_in_name.diff",
		"sample_file_extended_empty_new.diff",
		"500.diff",
	}

	var allPatches [][]byte
	for _, fileName := range inputFiles {
		contents, err := ioutil.ReadFile("testdata/" + fileName)
		assert.NoError(t, err)
		allPatches = append(allPatches, contents)
	}

	testCases := []struct {
		name         string
		withInverter bool
	}{
		{
			name: "normal",
		},
		{
			name:         "reverse",
			withInverter: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewUnidiff(NewBytesPatchReader(allPatches), zap.NewNop()).WithExpandedHunks()
			if tc.withInverter {
				d = d.WithInverter()
			}

			binaryDiffs, nonBinaryDiffs, err := d.DecorateSeparateBinary()
			assert.NoError(t, err)
			assert.Len(t, binaryDiffs, 2)
			assert.Len(t, nonBinaryDiffs, 2)
		})
	}
}

func TestFixLargeFilesDiffs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "new mode 100644",
			input: `diff --git a/aaa-100MB.dmg b/aaa-100MB.dmg
old mode 0
new mode 100644
Binary files /dev/null and /dev/null differ
`,
			expected: `diff --git a/aaa-100MB.dmg b/aaa-100MB.dmg
new file mode 100644
index 0000000..0000000
Binary files /dev/null and b/aaa-100MB.dmg differ
`,
		},
		{
			name: "space in name",
			input: `diff --git a/with a space.dmg b/with a space.dmg
old mode 0
new mode 100644
Binary files /dev/null and /dev/null differ
`,
			expected: `diff --git a/with a space.dmg b/with a space.dmg
new file mode 100644
index 0000000..0000000
Binary files /dev/null and b/with a space.dmg differ
`,
		},
		{
			name: "new mode 100777",
			input: `diff --git a/aaa-100MB.dmg b/aaa-100MB.dmg
old mode 0
new mode 100777
Binary files /dev/null and /dev/null differ
`,
			expected: `diff --git a/aaa-100MB.dmg b/aaa-100MB.dmg
new file mode 100777
index 0000000..0000000
Binary files /dev/null and b/aaa-100MB.dmg differ
`,
		},
		{
			name: "non changed file",
			input: `diff --git a/aaa-100MB.dmg b/aaa-100MB.dmg
new file mode 100777
index abc123..abc123
Binary files /dev/null and b/aaa-100MB.dmg differ
`,
			expected: `diff --git a/aaa-100MB.dmg b/aaa-100MB.dmg
new file mode 100777
index abc123..abc123
Binary files /dev/null and b/aaa-100MB.dmg differ
`,
		},
		{
			name: "normal diff",
			input: `diff --git a/500.txt b/500.txt
index 3f1dcfc..23c60d6 100644
--- a/500.txt
+++ b/500.txt
@@ -10,6 +10,13 @@
 10
 11
 12
+Binary files /dev/null and /dev/null differ
`,
			expected: `diff --git a/500.txt b/500.txt
index 3f1dcfc..23c60d6 100644
--- a/500.txt
+++ b/500.txt
@@ -10,6 +10,13 @@
 10
 11
 12
+Binary files /dev/null and /dev/null differ
`,
		},
		{
			name: "sample", // no changes
			input: `diff --git a/app/assets/bin/sturdy b/app/assets/bin/sturdy
old mode 100755
new mode 100644
index 16edd4f..9b8fb76
Binary files a/app/assets/bin/sturdy and b/app/assets/bin/sturdy differ
`,
			expected: `diff --git a/app/assets/bin/sturdy b/app/assets/bin/sturdy
old mode 100755
new mode 100644
index 16edd4f..9b8fb76
Binary files a/app/assets/bin/sturdy and b/app/assets/bin/sturdy differ
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := fixLargeFilesDiffs(tc.input)
			assert.Equal(t, tc.expected, out)
			assert.NoError(t, err)
		})
	}
}
