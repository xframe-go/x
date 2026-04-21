package requests

import "strings"

type FilterValue any

var (
	False = FilterValue(false)
	True  = FilterValue(true)
)

func convertValue(value string) FilterValue {
	if value == "true" {
		return True
	}

	if value == "false" {
		return False
	}

	if strings.Contains(value, ",") {
		return FilterValue(strings.Split(value, ","))
	}

	return FilterValue(value)
}
