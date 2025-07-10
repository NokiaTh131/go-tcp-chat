# Go Chat Server

A simple TCP-based chat server and client implementation in Go.

## Features

- Multi-client support
- Real-time messaging
- User join/leave notifications
- Built-in commands (`/help`, `/users`, `/quit`)
- Anonymous user support

## Building

```bash
# Build server
go build -o server server/server.go

# Build client
go build -o client client/client.go
```

## Usage

1. Start the server:
```bash
./server
```

2. Connect clients (in separate terminals):
```bash
./client
```

3. Enter your name when prompted and start chatting!

## Commands

- `/help` - Show available commands
- `/users` - List online users
- `/quit` - Leave the chat

## Alternative Connection

You can also connect using telnet:
```bash
telnet localhost 8080
```

i also use nc command to connect to the server:

```bash
nc localhost 8080
```

## Architecture

- **Server**: Handles multiple concurrent connections using goroutines
- **Client**: Simple terminal-based interface with real-time message display
- **Protocol**: Plain text over TCP on port 8080
