# ACME

Example project to obtain SSL certificate and private key via ACME protocol for example from Let's Encrypt.

## Usage

`docker compose up -d`

Check IP address of `docker0` interface - `ip a s docker0` 
Example output:

> 10: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default
link/ether 02:42:48:9f:ea:6c brd ff:ff:ff:ff:ff:ff
inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
valid_lft forever preferred_lft forever

From this output we see that the IP address of this interface is `172.17.0.1`.
In this example we use [nip.io](https://nip.io/) to use domain name which resolve to our docker0 IP address.

Run `go run ./main.go ./user.go -domain app.172.17.0.1.nip.io`
The files key.pem and cert.pem had been created.
