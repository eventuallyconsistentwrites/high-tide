# Load Testing `high-tide-server` with K6

This project uses **Grafana K6** running in Docker to generate load on the demo server.

## Architecture

This uses a **Hybrid Scaling** approach to maximize efficiency while maintaining unique IP addresses:

* **Scaling Containers:** We scale Docker containers to provide distinct source IP addresses (bypassing rate limits).
* **K6 Internals:** Each container runs multiple lightweight "Virtual Users" (VUs) internally.

**Example:**
Running **10 containers** with **100 VUs** each generates **1,000 Concurrent Users** originating from **10 distinct IPs**.

---

## Running the Test

### 1. Basic Run (Test Logic)
To verify the script works with a single container (1 IP):

```bash
docker-compose up k6-worker
```

### 2. Scaled Run (High Load, Multiple IPs)

To run a distributed test simulating traffic from multiple IP addresses, scale the `k6-worker` service:

```bash
# Starts 10 worker containers.
# If K6_VUS is set to 100, this generates 1,000 total concurrent users.
docker-compose up -d --scale k6-worker=10
```

> **Note:** You do not need the `--build` flag for the load tester anymore, as it pulls the official `grafana/k6` image and mounts your local script.

---

## Configuration

### Adjusting Load Intensity

You can control the number of users *per container* by setting the `K6_VUS` environment variable.

**Example: 5 containers x 50 users = 250 Total Users**

```bash
K6_VUS=50 docker-compose up --scale k6-worker=5
```

**Running for 5 minutes**

```bash
# Overriding the command to run for 5 minutes in detached mode
K6_VUS=50 K6_DURATION=5m docker-compose up -d --scale k6-worker=5
# Seeing logs
docker-compose logs -f k6-worker
```

### Modifying the Test Script

The test logic is defined in `loadtest.js`.

1. Edit `loadtest.js` locally (e.g., change the sleep duration or endpoint).
2. Restart the test. **No rebuild is required.**
---

## Monitoring

### View Logs

Watch the aggregate logs from all K6 workers to see the test progress:

```bash
docker-compose logs -f k6-worker
```

### Verify Source IPs

To confirm that your load is coming from different IP addresses, list the internal Docker IPs for all worker containers:

```bash
docker inspect --format '{{.Name}} - IP: {{.NetworkSettings.Networks.loadtester_default.IPAddress}}' $(docker-compose ps -q k6-worker)
```