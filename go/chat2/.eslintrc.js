module.exports = {
  env: {
    browser: true,
    es6: true,
  },
  extends: [
    'airbnb-base',
  ],
  globals: {
    '$': true
  },
  parserOptions: {
    ecmaVersion: 'latest',
  },
  rules: {
    'no-plusplus': 0
  },
};
