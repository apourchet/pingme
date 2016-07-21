# PingMe

Finally a notification service with a simple API, accessible from command line.

### Installation
    go install github.com/apourchet/pingme/cmd/pingmesrv
    go install github.com/apourchet/pingme/cmd/pingme

### Server Setup
    cd $GOPATH/github.com/apourchet/pingme
    make create && make start

### Quickstart
If you want to try it out without setting up the server yourself:

    pingme -h antoinepourchet.com:1025

In one terminal window, you listen for a notification:

    pingme -l firstchannel -n 1

In another, you can ping the listening process:

    pingme -p firstchannel "This is a ping!"
    

