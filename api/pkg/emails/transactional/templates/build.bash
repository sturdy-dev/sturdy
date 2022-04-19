#!/usr/bin/env bash

set -euo pipefail

CWD="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

mkdir -p "${CWD}/output"
yarn install

yarn run mjml "${CWD}/welcome.template.mjml" -o "${CWD}/output/welcome.template.html"
yarn run mjml "${CWD}/notification/github_repository_imported.template.mjml" -o "${CWD}/output/notification/github_repository_imported.template.html"
yarn run mjml "${CWD}/notification/comment.template.mjml" -o "${CWD}/output/notification/comment.template.html"
yarn run mjml "${CWD}/notification/new_suggestion.template.mjml" -o "${CWD}/output/notification/new_suggestion.template.html"
yarn run mjml "${CWD}/notification/requested_review.template.mjml" -o "${CWD}/output/notification/requested_review.template.html"
yarn run mjml "${CWD}/notification/review.template.mjml" -o "${CWD}/output/notification/review.template.html"
yarn run mjml "${CWD}/verify_email.template.mjml" -o "${CWD}/output/verify_email.template.html"
yarn run mjml "${CWD}/magic_link.template.mjml" -o "${CWD}/output/magic_link.template.html"
yarn run mjml "${CWD}/invite_to_codebase.template.mjml" -o "${CWD}/output/invite_to_codebase.template.html"
yarn run mjml "${CWD}/invite_to_organization.template.mjml" -o "${CWD}/output/invite_to_organization.template.html"
