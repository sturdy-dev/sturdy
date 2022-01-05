#!/usr/local/bin/bash

mkdir templates/output
yarn install

yarn run mjml ./templates/welcome.template.mjml -o templates/output/welcome.template.html
yarn run mjml ./templates/github_repository_imported.template.mjml -o ./templates/output/github_repository_imported.template.html
yarn run mjml ./templates/notification/comment.template.mjml -o ./templates/output/notification/comment.template.html
yarn run mjml ./templates/notification/new_suggestion.template.mjml -o ./templates/output/notification/new_suggestion.template.html
yarn run mjml ./templates/notification/requested_review.template.mjml -o ./templates/output/notification/requested_review.template.html
yarn run mjml ./templates/notification/review.template.mjml -o ./templates/output/notification/review.template.html
yarn run mjml ./templates/verify_email.template.mjml -o ./templates/output/verify_email.template.html
yarn run mjml ./templates/magic_link.template.mjml -o ./templates/output/magic_link.template.html
