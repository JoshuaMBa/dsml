# Compile protobuf:
- Check the makefile, but `make compile`

# Running the different servers (we'll move this to docker + kubernetes)
```
go run gpu_device/server/server.go
go run gpu_coordinator/server/server.go
```

# Infra details:
- logging: copied directly from lab1
- server: copied directly from lab1