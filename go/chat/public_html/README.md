# howto uglify

```
npm install uglify-js
npm install eslint
npx uglifyjs --compress --mangle -- public_html/core.js
```

```
npm install terser
npx terser --compress --mangle -- public_html/core.js
```

```
git rev-parse --short=12 HEAD
```
