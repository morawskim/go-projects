
generateRequestBody()
{
  cat <<EOF
{
  "user_pw": "password",
  "owner_pw": "owner_password",
  "file": "$1"
}
EOF
}

FILE=$(base64 fixture/test.pdf | tr -d '\n' )
REQUEST_BODY=$(generateRequestBody "$FILE")

curl -XPOST \
  -H "Accept: application/json" \
  -H "Content-Type:application/json" \
  -s \
  -d"$REQUEST_BODY" localhost:8080/encrypt | jq -r '.file' | base64 -d > api.pdf
