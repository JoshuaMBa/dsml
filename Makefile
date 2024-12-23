.PHONY: compile

compile:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative gpu_sim/proto/gpu_sim.proto

start-devices: kill-devices
	# sleep bc don't know when server thread starts
	mkdir -p logs
	> logs/device_pids
	go run gpu_sim/gpu_device/server/server.go --config=gpu_sim/run/device_configs/device1.json > logs/dev1.log 2>&1 &
	sleep 2 && pgrep -f "/exe/server --config=gpu_sim/run/device_configs/device1.json" | tail -n 1 >> logs/device_pids
	go run gpu_sim/gpu_device/server/server.go --config=gpu_sim/run/device_configs/device2.json > logs/dev2.log 2>&1 &
	sleep 2 && pgrep -f "/exe/server --config=gpu_sim/run/device_configs/device2.json" | tail -n 1 >> logs/device_pids
	go run gpu_sim/gpu_device/server/server.go --config=gpu_sim/run/device_configs/device3.json > logs/dev3.log 2>&1 &
	sleep 2 && pgrep -f "/exe/server --config=gpu_sim/run/device_configs/device3.json" | tail -n 1 >> logs/device_pids
	cat logs/device_pids

kill-devices:
	@if [ -f logs/device_pids ]; then \
		echo "Killing old device server PIDs:"; \
		cat logs/device_pids; \
		kill $$(cat logs/device_pids) || true; \
		rm logs/device_pids; \
	else \
		echo "No device servers running"; \
	fi

start-coordinator: kill-coordinator
	mkdir -p logs
	> logs/coordinator_pid
	go run gpu_sim/gpu_coordinator/server/server.go --device-list=gpu_sim/run/device_lists/three_devices.json > logs/coord1.log 2>&1 &
	sleep 2 && pgrep -f "/exe/server --device-list" | tail -n 1 >> logs/coordinator_pid
	cat logs/coordinator_pid

kill-coordinator:
	@if [ -f logs/coordinator_pid ]; then \
		echo "Killing old coordinator server PID:"; \
		cat logs/coordinator_pid; \
		kill $$(cat logs/coordinator_pid) || true; \
		rm logs/coordinator_pid; \
	else \
		echo "No coordinator servers running"; \
	fi

clean: kill-coordinator kill-devices
	rm -rf logs