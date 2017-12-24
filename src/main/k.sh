CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build daemon.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build client.go
scp daemon groot108:daemon
scp client groot108:client
