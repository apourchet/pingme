# PingMe

Finally a notification service with a simple API, accessible from command line.

### Installation
    go install github.com/apourchet/pingme/cmd/pingmesrv
    go install github.com/apourchet/pingme/cmd/pingme

### Systemd Setup
    cd $GOPATH/github.com/apourchet/pingme
    make create && make start

### Quickstart
In one terminal window, you listen for a notification:

    pingme -l firstchannel -n 1

In another, you can ping the listening process:

    pingme -p firstchannel "This is a ping!"
    

