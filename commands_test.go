package main

import (
	"strings"
	"testing"
)

func TestGenericVerbCmd(t *testing.T) {
	cmdResponse := GenericVerbCmd("slaps", "USER")
	if !strings.Contains(cmdResponse, "slaps") || !strings.Contains(cmdResponse, "USER") {
		t.Errorf("GenericVerbCmd(%v, %v), got %v, want 'slaps' and 'USER'", "slaps", "USER\n", cmdResponse)
	}
}

func TestCakeDayCmd(t *testing.T) {
	cakedayMsg := CakeDayCmd("darthlukan")
	if !strings.Contains(cakedayMsg, "darthlukan") || !strings.Contains(cakedayMsg, "9 June 2010") {
		t.Errorf("CakeDayCmd(%v), got %v, want 'Reddit cakeday for darthlukan is: 9 June 2010\n", "darthlukan", cakedayMsg)
	}
}

func TestHelpCmd(t *testing.T) {
	helpMsg := HelpCmd()
	if !strings.Contains(helpMsg, "Available commands: !help, !ddg/search") {
		t.Errorf("HelpCmd(), got %v, want 'Available commands: !help, !ddg/search'\n", helpMsg)
	}
}

func TestWikiCmd(t *testing.T) {
	config := &Config{}
	config.WikiLink = "WIKIURL"
	wikiUrl := WikiCmd(config)
	if !strings.Contains(wikiUrl, "WIKIURL") {
		t.Errorf("WikiCmd(%v), got %v, want 'WikiUrl'\n", config, wikiUrl)
	}
}

func TestHomePageCmd(t *testing.T) {
	config := &Config{}
	config.Homepage = "HOMEURL"
	homeUrl := HomePageCmd(config)
	if !strings.Contains(homeUrl, "HOMEURL") {
		t.Errorf("HomePageCmd(%v), got %v, want 'HOMEURL'\n", config, homeUrl)
	}
}

func TestForumCmd(t *testing.T) {
	config := &Config{}
	config.Forums = "FORUMURL"
	forumUrl := ForumCmd(config)
	if !strings.Contains(forumUrl, "FORUMURL") {
		t.Errorf("ForumCmd(%v), got %v, want 'FORUMURL'\n", config, forumUrl)
	}
}

func TestSearchCmd(t *testing.T) {
	topicalResult := SearchCmd("New York")
	if !strings.Contains(topicalResult, "First Topical Result:") {
		t.Errorf("SearchCmd(%v), got %v, want 'First Topical Result:'\n", "New York", topicalResult)
	}
	redirectResult := SearchCmd("!archwiki i3")
	if !strings.Contains(redirectResult, "Redirect result:") {
		t.Errorf("SearchCmd(%v), got %v, want 'Redirect result:'\n", "!archwiki i3", redirectResult)
	}
	noResult := SearchCmd("my face")
	if !strings.Contains(noResult, "returned no results") {
		t.Errorf("SearchCmd(%v), got %v, want 'Query: my face returned no results.'\n", "your face", noResult)
	}
}
