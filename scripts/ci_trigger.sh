#!/bin/sh
set -e

# run the test locally
sh scripts/ci_test.sh

git add cover.out coverage.out 
find . -type f -name 'checksum' -exec git add {} \;

# add all checksume
find . -type f -name 'checksum' -exec git add {} \;

NOW=$(date +%Y.%-m%d.%-H%M)
HAHSTAGS=${HAHSTAGS:-""}

git commit --allow-empty -m "ci($NOW): ✨🐛🚨 $HAHSTAGS"

TARGET=${1:-origin}
echo "---------------------------"
printf "Pushing... $NOW --> %s\n" "$TARGET"
echo "---------------------------"
git push "$TARGET"  