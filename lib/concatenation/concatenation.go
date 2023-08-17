package conc

import (
	"strconv"
	"strings"
)

// Создание пронумерованного списка в строке из списка строк
func EnumeratedJoin(elements []string) string {

	var enumerated strings.Builder

	for i, elem := range elements {
		enumerated.WriteString(strconv.Itoa(i+1) + ". " + elem + "\n\n")
	}

	return enumerated.String()
}
