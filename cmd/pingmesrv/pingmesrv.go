package main

import (
	"flag"
	"fmt"
	"github.com/apourchet/pingme/lib"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

const (
	unitFilename = "pingmesrv.service"
	systemdFile  = `[Unit]
Description=Pingme Server

[Service]
User=antoine
TimeoutStartSec=0
ExecStart=%s/bin/pingmesrv -p %s

[Install]
WantedBy=multi-user.target
`
)

var (
	systemd  *bool
	port     *string
	unitPath string // "/home/antoine/.config/systemd/user/pingme.service"
)

func init() {
	homeDir, err := getHomeDir()
	if err != nil {
		log.Println("Error: Could not locate home directory.")
		os.Exit(1)
	}
	unitPath = path.Join(homeDir, ".config/systemd/user/", unitFilename)

	systemd = flag.Bool("systemd", false, "pingmesrv -systemd\nInstalls systemd unit file.")
	port = flag.String("p", ping.DEFAULT_PORT, "Port to listen on")
	flag.Parse()
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func main() {
	if *systemd {
		file := fmt.Sprintf(systemdFile, os.Getenv("GOPATH"), *port)
		err := ioutil.WriteFile(unitPath, []byte(file), 0644)
		if err != nil {
			log.Println("Error: Could not write to", unitPath)
			log.Println(err)
			os.Exit(1)
		}
		log.Println("Enable this service:\n",
			"systemctl enable "+unitPath+"\n",
			"systemctl daemon-reload\n",
			"systemctl start "+unitFilename)
	} else {
		s := ping.NewServer(*port)
		log.Println("Listening on port " + *port)
		s.Serve()
	}
}
