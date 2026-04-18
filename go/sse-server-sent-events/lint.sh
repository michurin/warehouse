#!/bin/sh -xe

# *** LEGACY APPROACH ***

eslint \
    --no-config-lookup \
    --parser-options '{"ecmaVersion": 2018}' \
    \
    --global console \
    --global window \
    --global document \
    --global location \
    --global URLSearchParams \
    \
    --rule '{"arrow-body-style": "error"}' \
    --rule '{"eqeqeq": ["error", "always"]}' \
    --rule '{"indent": ["error", 2]}' \
    --rule '{"no-undef": "error"}' \
    --rule '{"no-unused-vars": "error"}' \
    --rule '{"no-var": "error"}' \
    --rule '{"one-var": ["error", "never"]}' \
    --rule '{"prefer-arrow-callback": "error"}' \
    --rule '{"prefer-const": "error"}' \
    --rule '{"quotes": ["error", "single"]}' \
    --rule '{"semi": ["error", "never"]}' \
    "$@" \
    $(find "$(dirname "$0")" -type f -name '*.js' | sort)
