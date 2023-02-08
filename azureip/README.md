# Usage

## To Run
./azureip_get-linuxamd64  | jq -r '.values[].properties.addressPrefixes[]' | sort -u | wc -l

## To Build
GOOS=linux GOARCH=amd64 go build main.go
GOOS=darwin GOARCH=arm64 go build main.go

