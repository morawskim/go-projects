package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
)

var JWT_HEADER = base64.URLEncoding.EncodeToString([]byte(`{"alg":"ES256","typ":"JWT"}`))
var JWT_PAYLOAD = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":"123456789","iss":"tpm","exp":%d}`, time.Now().Add(4*time.Hour).Unix())))

func main() {
	var tpmPath = flag.String("tpm", "/dev/tpmrm0", "The path to the TPM device to use")
	flag.Parse()

	thetpm, err := linuxtpm.Open(*tpmPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the TPM file %s: %s\n", *tpmPath, err)
		return
	}
	defer thetpm.Close()

	keyHandler := tpm2.TPMHandle(0x81010001)
	pubResp, err := tpm2.ReadPublic{
		ObjectHandle: keyHandler,
	}.Execute(thetpm)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read public key: %s\n", err)
		return
	}

	digest := sha256.Sum256([]byte(JWT_HEADER + "." + JWT_PAYLOAD))
	resp, err := tpm2.Sign{
		KeyHandle: tpm2.NamedHandle{
			Handle: keyHandler,
			Name:   pubResp.Name,
		},
		Digest: tpm2.TPM2BDigest{
			Buffer: digest[:],
		},
		InScheme: tpm2.TPMTSigScheme{
			Scheme: tpm2.TPMAlgECDSA,
			Details: tpm2.NewTPMUSigScheme(
				tpm2.TPMAlgECDSA,
				&tpm2.TPMSSchemeHash{
					HashAlg: tpm2.TPMAlgSHA256,
				},
			),
		},
		Validation: tpm2.TPMTTKHashCheck{
			Tag: tpm2.TPMSTHashCheck,
		},
	}.Execute(thetpm)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't sign: %s\n", err)
		return
	}

	signature, err := resp.Signature.Signature.ECDSA()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get ECDSA: %s\n", err)
		return
	}

	rBuffer := signature.SignatureR.Buffer
	sBuffer := signature.SignatureS.Buffer

	rawSig := append(leftPad(rBuffer, 32), leftPad(sBuffer, 32)...)
	jwtSig := base64.RawURLEncoding.EncodeToString(rawSig)

	fmt.Println(JWT_HEADER + "." + JWT_PAYLOAD + "." + jwtSig)
}

func leftPad(in []byte, size int) []byte {
	if len(in) >= size {
		return in
	}
	out := make([]byte, size)
	copy(out[size-len(in):], in)

	return out
}
