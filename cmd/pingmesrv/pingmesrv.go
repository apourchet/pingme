package main

import (
	"flag"
	"fmt"
	"github.com/apourchet/pingme/lib"
	"io/ioutil"
	"log"
)

const (
	systemdFilename = "/etc/systemd/system/pingme.service"
	systemdFile     = `[Unit]
Description=Pingme

[Service]
TimeoutStartSec=0
ExecStart=pingmesrv -p %s

[Install]
WantedBy=multi-user.target
`
)

var (
	systemd *bool
	port    *string
)

func init() {
	systemd = flag.Bool("systemd", false, "pingmesrv -systemd\nInstalls systemd unit file")
	port = flag.String("p", ping.DEFAULT_PORT, "Port to listen on")
	flag.Parse()
}

func main() {
	if *systemd {
		file := fmt.Sprintf(systemdFile, *port)
		err := ioutil.WriteFile(systemdFilename, []byte(file), 0644)
		if err != nil {
			log.Println("Error: Could not write to", systemdFilename)
			log.Println(err)
		}
	} else {
		s := ping.NewServer(*port)
		log.Println("Listening on port " + *port)
		s.Serve()
	}
}
