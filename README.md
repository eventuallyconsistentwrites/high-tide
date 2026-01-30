# 🌊 high-tide
Stateless SYN Flood detection using Count-Min Sketch. Tracks net connection imbalance with constant memory efficiency

Table of Contents
---
- [Building and Running](#building-and-running)
- [Directory Structure](#directory-structure)

---
## Building and Running

### Building the Project

```bash
go mod init github.com/eventuallyconsistentwrites/high-tide-server
go mod tidy
```

```bash
go build -o hts-main ./cmd/api/main.go
```

### Install Dependencies

```bash
go get github.com/mattn/go-sqlite3
go mod tidy
```

### Running the Project

```bash
go run ./cmd/api/main.go

# Or in case of pre-built binary
./hts-main
```

### Testing Count-min Sketch

```bash
go test -v ./countmin
```

---

## Directory Structure

### `api`
Contains server contract information

### `cmd`
`api/main.go` contains application entry point

### `countmin`
Contains data structure logic of Count-Min Sketch

### `examples`
Contains demo go script to use `CountMinSketch` data structure

### `internal`

#### `domain`
Contains definitions for data structures that will be used to structure the incoming request bodies as well as define sqlite tables.

#### `repository`
Related to ORM system that is used to make the sqlite queries.

#### `routes`
Defines APIs

#### `server`
Contains business logic that makes use of repositories. The results are returned to the APIs that queried the server.

---