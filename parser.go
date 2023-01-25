package main

import (
	"regexp"
)

var pex, _ = regexp.Compile(`(?P<omis>\d+\.\d+(\.\d)?)\((?P<data>\d+\.\d+)\*(?P<unit>\S+)\)`)

func matchData(str string) []map[string]string {

	match := pex.FindAllStringSubmatch(str, -1)
	var data []map[string]string
	for ind := range match {
		subMatchMap := make(map[string]string)
		for i, name := range pex.SubexpNames() {

			if name == "omis" || name == "data" || name == "unit" {
				subMatchMap[name] = match[ind][i]
			}

		}
		data = append(data, subMatchMap)
	}

	return data
}
