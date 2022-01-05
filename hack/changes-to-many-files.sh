#!/bin/bash

set -euxo pipefail

mkdir -p nested/{this,there}

date > nested/this/that.txt
date > nested/this/this-one.txt
date > nested/this/the-other-one.txt
date > nested/there/that.txt
date > nested/there/this-one.txt
date > nested/there/the-other-one.txt
date > nested/here.txt
date > nested/there.txt
date > root.txt

