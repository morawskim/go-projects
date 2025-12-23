package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var JWT_HEADER = base64.URLEncoding.EncodeToString([]byte(`{"alg":"ES256","typ":"JWT"}`))
var JWT_PAYLOAD = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":"123456789","iss":"tpm","exp":%d}`, time.Now().Add(4*time.Hour).Unix())))

func main() {
	var tpmPath = flag.String("tpm", "/dev/tpmrm0", "The path to the TPM device to use")
	flag.Parse()

	tpmRwc, err := tpm2.OpenTPM(*tpmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmPath, err)
		return
	}
	defer tpmRwc.Close()

	digest := sha256.Sum256([]byte(JWT_HEADER + "." + JWT_PAYLOAD))

	signature, err := tpm2.Sign(
		tpmRwc,
		tpmutil.Handle(0x81010001),
		"",
		digest[:],
		nil,
		&tpm2.SigScheme{
			Alg:  tpm2.AlgECDSA,
			Hash: tpm2.AlgSHA256,
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't sign: %s\n", err)
		return
	}

	rBytes := signature.ECC.R.FillBytes(make([]byte, 32))
	sBytes := signature.ECC.S.FillBytes(make([]byte, 32))

	rawSig := append(rBytes, sBytes...)
	jwtSig := base64.RawURLEncoding.EncodeToString(rawSig)

	fmt.Println(JWT_HEADER + "." + JWT_PAYLOAD + "." + jwtSig)
}
