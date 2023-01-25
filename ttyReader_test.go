package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataParsing(t *testing.T) {
	f, err := os.Open("./resources/test-data")
	reader := bufio.NewReader(f)
	readData, err := reader.ReadString('\x21')
	if err != nil {
		panic(err)
	}

	matches := matchData(readData)

	assert.Len(t, matches, 48)

	kwhMap := map[string]string{
		"omis": "2.8.1",
		"data": "795.989",
		"unit": "kWh",
	}

	assert.Contains(t, matches, kwhMap)

	kwMap := map[string]string{
		"omis": "1.7.0",
		"data": "0.348",
		"unit": "kW",
	}

	assert.Contains(t, matches, kwMap)

	kvarMap := map[string]string{
		"omis": "4.7.0",
		"data": "0.296",
		"unit": "kvar",
	}

	assert.Contains(t, matches, kvarMap)

	vMap := map[string]string{
		"omis": "32.7",
		"data": "232.4",
		"unit": "V",
	}

	assert.Contains(t, matches, vMap)

	aMap := map[string]string{
		"omis": "31.7",
		"data": "0.58",
		"unit": "A",
	}

	assert.Contains(t, matches, aMap)
}
