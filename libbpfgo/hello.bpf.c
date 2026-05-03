// +build ignore
#include "hello.bpf.h"

// Example: tracing a message on a kprobe
SEC("kprobe/sys_execve")
int hello(void *ctx)
{
    bpf_printk("I'm alive!");
    return 0;
}
