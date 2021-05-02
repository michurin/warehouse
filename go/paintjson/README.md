# Simple JSON painter

## Install

    go get github.com/michurin/warehouse/go/paintjson

## Examples

    fmt.Println(paintjson.PJ(`{"x":12}`))

## Description

In fact, it doesn't perform full JSON parsing. It consider
spaces, quoted strings, brackets, colons, commas... In addition,
it emphasizes quoted strings right before colons and mark them
as keys.

Thanks to this, it can treat semi-JSON strings like this:

    Body: {"ok": true}

## Todo

- Refactoring: separate FSM
- Streaming: obtain `io.Reader`
- Make colors tunable
- CLI tool
