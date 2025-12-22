package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-tpm/legacy/tpm2"
)

func main() {
	var tpmname = flag.String("tpm", "/dev/tpmrm0", "The path to the TPM device to use")
	flag.Parse()

	tpmRwc, err := tpm2.OpenTPM(*tpmname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmname, err)
		return
	}
	defer tpmRwc.Close()

	randomBytes, err := tpm2.GetRandom(tpmRwc, 20)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get randomBytes: %s\n", err)
		return
	}

	fmt.Println(hex.EncodeToString(randomBytes))
}
