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

# nccl communicator / device relationship
- device is the lowest level, contains stuff like `{deviceId, memStart, memEnd}`
- communicator is the second lowest level. they wrap multiple raw devices to make it easier on end users
  - in nccl, a communicator is composed of many `ncclComm` structs, which represent individual devices
  - they know which `ncclComm` struct knows which communicator they belong to using commHash

  - we can use *groups* to group multiple gpu commands together

- my key takeaway from this is that: **gpu_device is a misleading name**; it doesn't represent the lowest level gpu device, but rather the communicator which wraps it
- this is how communicators can identify each other by rank alone, b/c they're not just naked devices

- reference, ncclComm struct: [https://github.com/NVIDIA/nccl/blob/2ea4ee94bfb04c886c79ccae60ac9961000fdee2/src/include/comm.h#L399]


# Implementation Detail
- we need to make some simplfying assumptions
  - based on how the lab is set up, spinning up a `gpu_device` process thread is equivalent to creating a communicator
  - basically, you don't start the server, then register it with a communicator. when you start the server, it is already assigned a communicator