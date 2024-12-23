.PHONY: run curl billing ngrok

run:
	go run main.go

curl:
	@curl 'localhost:8080/?url=https%3A%2F%2Fhttpbin.org%2Fdelay%2F8'

billing:
	@curl --no-buffer 'localhost:8080/?url=http://api.zeme.local:8096/data/v1/orders/supplier-summary?start_date=2024-10-01&end_date=2024-10-31' \
		-H 'Authorization: Bearer QeY<)(RHYZw9(]2+'

ngrok:
	 ngrok http 8080

# ngrok config add-authtoken aaa

# ngrok http --host-header=api.zeme.local 8096