package conc

import (
	"strconv"
	"strings"
)

// This function creates string like
// ---------------------------------------
// i. elements[i - 1] \n\n
// ---------------------------------------
// Where i - index of elements in the slice.
func EnumeratedJoin(elements []string) string {

	var enumerated strings.Builder

	for i, elem := range elements {
		enumerated.WriteString(strconv.Itoa(i+1) + ". " + elem + "\n\n")
	}

	return enumerated.String()
}

// This function creates string like
// ---------------------------------------
// i. tags[i - 1]
// data[i - 1] \n\n
// ---------------------------------------
// Where i - index of elements in the slice.
// len(data) must be equal to the len(tags), otherwise returns empty string
//
// If data[i] == tags[i], an empty string is written instead of tags[i]
func EnumeratedJoinWithTags(data []string, tags []string) string {
	if len(data) != len(tags) {
		return ""
	}
	var enumerated strings.Builder

	for i := 0; i < len(data); i++ {
		if data[i] == tags[i] {
			enumerated.WriteString(strconv.Itoa(i+1) + ".\n" + data[i] + "\n\n")
		} else {
			enumerated.WriteString(strconv.Itoa(i+1) + ". " + tags[i] + "\n" + data[i] + "\n\n")
		}
	}

	return enumerated.String()
}

// returns a cropped string. an 3 dots are added at the end
func trimData(data string, maxLen int) string {
	if len(data) > maxLen {
		data = data[:maxLen-3] + "..."
	}
	return data
}
