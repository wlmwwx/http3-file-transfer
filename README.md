# HTTP3 File Upload/Download Application

This is a Go application that implements file upload and download functionality using HTTP3 protocol.

## Prerequisites

- Go 1.22 or later
- OpenSSL (for generating certificates)

## Setup

1. Generate certificates:
```bash
# Make the script executable
chmod +x scripts/generate_cert.sh
# Generate certificates
./scripts/generate_cert.sh
```

2. Create necessary directories:
```bash
mkdir -p uploads downloads
```

3. Install dependencies:
```bash
go mod tidy
```

## Running the Application

1. Start the server:
```bash
go run cmd/h3_server/main.go
```

2. Use the client to upload/download files:

Upload a file:
```bash
go run cmd/h3_client/main.go -upload path/to/file
```

Download a file:
```bash
go run cmd/h3_client/main.go -download filename
```

## Configuration

The application uses YAML configuration files located in the `configs` directory:

- `configs/server.yaml`: Server configuration
- `configs/client.yaml`: Client configuration

You can modify these files to change settings such as:
- Server host and port
- Maximum file size
- Allowed file types
- Upload/download directories
- TLS certificate locations

## Security Notes

- For production use, replace the self-signed certificates with proper TLS certificates
- Review and adjust the allowed file types in server configuration
- Set appropriate file size limits
- Implement proper authentication and authorization