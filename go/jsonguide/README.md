# jsonguide

## It does the trick

```sh
echo '{"K":"V","A":[1,2,{"e":true}]}' | jsonguide
```

```
.K = V
.A[0] = 1 (float64)
.A[1] = 2 (float64)
.A[2].e = true (bool)
```

## It is error-tolerant

```sh
echo '{"A":[1,{"q":[2' | jsonguide
```

```
.A[0] = 1 (float64)
.A[1].q[0] = 2 (float64)
.A[1].q[1]: [array] Unexpected EOF
```

In some cases it shows context of error:

```sh
echo '{"data":{"key-a":"a","key-b":***}}' | jsonguide
```

```
.data.key-a = a
.data.key-b: [value] Parse error: ("key-a":"a","key-b":***}}\n) invalid character '*' looking for beginning of value
```

## It supports multiple JSON objects

```sh
echo '{"X":1} {"X":2} {"X":3}' | jsonguide
```

```
.X = 1 (float64)
---
.X = 2 (float64)
---
.X = 3 (float64)
```

## It supports embedded JSONs

```sh
echo '{"A":false,"B":"{\"x\":[1,2],\"y\":true}","C":"just str"}' | jsonguide
```

```
.A = false (bool)
.B | .x[0] = 1 (float64)
.B | .x[1] = 2 (float64)
.B | .y = true (bool)
.C = just str
```

## It supports embedded base64 JSONs

```sh
echo '[1,2,3]' | openssl enc -base64 # WzEsMiwzXQo=
echo '{"A":"B","V":"WzEsMiwzXQo="}' | jsonguide
```

```
.A = B
.V # .[0] = 1 (float64)
.V # .[1] = 2 (float64)
.V # .[2] = 3 (float64)
```

## Supports timestamps in seconds and milliseconds

```sh
echo '[1777777777, 1777777777777]' | jsonguide
```

```
.[0] = 1777777777 (2026-05-03 03:09:37 UTC) (float64/timestamp)
.[1] = 1777777777777 (2026-05-03 03:09:37.777 UTC) (float64/timestamp)
```

## Supports UUIDv7

```sh
echo '{"id": "019ddddd-dddd-7ddd-0123-456789abcdef"}' | jsonguide
```

```
.id = 019ddddd-dddd-7ddd-0123-456789abcdef (2026-04-30 10:09:58.237 UTC) (string/UUIDv7)
```

## Supports non-unique keys

```sh
echo '{"A":"a","A":"b"}' | jsonguide
```

```
.A = a
.A = b
```

## Install it and enjoy

```sh
go install github.com/michurin/warehouse/go/jsonguide@latest
```
