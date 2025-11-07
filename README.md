# Service - Go REST API

A lightweight REST API built with Go standard library, featuring CRUD operations, JSON file persistence, and full test coverage.

## Features

- RESTful API with full CRUD operations
- Thread-safe in-memory storage with JSON persistence
- Built with Go standard library (minimal dependencies)
- Comprehensive unit tests
- Docker containerization with multi-stage builds
- Automatic ID generation and timestamp management

## Project Structure

```
.
├── main.go                    # HTTP server and routing
├── handlers/
│   └── item_handler.go       # CRUD endpoint handlers
├── models/
│   └── item.go               # Data model
├── storage/
│   ├── storage.go            # Storage implementation
│   └── storage_test.go       # Unit tests
├── Dockerfile                # Container configuration
└── data.json                 # Persistent data file (auto-generated)
```

## Prerequisites

- Go 1.24.5 or higher
- Docker (optional, for containerized deployment)

## Building and Running

### Option 1: Run Locally

```bash
# Install dependencies
go mod download

# Run the service
go run main.go

# Or build and run the binary
go build -o service
./service
```

The server will start on `http://localhost:8080`

### Option 2: Run with Docker

```bash
# Build the Docker image
docker build -t service:latest .

# Run the container
docker run -d -p 8080:8080 --name service service:latest

# View logs
docker logs service

# Stop the container
docker stop service
docker rm service
```

### Option 3: Run with Docker (with persistent data)

To persist data outside the container, mount a volume:

```bash
docker run -d -p 8080:8080 -v $(pwd)/data:/root --name service service:latest
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/items` | Create a new item |
| GET | `/items` | Get all items |
| GET | `/items/{id}` | Get item by ID |
| PUT | `/items/{id}` | Update an item |
| DELETE | `/items/{id}` | Delete an item |

## Testing

```bash
go test ./...
```

## Data Persistence

Data is automatically persisted to `data.json` in the working directory. The file is created on first write and loaded on startup.

## Customizing the Data Model

To use your own data structure instead of the default `Item` model:

1. Update `models/item.go` with your struct definition
2. Ensure your struct has JSON tags for serialization
3. The CRUD operations will automatically work with your new model

## Configuration

- **Port:** Default is `8080` (can be modified in `main.go:21`)
- **Data file:** Default is `data.json` (can be modified in `storage/storage.go:27`)
