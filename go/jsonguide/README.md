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

## Install it and enjoy

```sh
go install github.com/michurin/warehouse/go/jsonguide@latest
```
