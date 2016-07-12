package main

import (
	random "crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	pingme "github.com/apourchet/pingme/lib"
	"io/ioutil"
	"os"
	"strings"
)

const (
	RAND_ID_SIZE = 16
)

var (
	aliasFlag = flag.Bool("alias", false,
		"pingme -alias <alias> <id>\n\tCreates an alias for this channel id.")
	rand_flag   = flag.Bool("r", false, "Creates a random channel and waits for it to be pinged.")
	listen_flag = flag.Bool("l", false,
		"pingme -l <id/alias>\n\tListens to a given channel.")
	numberFlag = flag.Int("n", -1, "pingme -l -n=<num> <id/alias>\n\tNumber of messages to wait for.")

	config pingme.ClientConfig
)

func listen(id string) {
	c := &pingme.Client{config.Host, config.Port}

	out := func(msg string) bool {
		fmt.Println(msg)
		*numberFlag -= 1
		if *numberFlag == 0 {
			return false
		}
		return true
	}

	err := c.Listen(id, out)
	exitOnError(err)
}

func ping(id string, msg string) {
	c := &pingme.Client{config.Host, config.Port}
	err := c.Ping(id, msg)
	exitOnError(err)
}

func resolveAlias(alias string) string {
	if id, ok := config.Aliases[alias]; ok {
		return id
	}
	return alias
}

func parseEvent(evt string) string {
	if !strings.Contains(evt, "data: ") {
		exitOnError(nil)
		return ""
	}
	msg := strings.Split(evt, "data: ")[1]
	return strings.Trim(msg, "\n \n")
}

func getMsgFromStdin() string {
	msgBytes, err := ioutil.ReadAll(os.Stdin)
	msg := strings.TrimSpace(string(msgBytes))
	exitOnError(err)
	return msg
}

func checkArgLength(l int) {
	if len(flag.Args()) < l {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("An error has occured")
		fmt.Println("\t", err)
		os.Exit(1)
	}
}

func aliasAction() {
	checkArgLength(2)
	id := flag.Args()[0]
	alias := flag.Args()[1]
	config.Aliases[alias] = id
	config.WriteOut()
}

func listenAction() {
	checkArgLength(1)
	id := resolveAlias(flag.Args()[0])
	listen(id)
}

func randAction() {
	b := make([]byte, RAND_ID_SIZE)
	_, err := random.Read(b)
	exitOnError(err)

	id := hex.EncodeToString(b)
	fmt.Println("Channel id: ", id)
	listen(id)
}

func pingAction() {
	checkArgLength(1)
	id := resolveAlias(flag.Args()[0])
	if len(flag.Args()) > 1 {
		ping(id, flag.Args()[1])
	} else {
		msg := getMsgFromStdin()
		ping(id, msg)
	}
}

func init() {
	flag.Parse()
}

func main() {
	config = pingme.GetClientConfig()

	if *aliasFlag {
		aliasAction()
	} else if *rand_flag {
		randAction()
	} else if *listen_flag {
		listenAction()
	} else {
		pingAction()
	}
}
