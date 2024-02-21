package tools

import (
	"fmt"
	"regexp"
)

// data_source zstack image validate
func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}
