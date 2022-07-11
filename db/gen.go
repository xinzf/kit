package db

import "strings"

var (
	GenTypesMapping = map[string]func(detailType string) (dataType string){
		"int": func(detailType string) (dataType string) { return "int" },
		"tinyint": func(detailType string) (dataType string) {
			if strings.HasPrefix(detailType, "tinyint(1)") {
				return "bool"
			}
			return "int8"
		},
	}
)
