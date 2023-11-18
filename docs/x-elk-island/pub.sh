#/bin/sh

set -x
set -e

PUB=ei__

git diff --exit-code . 2>&1 >/dev/null || { echo "ERROR: code is changed: see git diff"; exit 1; }
git diff --cached --exit-code . 2>&1 >/dev/null || { echo "ERROR: code ain't committed: see git diff --cached"; exit 1; }
test -d "$PUB" && { echo "ERROR: $PUB exists!"; exit 1; }

mkdir -p "$PUB"
git rev-parse --short HEAD >"$PUB/version"
sed -e '/<head>/ r pub_ga_tag' -e '/script/ s/main\.js/main.min.js/' index.html >"$PUB/index.html"
uglifyjs -c -m toplevel --mangle <main.js >"$PUB/main.min.js" # DO NOT MANGLE PROPS
cp favicon128.png favicon.ico style.css "$PUB"
pushd "$PUB"
find . -type f -print0 | grep -zv 'version.md5' | sort -z | xargs -0 sha256sum >'version.md5'
popd
