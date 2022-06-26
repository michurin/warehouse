# chat2

## About

(TODO)

## Start

```sh
go run ./examples/simple/...
```

You can specify binding address

```sh
go run ./examples/simple/... :9000
```

## Dev

### Cheap stress test

```javascript
t=$('#text'); setInterval(()=>{p = $.Event('keypress'); p.which=13; t.val(Math.random()); t.trigger(p)}, 100)
```
