# libbpfgo eBPF Demo

This project is a small demo showing how to use [libbpfgo](https://github.com/aquasecurity/libbpfgo) from Go to load and attach an eBPF program.

The demo builds a simple eBPF program that attaches to the `execve` syscall using a kprobe.
Every time a new process is executed on the system, the eBPF program writes a short message to the kernel tracing pipe. 
The Go application loads the compiled eBPF object, attaches it to the kernel probe, and prints tracing output to the terminal.

## Using Vagrant

This project includes a `Vagrantfile` that creates a Rocky Linux VM with the required dependencies installed.

Start the VM:
`vagrant up`

Connect to virtual machine:
`vagrant ssh`

Go to the shared project directory in virtual machine:
`cd /vagrant`

Build the Go application and the eBPF object file:
`make all`

Run the demo with root privileges:
`sudo ./hello`

In another terminal, execute any command to trigger the `execve` syscall, for example: `ls`

You should see tracing output printed by the running demo, including the message from the eBPF program:

> run-parts-56134   [001] ....2.1  4300.731574: bpf_trace_printk: I'm alive!

Stop the demo with: `Ctrl+C`
