# Benchmarking — Distributed Routing Engine

## Goal

Validate the routing engine under extreme load:
- **15,000 concurrent connections**
- **30 seconds sustained**
- **Measure latency distribution (p50/p90/p99) and RPS**

## Quick Start (WSL)

### Step 1: Open WSL Terminal
```bash
wsl
```

### Step 2: Navigate to Benchmark Folder
```bash
cd /mnt/c/Users/daksh/Projects/DistributedRoutingEngine/benchmark
```

### Step 3: Set File Descriptor Limit
```bash
ulimit -n 65535
```
Verify:
```bash
ulimit -n  # Should return 65535
```

### Step 4: Make Server Binary Executable
```bash
chmod +x test_server_linux
```

### Step 5: Start Server in Background
```bash
./test_server_linux &
```
(Press Enter if prompt doesn't return immediately)

### Step 6: Verify Server is Running
```bash
curl http://localhost:8080/health
```

### Step 7: Install `hey` Benchmark Tool (if not present)
```bash
go install github.com/rakyll/hey@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### Step 8: Run the Full 15k Concurrent Benchmark
```bash
hey -z 30s -c 15000 -m POST -D dummy_route_request.json http://localhost:8080/route
```

## Expected Output

The `hey` command will output:
- **Requests/sec**: Should be 78,000+ 
- **Latency distribution**: p50, p90, p99 percentiles
- **Success rate**: All 1,000,000 requests should be HTTP 200

## Result Interpretation

See [../README.md](../README.md#performance--benchmarks) for the system architecture analysis and why p99 latency may spike while p50 remains under 5ms.

## Troubleshooting

**Port 8080 already in use:**
```bash
killall -9 test_server_linux
sleep 1
./test_server_linux &
```

**Permission denied on test_server_linux:**
```bash
chmod +x test_server_linux
```

**Go not installed in WSL:**
```bash
sudo apt update && sudo apt install -y golang-go
```
