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


# GPU Communication Semantics
- GPUs perform computation
- Coordinator initiates data movement
- Coordinator specifies which GPU is the sender and which GPU is the receiver
  - Coordinator first tells sender 
    - what data that the sender owns that it should send
    - which device (specified by rank) the sender should send its data to
  - Coordinator (after communicating with sender) then tells receiver 
    - where the receiver should store received data
    - who is sending the data (specified by rank)
    - which channel of communication the current receive operation belongs to (specified by a stream id)
- When coordinator communicates with sender, the sender
  - generates a new stream object identified by a stream id
  - kicks of sending data into the stream
  - returns stream id to coordinator
- Coordinator uses stream id received from sender to inform receiver of the channel of communication
- If the sender begins sending messages into a new stream before the receiver learns (from the coordinator) where to store the received data, the receiver hangs until the coordinator indicates where to store the received data.
