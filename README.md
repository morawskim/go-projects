# Go projects

## encrypt-pdf

HTTP server with one endpoint "/encrypt" which encrypt provided PDF file (encoded via base64).

* Integration with external Go package pdfcpu to encrypt PDF file
* Using gorilla/mux package to create REST POST endpoint
* Decode income JSON

Example request
```
{
  "user_pw": "password",
  "owner_pw": "owner_password",
  "file": "<PDF_FILE_ENCODED_IN_BASE64>"
}
```

Example response
```
{
    "file": "<ENCRYPTED_PDF_FILE_ENCODED_IN_BASE64>"
}
```

Curl (see also test.sh)

```
curl -XPOST \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -s \
  -d"$REQUEST_BODY" localhost:8080/encrypt | jq -r '.file' | base64 -d > api.pdf

```
