$env:GOOS="linux"
$env:GOARCH="amd64"

go build -o UbuntuSharedFolder/pglogalyze ./cmd/pglogalyze