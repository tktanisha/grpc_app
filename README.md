# gRPC User CRUD example (Go)

This repository contains a simple Go gRPC server and client demonstrating Create/Get/Update/Delete/List for users.

Structure:
- `proto/` - the proto file (`user.proto`)
- `pb/` - minimal hand-written protobuf-like bindings for example (in real projects generate with protoc)
- `server/` - gRPC server
- `client/` - simple CLI client

Run server:

```powershell
cd c:\Users\khand\OneDrive\Desktop\grpc_app
go mod tidy
go run server\main.go
```

Run client examples (in another terminal):

```powershell
go run client\main.go -cmd=create -name="Alice" -email="alice@example.com"
go run client\main.go -cmd=list
```

Notes:
- The `pb` package here is hand-written for convenience. For production, generate code via `protoc --go_out` and `--go-grpc_out`.
- Uses an in-memory store; data is lost when the server restarts.
