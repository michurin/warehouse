#!/bin/sh

GITHUB_USER=michurin

curl -qs 'https://api.github.com/users/'"$GITHUB_USER"'/repos?per_page=100' |
    tee generated_data_raw.json |
    jq '{repos: [.[] | {url: .clone_url, fork: .fork, name: .name, pushed_at: .pushed_at}]}' |
    tee generated_data.json

j2 script.j2 generated_data.json |
    tee generated_backup_script.sh
