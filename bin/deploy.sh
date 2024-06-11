#!/usr/bin/env bash
die() { echo "${1:-argh}"; exit ${2:-1}; }

hash git  2>/dev/null || die "missing dep: git"
hash sam  2>/dev/null || die "missing dep: sam"
hash aws  2>/dev/null || die "missing dep: aws"

API_NAME="GameDay-api"
STACK_NAME="GameDay-api"
GIT_BRANCH=$(git branch --show-current) || die "git branch failed"

echo "~~~ :git: Get GIT info"
git_hash=$(git rev-parse HEAD)
git_repo=$(git config --get remote.origin.url)
tags="git:hash=${git_hash} git:branch=${GIT_BRANCH} ops:origin=${git_repo} ops:name=${API_NAME}"

echo "~~~ :aws: Build code using iterative golang compiling"
for path in ./cmd/*/; do
  dirname=$(basename "$path")

  echo "Building $dirname..."

  cd "./cmd/$dirname" || die "failed to cd into ./cmd/$dirname"
  GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o "../bin/$dirname/bootstrap" "$dirname.go"
  cd ../..
done

echo "~~~ :aws: Deploy SAM Stack for branch $GIT_BRANCH"
sam deploy \
  --tags "${tags}" \
  --no-fail-on-empty-changeset \
  --stack-name "${STACK_NAME}" \
  --s3-prefix "${STACK_NAME}" \
  --resolve-s3 \
  || die "sam deploy failed"