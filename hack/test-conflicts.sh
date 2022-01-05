#!/bin/bash

set -euo pipefail

git rebase --abort || true

root_commit=$(git rev-list --max-parents=0 HEAD)
workspace_name=$(git branch --show-current)

git checkout $root_commit

# Create common ancestor
git branch -D sturdytrunk && git checkout -b sturdytrunk
# echo "func foo(a, b int) int" > foo.go
cat <<EOT > foo.go
package math

func Add(a, b int) int {
  return a + b
}
EOT

echo "Hey In Ancestor (2)" > file2.txt
git add foo.go file2.txt && git commit -m "Added file (in ancestor)"
git push -u origin sturdytrunk --force

# Based on the new sturdytrunk, create a new workspace
git branch -D $workspace_name && git checkout -b $workspace_name
cat <<EOT > foo.go
package math

func Add(x, y int) int {
  return y + x
}
EOT
echo "Hey In Workspace (2)" > file2.txt
git add foo.go file2.txt && git commit -m "Updated file (in workspace)"

echo "Hey In Workspace AGAIN (1)" > foo.go
git add foo.go && git commit -m "Updated file 1 (in workspace, again)"

echo "1" > new_file.txt && git add new_file.txt && git commit -m "New File 1"
echo "2" > new_file.txt && git add new_file.txt && git commit -m "New File 2"
echo "3" > new_file.txt && git add new_file.txt && git commit -m "New File 3"
echo "4" > new_file.txt && git add new_file.txt && git commit -m "New File 4"

git push -u origin $workspace_name --force

# Go back to the trunk, and make another change
git checkout sturdytrunk
cat <<EOT > foo.go
package math

func Add(a, b float64) float64 {
  return a + b
}
EOT
echo "Hey In Trunk (2)" > file2.txt
git add foo.go file2.txt && git commit -m "Updated file (in trunk)"
git push -u origin sturdytrunk

git checkout $workspace_name

