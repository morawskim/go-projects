## Usage

1. Start mysql: `docker compose up -d`
1. Start proxy: `go run ./main.go`
1. Connect to proxy: `mysql -h 127.0.0.1 -P3307 -udbname -puserpassword --ssl-mode=disabled dbname`
