package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"syscall"

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
		token, err := TopicName2Token(editedTopic)
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
		token, err := TopicName2Token(editedTopic)
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
		token, err := TopicName2Token(editedTopic)
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

func TopicName2Token(topic string) (string, error) {
	tmp := strings.Replace(topic, "/", "", -1)
	if len(tmp) == 0 {
		return "", TopicNameError{fmt.Sprintf("Invalid topic name (inputed topic name: %v)", topic)}
	}
	var token uint64
	switch string(tmp[0]) {
	case "0":
		token = 0b0000000000000000000000000000000000000000000000000000000000000000
	case "1":
		token = 0b0010000000000000000000000000000000000000000000000000000000000000
	case "2":
		token = 0b0100000000000000000000000000000000000000000000000000000000000000
	case "3":
		token = 0b0110000000000000000000000000000000000000000000000000000000000000
	case "4":
		token = 0b1000000000000000000000000000000000000000000000000000000000000000
	case "5":
		token = 0b1010000000000000000000000000000000000000000000000000000000000000
	default:
		return "", TopicNameError{fmt.Sprintf("Invalid topic name (inputed topic name: %v)", topic)}
	}
	maskTail := uint64(0b0001000000000000000000000000000000000000000000000000000000000000)
	masks := [3]uint64{
		0b0000100000000000000000000000000000000000000000000000000000000000,
		0b0001000000000000000000000000000000000000000000000000000000000000,
		0b0001100000000000000000000000000000000000000000000000000000000000,
	}
	for _, v := range tmp[1:] {
		switch string(v) {
		case "0":
			// 何もしない
		case "1":
			token = token | masks[0]
		case "2":
			token = token | masks[1]
		case "3":
			token = token | masks[2]
		default:
			return "", TopicNameError{fmt.Sprintf("Invalid topic name (inputed topic name: %v)", topic)}
		}

		for j := 0; j < 3; j++ {
			masks[j] = masks[j] >> 2
		}
		maskTail = maskTail >> 2
	}
	tokenString := uint2Token(token | maskTail)
	tokenLen := 1
	tokenLen += int(len(tmp) / 2)
	return tokenString[:tokenLen], nil
}

func uint2Token(ui uint64) string {
	token := ""
	mask := uint64(0b1111000000000000000000000000000000000000000000000000000000000000)

	for i := 0; i < 16; i++ {
		tmp := (ui & mask)
		for j := i + 1; j < 16; j++ {
			tmp = tmp >> 4
		}
		token += fmt.Sprintf("%x", tmp)
		mask = mask >> 4
	}
	return token
}

//////////////           以下、エラー 関連                 //////////////
type TopicNameError struct {
	Msg string
}

func (e TopicNameError) Error() string {
	return fmt.Sprintf("Error: %v", e.Msg)
}

//////////////           以上、エラー 関連                 //////////////
