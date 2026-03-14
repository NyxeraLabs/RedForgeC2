#!/bin/bash
# setup_tls.sh - Generate self-signed certificates for RedForgeC2

set -e

CERT_DIR="./certs"
mkdir -p "$CERT_DIR"

echo "[+] Generating self-signed certificates in $CERT_DIR..."

openssl req -x509 -newkey rsa:4096 -keyout "$CERT_DIR/server.key" -out "$CERT_DIR/server.crt" -days 365 -nodes -subj "/CN=localhost"

chmod 600 "$CERT_DIR/server.key"
chmod 644 "$CERT_DIR/server.crt"

echo "[+] Certificates generated successfully."
echo "[!] Ensure your docker-compose.yml mounts this directory and environment variables point to these files."