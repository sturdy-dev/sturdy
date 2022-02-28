#!/bin/bash

set -euo pipefail
set -x

# generate-repo.sh generates a small git repository with some minor history
# The output is deterministic, and all commit IDs are always the same

# *   commit 9ae1f77ec4f86b7d23b6fd65123542cf61be4407 (HEAD -> sturdytrunk)
# |\  Merge: 60fd0d0 41cca62
# | | Author: Sturdy Testdata <support@getsturdy.com>
# | | Date:   Thu Feb 17 13:29:40 2022 +0100
# | |
# | |     Merge ws2
# | |
# | * commit 41cca6270030201cff9f2f0505caa172adc1f54c (ws2)
# | | Author: Sturdy Testdata <support@getsturdy.com>
# | | Date:   Thu Feb 17 13:29:40 2022 +0100
# | |
# | |     5
# | |
# * |   commit 60fd0d0ef0ce162ad77d9943c007f694507fd762
# |\ \  Merge: 6bef3d5 88e25ff
# | | | Author: Sturdy Testdata <support@getsturdy.com>
# | | | Date:   Thu Feb 17 13:29:40 2022 +0100
# | | |
# | | |     Merge ws1
# | | |
# | * | commit 88e25fff7846b1a2337bea9764d458314bbbcf94 (ws1)
# | |/  Author: Sturdy Testdata <support@getsturdy.com>
# | |   Date:   Thu Feb 17 13:29:40 2022 +0100
# | |
# | |       4
# | |
# * | commit 6bef3d58c51b96b85466bef9402a39bb27018495
# |/  Author: Sturdy Testdata <support@getsturdy.com>
# |   Date:   Thu Feb 17 13:29:40 2022 +0100
# |
# |       6
# |
# * commit d814344b187b430bf5e9f808b64ac9b79de9b9fc
# | Author: Sturdy Testdata <support@getsturdy.com>
# | Date:   Thu Feb 17 13:29:40 2022 +0100
# |
# |     3
# |
# * commit 93385ffb9bfd51168be623c62f29c5c14fafd924
# | Author: Sturdy Testdata <support@getsturdy.com>
# | Date:   Thu Feb 17 13:29:40 2022 +0100
# |
# |     2
# |
# * commit 8b2ba42b367455a2b1f55a8862951a9dc3c53f1e
#   Author: Sturdy Testdata <support@getsturdy.com>
#   Date:   Thu Feb 17 13:29:40 2022 +0100
#
#       1

# DIR=$(mktemp -d)
# pushd "$DIR"
# git init .

pushd "$1"

git config user.email "support@getsturdy.com"
git config user.name "Sturdy Testdata"
git config commit.gpgsign false

# This script is expected to be executed in a repository initialized from CreateBareRepoWithRootCommit and CloneRepo
git checkout -b tmp
git branch -D sturdytrunk
git checkout --orphan sturdytrunk
git branch -D tmp

commit() {
	num=$1
	ts="Thu Feb 17 13:29:40 CET 2022"
	GIT_AUTHOR_DATE=$ts \
		GIT_COMMITTER_DATE=$ts \
		GIT_AUTHOR_NAME="Sturdy Testdata" \
		GIT_COMMITTER_NAME="Sturdy Testdata" \
		GIT_AUTHOR_EMAIL="support@getsturdy.com" \
		GIT_COMMITTER_EMAIL="support@getsturdy.com" \
		git commit -m "$num"
}

write_and_commit() {
	num=$1
	file=$2
	echo "$num" >$file
	git add $file
	commit "$num"
}

write_and_commit 1 file.txt
write_and_commit 2 file.txt
write_and_commit 3 file.txt

git checkout -b ws1
write_and_commit 4 ws1.txt

git checkout sturdytrunk
git checkout -b ws2
write_and_commit 5 ws2.txt

git checkout sturdytrunk
write_and_commit 6 file.txt
git merge --no-commit ws1 && commit "Merge ws1"
git merge --no-commit ws2 && commit "Merge ws2"

git show HEAD
