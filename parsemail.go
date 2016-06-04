package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var (
	reClient *regexp.Regexp
	reOrigTo *regexp.Regexp
	reSender *regexp.Regexp
	reTo     *regexp.Regexp

	flagUniq    bool
	flagClient  bool
	flagOrigTo  bool
	flagSender  bool
	flagTo      bool
	flagLogfile string
)

func init() {
	// Declare and compile Regular Expressions
	var err error
	reTxtClient := "client=([^, \n]+)"
	reTxtOrigTo := "orig_to=<([^>]+)"
	reTxtSender := "from=<([^>]+)"
	reTxtTo := " to=<([^>]+)"
	reClient, err = regexp.Compile(reTxtClient)
	errTest(err)
	reOrigTo, err = regexp.Compile(reTxtOrigTo)
	errTest(err)
	reSender, err = regexp.Compile(reTxtSender)
	errTest(err)
	reTo, err = regexp.Compile(reTxtTo)
	errTest(err)

	// Specify command line flags
	flag.BoolVar(&flagUniq, "uniq", false, "Dedupe results")
	flag.BoolVar(
		&flagClient,
		"client",
		false,
		"Match Client MTA's",
	)
	flag.BoolVar(
		&flagOrigTo,
		"origto",
		false,
		"Match Original Recipient (orig_to)",
	)
	flag.BoolVar(
		&flagSender,
		"sender",
		false,
		"Match Envelope senders",
	)
	flag.StringVar(
		&flagLogfile,
		"log",
		"/var/log/mail.info",
		"Mail logfile to parse",
	)
	flag.BoolVar(
		&flagTo,
		"to",
		false,
		"Match Recipient (to)",
	)
	flag.Parse()
}

// Kludgy panic on error function
func errTest(err error) {
	if err != nil {
		panic(err)
	}
}

// IsMemberStr tests for the membership of a string in a slice
func IsMemberStr(s string, slice []string) bool {
	for _, n := range slice {
		if n == s {
			return true
		}
	}
	return false
}

/*
	Client  - Domain name of the sending server
	Sender - Email address of the Envelope Sender
*/
func main() {
	// Open the logfile for reading
	f, err := os.Open(flagLogfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var uniqChk []string // Used when testing for unique results
	scanner := bufio.NewScanner(f)
	// Iterate file line-by-line
	for scanner.Scan() {
		line := scanner.Text()
		var res []string // List of sub-matches on a line
		if flagClient {
			res = reClient.FindStringSubmatch(line)
		}
		if flagOrigTo && len(res) == 0 {
			res = reOrigTo.FindStringSubmatch(line)
		}
		if flagSender && len(res) == 0 {
			res = reSender.FindStringSubmatch(line)
		}
		// to and orig_to always occur on the same line so this will
		// never be parsed if orig_to is tested.
		if flagTo && len(res) == 0 {
			res = reTo.FindStringSubmatch(line)
		}
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
