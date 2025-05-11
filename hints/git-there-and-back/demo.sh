#!/bin/bash

set -f
set -e
set -xv

[[ -e .git ]] && { echo "Could you please delete .git?"; exit 1; }

export GIT_CONFIG_NOSYSTEM=true

git init -q -b master

## FIRST COMMIT (prepare state)

echo 'line 1' >FILE

git status -s

git add FILE

git status -s

git commit -m 'Line 1'

git status -s

git log --oneline

## SECOND COMMIT (our workflow)

echo 'line 2' >>FILE # (1) editing

git status -s

git add FILE # (2) adding

git status -s

git commit -m 'Line 2' # (3) commit

git status -s

git log --oneline

second_head_hash=$(git rev-parse --short HEAD) # just remember for final repairing

## UNDO STEP BY STEP

git reset --soft HEAD^ # (3) commit, and moves HEAD!

git status -s

git diff # empty
cat FILE # both Line 1 and Line 2
git log --oneline --reflog # previous HEAD is not evaporated

git reset # (2) undo add

git diff # diff is here
cat FILE # both Line 1 and Line 2

git reset --hard # (1) undo changes (DANGEROUS but mostly still not fatal)

git diff # nothing
cat FILE # only Line 1
git log --oneline --reflog # everything is repairable

## REPAIRING ALL UNDOS

git reset --hard $second_head_hash # repairing everything committed in second commit

git log --oneline
git status -s
git diff
cat FILE
