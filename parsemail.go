package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

// IsMemberStr tests for the membership of a string in a slice
func IsMemberStr(s string, slice []string) bool {
	for _, n := range slice {
		if n == s {
			return true
		}
	}
	return false
}

// makeSenderRegexp compiles a Regex suitable for Sender matching
func makeSenderRegexp() (reSender *regexp.Regexp) {
	reSender, err := regexp.Compile("from=<([^>]+)")
	if err != nil {
		panic(err)
	}
	return
}

// makeClientRegexp compiles a Regex suitable for Client matching
func makeClientRegexp() (reClient *regexp.Regexp) {
	reClient, err := regexp.Compile("client=([^, \n]+)")
	if err != nil {
		panic(err)
	}
	return
}

/*
	Client  - Domain name of the sending server
	Sender - Email address of the Envelope Sender
*/
func main() {
	var flagUniq bool
	var flagClient bool
	var flagLogfile string
	flag.BoolVar(&flagUniq, "uniq", false, "Dedupe results")
	flag.BoolVar(
		&flagClient,
		"client",
		false,
		"Match Clients instead of Senders",
	)
	flag.StringVar(
		&flagLogfile,
		"log",
		"/var/log/mail.info",
		"Mail logfile to parse",
	)
	flag.Parse()

	// Define the regex to match log records against
	var hit *regexp.Regexp
	if flagClient {
		hit = makeClientRegexp()
	} else {
		hit = makeSenderRegexp()
	}

	// Open the logfile for reading
	f, err := os.Open(flagLogfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var uniqChk []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		res := hit.FindStringSubmatch(line)
		if len(res) < 1 {
			// No submatches found. Move on!
			continue
		}
		if flagUniq {
			if IsMemberStr(res[1], uniqChk) {
				// Duplicate string
				continue
			}
			uniqChk = append(uniqChk, res[1])
		}
		fmt.Println(res[1])
	}
}
