$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o polydictor ./cmd/polydictor

scp .\polydictor chetyredva@YOUR_SERVER_IP:/home/chetyredva/Polydictor/