package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"tool/pkg/topicname2token"

	"github.com/golang/geo/s2"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	fName := flag.String("f", "", "File name")
	isHtml := flag.Bool("html", false, "Output by html format")
	level := flag.Int("l", -1, "Zoom level (0-30)")
	isCenterHold := flag.Bool("h", false, "Hold center position")
	flag.Parse()

	var filePointer *os.File = nil

	if *fName != "" {
		var err error
		filePointer, err = os.Open(*fName)
		if err != nil {
			panic(err)
		}
		defer filePointer.Close()

	} else if terminal.IsTerminal(syscall.Stdin) {
		var stdin string
		fmt.Print("Type filename then press the enter key: ")
		fmt.Scan(&stdin)

		var err error
		filePointer, err = os.Open(stdin)
		if err != nil {
			panic(err)
		}
		defer filePointer.Close()
	}

	isAdded := false
	subscribingTokens := []string{}
	counter := 1
	step := 1
	minLevel := 30
	centerLat, centerLng := 0., 0.
	isUpdatedCenter := false
	htmlString := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>`
	htmlLinkFormat := "<p><a href=\"https://s2.sidewalklabs.com/regioncoverer/?center=%f%%2C%f&zoom=%d&cells=%s\" target=\"_blank\" rel=\"noopener noreferrer\">%d</a></p>\n"
	if filePointer != nil {
		scanner := bufio.NewScanner(filePointer)
		for scanner.Scan() {
			input := regexp.MustCompile("\r\n|\n\r|\n|\r").ReplaceAllString(scanner.Text(), "")
			splitedInput := strings.Split(input, ",")
			if len(splitedInput) != 2 {
				panic("File format error")
			}
			cmd, topic := splitedInput[0], splitedInput[1]
			editedTopic := strings.Replace(topic, "/#", "", 1)
			token, err := topicname2token.TopicName2Token(editedTopic)
			if err != nil {
				panic(fmt.Sprintf("%v : %v, line %v", err, *fName, counter))
			}
			if cmd == "a" || cmd == "add" {
				subscribingTokens = append(subscribingTokens, token)
				isAdded = true
				id := s2.CellIDFromToken(token)
				if minLevel > id.Level() {
					minLevel = id.Level()
					if !*isCenterHold || !isUpdatedCenter {
						centerLat, centerLng = id.LatLng().Lat.Degrees(), id.LatLng().Lng.Degrees()
						isUpdatedCenter = true
					}
				}
			} else if cmd == "r" || cmd == "remove" {
				if isAdded {
					out := ""
					for _, v := range subscribingTokens {
						out += v + ","
					}
					if len(out) > 0 {
						if !*isHtml {
							fmt.Println(out)
						}
						tmpLevel := minLevel
						if *level > -1 {
							tmpLevel = *level
						}
						htmlString += fmt.Sprintf(htmlLinkFormat, centerLat, centerLng, tmpLevel, strings.Replace(out[:len(out)-2], ",", "%2C", -1), step)
						step++
						minLevel = 30
					}
				}
				l := len(subscribingTokens)
				for i := l - 1; i >= 0; i-- {
					if subscribingTokens[i] == token {
						subscribingTokens = remove(subscribingTokens, i)
					}
				}
				isAdded = false
			} else {
				panic(fmt.Sprintf("Unknown command in file: %v, line %v", *fName, counter))
			}
			counter++
		}
		out := ""
		for _, v := range subscribingTokens {
			id := s2.CellIDFromToken(v)
			if minLevel > id.Level() {
				minLevel = id.Level()
				centerLat, centerLng = id.LatLng().Lat.Degrees(), id.LatLng().Lng.Degrees()
			}
			out += v + ","
		}
		if len(out) > 0 {
			if !*isHtml {
				fmt.Println(out)
			}
			tmpLevel := minLevel
			if *level > -1 {
				tmpLevel = *level
			}
			htmlString += fmt.Sprintf(htmlLinkFormat, centerLat, centerLng, tmpLevel, strings.Replace(out[:len(out)-2], ",", "%2C", -1), step)
		}
		htmlString += `</body>
		</html>`
		if *isHtml {
			fmt.Println(htmlString)
		}
		return
	}
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
