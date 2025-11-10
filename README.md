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

### Option 4: Deploy to Google Cloud Run

For production deployment with automatic CI/CD:

1. Set up GitHub Actions secrets (see [DEPLOYMENT.md](./DEPLOYMENT.md))
2. Push to `main` branch or trigger workflow manually
3. Service deploys automatically to Cloud Run

See **[DEPLOYMENT.md](./DEPLOYMENT.md)** for complete setup instructions.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/items` | Create a new mushroom sighting |
| GET | `/items` | Get all sightings |
| GET | `/items/{id}` | Get sighting by ID |
| PUT | `/items/{id}` | Update a sighting |
| DELETE | `/items/{id}` | Delete a sighting |

## Data Model

The service uses a `MushroomSighting` model designed to match the UI's data structure:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "image": "data:image/jpeg;base64,...",
  "mushroomName": "Boletus edulis",
  "dateTime": "2025-11-09T19:24:00Z",
  "location": "Pacific Northwest forest",
  "count": 5,
  "created_at": "2025-11-09T19:24:10Z",
  "updated_at": "2025-11-09T19:24:10Z"
}
```

### Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Auto-generated | Unique identifier (UUID) |
| `image` | string | Optional | Base64 encoded image of the mushroom |
| `mushroomName` | string | Optional | User's identification of the mushroom species |
| `dateTime` | timestamp | **Required** | When the mushroom was found (ISO 8601) |
| `location` | string | **Required** | Where the mushroom was found |
| `count` | integer | **Required** | Number of mushrooms found (minimum 1) |
| `created_at` | timestamp | Auto-generated | When the record was created |
| `updated_at` | timestamp | Auto-generated | When the record was last updated |

## Example Requests

### Create a new sighting

```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
    "mushroomName": "Boletus edulis",
    "dateTime": "2025-11-09T14:30:00Z",
    "location": "Pacific Northwest forest",
    "count": 3
  }'
```

### Get all sightings

```bash
curl http://localhost:8080/items
```

### Get a specific sighting

```bash
curl http://localhost:8080/items/550e8400-e29b-41d4-a716-446655440000
```

### Update a sighting

```bash
curl -X PUT http://localhost:8080/items/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "mushroomName": "King Bolete",
    "dateTime": "2025-11-09T14:30:00Z",
    "location": "Pacific Northwest forest, near oak trees",
    "count": 5
  }'
```

### Delete a sighting

```bash
curl -X DELETE http://localhost:8080/items/550e8400-e29b-41d4-a716-446655440000
```

## Testing

```bash
go test ./...
```

## Data Persistence

Data is automatically persisted to `data.json` in the working directory. The file is created on first write and loaded on startup.

## Customizing the Data Model

The service currently uses a `MushroomSighting` model optimized for mushroom identification tracking. An `Item` type alias is maintained for backwards compatibility.

To further customize the data model:

1. Update `models/item.go` with your struct definition
2. Ensure your struct has JSON tags for serialization
3. Update validation in `handlers/item_handler.go` if needed
4. The CRUD operations will automatically work with your new model

## Configuration

- **Port:** Default is `8080` (configurable via `PORT` environment variable for Cloud Run)
- **Data file:** Default is `data.json` (can be modified in `storage/storage.go:27`)

## Production Deployment

For production deployment to Google Cloud Run with automated CI/CD via GitHub Actions, see **[DEPLOYMENT.md](./DEPLOYMENT.md)**.

Key features:
- Automatic deployment on push to `main`
- Container image storage in Google Artifact Registry
- Auto-scaling from 0 to 10 instances
- Pay-per-request pricing (~$0.50-$2/month for low traffic)

**Note:** Cloud Run is stateless. For production, consider migrating from `data.json` to:
- Cloud Storage (file-based)
- Cloud SQL (relational database)
- Firestore (NoSQL database)

See DEPLOYMENT.md for data persistence options.
