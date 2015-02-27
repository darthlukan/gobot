package main

import (
	"fmt"
	"github.com/darthlukan/cakeday"
	"github.com/darthlukan/goconvtemps"
	"github.com/darthlukan/goduckgo/goduckgo"
	"strconv"
	"strings"
)

// GenericVerbCmd returns a message string based on the supplied cmd (a verb).
func GenericVerbCmd(cmd, extra string) string {
	randQuip := RandomQuip()
	return fmt.Sprintf("\x01"+"ACTION %v %v, %v\x01", cmd, extra, randQuip)
}

// CakeDayCmd returns a string containing the Reddit cakeday of a user
// upon success, or an error string on failure.
func CakeDayCmd(user string) string {
	var msg string

	responseString, err := cakeday.Get(user)
	if err != nil {
		msg = fmt.Sprintf("I caught an error: %v\n", err)
	} else {
		msg = fmt.Sprintf("%v\n", responseString)
	}
	return msg
}

// WebSearch takes a query string as an argument and returns
// a formatted string containing the results from DuckDuckGo.
func SearchCmd(query string) string {
	msg, err := goduckgo.Query(query)
	if err != nil {
		return fmt.Sprintf("DDG Error: %v\n", err)
	}

	switch {
	case len(msg.RelatedTopics) > 0:
		return fmt.Sprintf("First Topical Result: [ %s ]( %s )\n", msg.RelatedTopics[0].FirstURL, msg.RelatedTopics[0].Text)
	case len(msg.Results) > 0:
		return fmt.Sprintf("First External result: [ %s ]( %s )\n", msg.Results[0].FirstURL, msg.Results[0].Text)
	case len(msg.Redirect) > 0:
		return fmt.Sprintf("Redirect result: %s\n", UrlTitle(msg.Redirect))
	default:
		return fmt.Sprintf("Query: '%s' returned no results.\n", query)
	}
}

func ConvertTempCmd(query string) string {
	var unit string
	var converted string

	input := strings.ToLower(query)

	if strings.Index(input, "c") != -1 {
		unit = "c"
	} else if strings.Index(input, "f") != -1 {
		unit = "f"
	} else {
		return fmt.Sprintf("Invalid unit input, please use either 'F' or 'C'.\n")
	}

	temp, err := strconv.ParseFloat(fmt.Sprintf("%v", string(strings.Split(input, unit)[0])), 64)

	if err != nil {
		return fmt.Sprintf("Caught error '%v' trying to convert '%v'.\n", err, query)
	}

	converted = goconvtemps.ConvertTemps(temp, unit)

	return fmt.Sprintf("%v is %v.\n", strings.ToUpper(input), converted)
}

func WeatherCmd(query string) string {
	// Cheezy workaround until YQL lib is complete
	return fmt.Sprintf("Try '!search !google %v'\n", query)
}

func HelpCmd() string {
	return fmt.Sprintf("Available commands: !help, !ddg/search !weather, !convtemp, !cakeday, !VERB\n")
}

func WikiCmd(config *Config) string {
	return fmt.Sprintf("(Channel Wiki)[ %s ]\n", config.WikiLink)
}

func HomePageCmd(config *Config) string {
	return fmt.Sprintf("(Channel Homepage)[ %s ]\n", config.Homepage)
}

func ForumCmd(config *Config) string {
	return fmt.Sprintf("(Channel Forums)[ %s ]\n", config.Forums)
}
