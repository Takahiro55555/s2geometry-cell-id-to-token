package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"syscall"
	"tool/pkg/topicname2token"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	topic := flag.String("t", "", "Topic name ex) /5/0/1/2/3/0/1/2/3/#")
	delimiter := flag.String("d", "\n", "Delimiter")
	flag.Parse()

	editedTopic := strings.Replace(*topic, "/#", "", 1)
	if *topic == "" && len(flag.Args()) == 1 {
		editedTopic = strings.Replace(flag.Args()[0], "/#", "", 1)
	}

	if editedTopic != "" {
		token, err := topicname2token.TopicName2Token(editedTopic)
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
		return
	}

	if terminal.IsTerminal(syscall.Stdin) {
		// Execute: go run main.go
		fmt.Print("Type topicname then press the enter key: ")
		var stdin string
		fmt.Scan(&stdin)
		editedTopic = strings.Replace(stdin, "/#", "", 1)
		token, err := topicname2token.TopicName2Token(editedTopic)
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
		return
	}

	// Execute: echo "foo" | go run main.go
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return
	}
	for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(body), -1) {
		if v == "" {
			continue
		}
		editedTopic = strings.Replace(v, "/#", "", 1)
		token, err := topicname2token.TopicName2Token(editedTopic)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		fmt.Printf("%s%s", token, *delimiter)
	}
	if *delimiter != "\n" {
		fmt.Println("")
	}
}
