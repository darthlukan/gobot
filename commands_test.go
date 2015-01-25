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
