CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build netaddr.go
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build daemon.go
scp daemon coffee:daemon
scp netaddr groot108:netaddr