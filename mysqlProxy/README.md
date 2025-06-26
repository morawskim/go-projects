This project is a lightweight MySQL proxy server implemented in Go.
The proxy listens for client connections and forwards queries to a backend MySQL database.
Additionally, it logs the original SQL queries being executed.
Designed as a starting point for more advanced database proxy functionalities.

## Usage

1. Start mysql: `docker compose up -d`
1. Start proxy: `go run ./main.go`
1. Connect to proxy: `mysql -h 127.0.0.1 -P3307 -udbname -puserpassword --ssl-mode=disabled dbname`
