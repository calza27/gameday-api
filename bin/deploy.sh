#!/usr/bin/env bash
die() { echo "${1:-argh}"; exit "${2:-1}"; }

hash sam  2>/dev/null || die "missing dep: sam"
hash aws  2>/dev/null || die "missing dep: aws"
hash ./bin/parse-yaml.sh || die "parse-yaml.sh not found."

profile=$1
[[ -z $profile ]] && die "Usage: $0 <profile>"

STACK_NAME="GameDay-api"

tags=$(./bin/parse-yaml.sh ./params/tags.yaml) || die "failed to parse tags"

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
  --s3-bucket "cfn-templates-3903---4367" \
  --s3-prefix "${STACK_NAME}" \
  --region "ap-southeast-2" \
  --capabilities "CAPABILITY_IAM" "CAPABILITY_AUTO_EXPAND" \
  --resolve-s3 \
  --profile "$profile" || die "sam deploy failed"