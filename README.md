# Simple Insecure File Transfer Protocol

A simple app to send and receive files over a TCP connection.

Each upload has a header which contains the protocol version, the file name, the file size, and the checksum of the file. After the header is sent, the file is sent over. The server receives the request and uses the header to verify the subsequently transferred file data.

## Running

```bash
# Start server
go run cmd/server/main.go
# Create a test file
echo "Hello world" > ./files/test.txt
# Run client
go run cmd/client/main.go ./files/test.txt
# Verify file was downloaded correctly
cat ./uploads/test.txt # Expecting 'Hello World'
```