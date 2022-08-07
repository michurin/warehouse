# chat2

## About

Tools, components, ideas, drafts for chatting web app.

## Content overview

### stream (go)

Storage for chat messages.

The core thing, that well done, tested, and even little bit documented.

Operations:

- `put` messages
- `wait` for new messages with `ctx`
- `get` new messages

Approaches:

- Bounds to detect new messages
- Ring buffer
- No private messages (if you need, just use separate storage)
- Message structure agnostic

### js/kit.js

Simple client-side chat adapter.

Oversimplification: designed for only one source of messages. One stream, one bound etc.

### text (go)

Validation. Drafts, ideas, examples.

### handler (go)

Just http handlers. Examples.

### htdocs

Common static.

### examples

#### examples/effortless + local htdocs

Very simple demo.

#### examples/simple + local htdocs

Simple demo.

#### examples/minesweeper + local htdocs (upcoming)

Example, how to use many streams in same web application.
And streams has individual `/pub` endpoints and *share* one `/sub` endpoint to safe connections.

## Hints

### Start

```sh
go run ./examples/simple/...
```

You can specify binding address

```sh
go run ./examples/simple/... :9000
```

### Dev

#### Cheap stress test

```javascript
t=$('#text'); setInterval(()=>{p = $.Event('keypress'); p.which=13; t.val(Math.random()); t.trigger(p)}, 100)
```

