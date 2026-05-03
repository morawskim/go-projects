package main

import (
	"C"
	"bufio"
	"fmt"
	"os"
	"os/signal"

	bpf "github.com/aquasecurity/libbpfgo"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	bpfModule, err := bpf.NewModuleFromFile("hello.bpf.o")
	must(err)
	defer bpfModule.Close()

	err = bpfModule.BPFLoadObject()
	must(err)

	prog, err := bpfModule.GetProgram("hello")
	must(err)

	_, err = prog.AttachKprobe("__x64_sys_execve")
	must(err)

	go tracePrint()
	<-sig
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// TracePrint reads data from the trace pipe that bpf_trace_printk() writes to,
// and writes it to stdout. The pipe is global, so this function is not
// associated with any BPF program. It is recommended to use bpf_trace_printk()
// and this function for debug purposes only.
// This is a blocking function intended to be called from a goroutine, for example:
//
//	go tracePrint()
func tracePrint() {
	f, err := os.Open("/sys/kernel/debug/tracing/trace_pipe")
	if err != nil {
		fmt.Println("TracePrint failed to open trace pipe: %v", err)
		return
	}

	r := bufio.NewReader(f)
	b := make([]byte, 1000)
	for {
		len, err := r.Read(b)
		if err != nil {
			fmt.Println("TracePrint failed to read from trace pipe: %v", err)
			return
		}

		s := string(b[:len])
		fmt.Println(s)
	}
}
