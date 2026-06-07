# eslint tips short list

## install

```
npm i --prefix=$HOME/.local eslint
```

```
ln -s ../.local/node_modules/.bin/eslint $HOME/bin
```

## run on single file

```
eslint --no-config-lookup static/script.js
```

## useful comments

```
/*global fetch, EventSource, requestAnimationFrame, document:readable, setTimeout:readable, localStorage:readable*/
/*eslint indent: ["error", 2]*/
/*eslint eqeqeq: ["error", "always"]*/
/*eslint prefer-const: "error"*/
/*eslint no-var: "error"*/
/*eslint no-undef: "error"*/
/*eslint one-var: ["error", "never"]*/
/*eslint semi: ["error", "never"]*/
/*eslint quotes: ["error", "single"]*/
/*eslint prefer-arrow-callback: "error"*/
/*eslint arrow-body-style: "error"*/
```
