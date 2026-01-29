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
go build -o api ./cmd/api/main.go
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
./api
```

### Testing Count-min Sketch

```bash
go test -v ./countmin
```

---

## Directory Structure

### `cmd`
`api/main.go` contains application entry point

### `countmin`
Contains data structure logic of Count-Min Sketch

### `examples`
Contains demo go script to use `CountMinSketch` data structure

### `internal`

#### `post`
Interfaces for DB connection

#### `server`
Contains API logic

---