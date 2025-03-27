#!/bin/bash

# Create configs directory if it doesn't exist
mkdir -p configs

# Generate private key
openssl genrsa -out configs/key.pem 2048

# Generate self-signed certificate
openssl req -new -x509 -sha256 -key configs/key.pem -out configs/cert.pem -days 365 \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

echo "Certificate and key files have been generated in the configs directory."