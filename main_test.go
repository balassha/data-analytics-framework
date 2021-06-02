package main

import (
	"hellofresh/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostcodeInput(t *testing.T) {
	cases := []struct {
		postCode string
		expected bool
	}{
		{
			postCode: "",
			expected: false,
		},
		{
			postCode: "1",
			expected: true,
		},
		{
			postCode: "a",
			expected: false,
		},
		{
			postCode: "aaaa",
			expected: false,
		},
		{
			postCode: "#$",
			expected: false,
		},
		{
			postCode: "#$$",
			expected: false,
		},
		{
			postCode: "10120",
			expected: true,
		},
	}
	for _, c := range cases {
		have := utils.IsNumeric(c.postCode)
		assert.Equal(t, c.expected, have, c.postCode)
	}
}

func TestStartingTimeInput(t *testing.T) {
	cases := []struct {
		startingTime string
		expected     bool
	}{
		{
			startingTime: "",
			expected:     false,
		},
		{
			startingTime: "1",
			expected:     false,
		},
		{
			startingTime: "a",
			expected:     false,
		},
		{
			startingTime: "11",
			expected:     false,
		},
		{
			startingTime: "aa",
			expected:     false,
		},
		{
			startingTime: "7AM",
			expected:     true,
		},
		{
			startingTime: "7PM",
			expected:     false,
		},
		{
			startingTime: "0AM",
			expected:     false,
		},
		{
			startingTime: "13AM",
			expected:     false,
		},
		{
			startingTime: "",
			expected:     false,
		},
	}
	for _, c := range cases {
		have := utils.IsTimeValid(c.startingTime, "AM")
		assert.Equal(t, c.expected, have, c.startingTime)
	}
}

func TestEndingTimeInput(t *testing.T) {
	cases := []struct {
		endingTime string
		expected   bool
	}{
		{
			endingTime: "",
			expected:   false,
		},
		{
			endingTime: "1",
			expected:   false,
		},
		{
			endingTime: "a",
			expected:   false,
		},
		{
			endingTime: "11",
			expected:   false,
		},
		{
			endingTime: "aa",
			expected:   false,
		},
		{
			endingTime: "7AM",
			expected:   false,
		},
		{
			endingTime: "7PM",
			expected:   true,
		},
		{
			endingTime: "0AM",
			expected:   false,
		},
		{
			endingTime: "13PM",
			expected:   false,
		},
		{
			endingTime: "",
			expected:   false,
		},
	}
	for _, c := range cases {
		have := utils.IsTimeValid(c.endingTime, "PM")
		assert.Equal(t, c.expected, have, c.endingTime)
	}
}
