package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
)

func main() {
	var tpmname = flag.String("tpm", "/dev/tpmrm0", "The path to the TPM device to use")
	flag.Parse()

	thetpm, err := linuxtpm.Open(*tpmname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmname, err)
		return
	}
	defer thetpm.Close()

	response, err := tpm2.GetRandom{BytesRequested: 20}.Execute(thetpm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get randomBytes: %s\n", err)
		return
	}

	fmt.Println(hex.EncodeToString(response.RandomBytes.Buffer))
}
