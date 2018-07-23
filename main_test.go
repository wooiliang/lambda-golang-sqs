package main

import "testing"

func TestGetJSON(t *testing.T) {
	jsonString := `
		{
			"foo": "bar"
		}
	`
	if _, err := getJSON(jsonString); err != nil {
		t.Errorf(`Expected JSON, got "%v"`, err)
	}
}
