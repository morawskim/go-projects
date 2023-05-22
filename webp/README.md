## Info

Simple CLI application to convert image to webp and create thumbnail also in webp.

## Usage

`go run ./main.go -image ./3000.jpg`

After run, you should see files `3000.webp` and `3000_thumbnail.webp`

Image `3000.jpg` downloaded from [picsum.photos](https://picsum.photos/)

## Image sizes

```
422K mar 13 16:25 3000.jpg
436 mar 13 18:11 3000_thumbnail.webp
200K mar 13 18:11 3000.webp
```

## Benchmark (quality of webp)

To see the current performance governor execute command `cat /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor`.
In my case all cores use schedutil policy.

To see available governors execute `cpupower frequency-info`.
In the output you should see a line similar to my:
> available cpufreq governors: ondemand performance schedutil

Next execute `sudo cpupower frequency-set -g performance`
You can again execute command `cat /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor` 
to see current governor policy. You should see a new policy.

Execute command to start benchmark: `go test -v -bench=. -benchtime=60s`

After benchmarking, you can restore your CPU governor policy.
