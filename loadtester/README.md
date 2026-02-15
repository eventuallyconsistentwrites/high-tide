# Load Testing `high-tide-server`

## Running the Test

The most common way to run the load test is to build the latest code and immediately scale the number of clients.

To build and start 10 "clients" at once, run:

```bash
docker-compose up -d --build --scale tester=10
```

> Why this works:
>
> - Isolation: Each tester instance is a separate container with its own virtual eth0 interface and unique IP.
> - Scalability: You can scale this to 50 or 100 instances (resource permitting) instantly.

Logs: You can watch the tester logs to see it making requests to the server:

```bash
docker-compose logs -f tester
```

List the docker container IPs

```bash
docker inspect --format '{{.Name}} - IP: {{.NetworkSettings.Networks.loadtester_default.IPAddress}}' $(docker-compose ps -q tester)
```