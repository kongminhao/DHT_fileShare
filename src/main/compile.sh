CGO_ENABLED=0 GOOS=linux GOARCH=386 go build daemon.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build client.go
scp daemon coffee:daemon
scp client coffee:client
