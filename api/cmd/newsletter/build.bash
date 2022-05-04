#!/usr/local/bin/bash

mkdir output
yarn install
yarn run mjml letters/2021-09-09.mjml -o output/2021-09-09.html
yarn run mjml letters/2021-09-29.mjml -o output/2021-09-29.html
yarn run mjml letters/2021-11-22.mjml -o output/2021-11-22.html
yarn run mjml letters/2021-12-07.mjml -o output/2021-12-07.html
yarn run mjml letters/2022-02-21.mjml -o output/2022-02-21.html
yarn run mjml letters/2022-03-10.mjml -o output/2022-03-10.html
yarn run mjml letters/2022-04-13.mjml -o output/2022-04-13.html
yarn run mjml letters/2022-05-03.mjml -o output/2022-05-03.html
