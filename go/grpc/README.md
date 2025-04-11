```
protoc --go_out=kit --go-grpc_out=kit api/api.proto
```

```
grpcurl -plaintext localhost:9898 list
grpcurl -plaintext -d '{"x":111}' localhost:9898 api.v1.CalsService/Square
```

```
grpcurl -plaintext localhost:9898 list api.v1.CalsService
grpcurl -plaintext localhost:9898 list grpc.reflection.v1.ServerReflection
grpcurl -plaintext localhost:9898 describe api.v1.CalsService
grpcurl -plaintext localhost:9898 describe api.v1.Number
```
