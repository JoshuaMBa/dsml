
Run `make start-devices start-coordinator`. This creates a coordinator and three GPU's. The GPU's specs are found in gpu_sim/run/device_configs. Changing ports also needs to be modified in gpu_sim/run/device_lists.

Then, run `go run gpu_sim/run/simulation.go`, which tests various operations. Alternatively, `go test gpu_sim/tests/coordinator_test.go`.

*Group work*

Ayush worked primarily on the integration between GPU's and coordinators. This includes functions such as CommInit, MemCpy, and writing tests.

Michael worked primarily on the gpu_device side. This includes functions such as StreamSend, BeginSend, and BeginReceive.

Josh worked primarily on the gpu_coordinator side. This includes functions such as AllReduce and failure detection.
