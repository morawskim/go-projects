# SMTP Server

A simple SMTP server written in Go that receives emails 
and forwards their content as notifications using the [shoutrrr](https://github.com/containrrr/shoutrrr) library.

## Build

To build the software, ensure you have Go installed on your system. 
Run the following command in the project root:

```bash
go build -o smtp-server server.go
```

## Run

After building, you can start the server by providing a notification URL.
For example, to send notifications to Telegram:

```bash
./smtp-server -notification-url "telegram://token@telegram?channels=channel-id"
```

### Configuration

The notification URL can be provided in several ways (checked in this order):

1. Command-line flag: `-notification-url`
2. Environment variable: `NOTIFICATION_URL`
3. File path via environment variable: `NOTIFICATION_URL_FILE` (reads the notification URL from the specified file)
4. Systemd credentials: looks for a credential named `notification_url` in the `CREDENTIALS_DIRECTORY` (if set).

### Options

- `-notification-url`: Notification URL.
- `-addr`: SMTP server listening address in `host:port` format (default: `127.0.0.1:25`).

The server listens on `127.0.0.1:25` by default.

## Test

You can test the SMTP server using the `swaks` tool. 
Use the following Docker command to send a test email:

```bash
docker run --network host --rm morawskim/swaks:latest --to user@example.com --server 127.0.0.1:25
```

If everything is configured correctly, the server will receive the email 
and attempt to send a notification to the specified service.
