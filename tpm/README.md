# TPM 

## Features
The `getRandom.go` program connects to a TPM device (defaulting to `/dev/tpmrm0`) and retrieves 20 random bytes, which are then displayed in hexadecimal format.
In bash `tpm2_getrandom 20 --hex`

## Requirements
To run this program, you will need:
*   VirtualBox 7+ (supports TPM 2.0 emulation) or TPM 2.0

## Setup
The project includes a `Vagrantfile` that automatically prepares the testing environment:
1. Boots an Ubuntu Jammy 64-bit system.
1. Configures VirtualBox to emulate a TPM 2.0 device.
1. Installs necessary tools: `tpm2-tools` and the Go compiler.
1. Adds the user to the `tss` group to grant access to `/dev/tpmrm0`.

## How to Run getRandom.go

1. Start the virtual machine: `vagrant up`
1. Log in to the machine: `vagrant ssh`
1. Navigate to the project directory (mounted at `/vagrant`): `cd /vagrant`
1. Run the program: `go run getRandom.go`
