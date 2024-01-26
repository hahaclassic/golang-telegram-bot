package telegram

import (
	"strconv"
	"strings"
)

func linkList(links []string, names []string) string {
	var enumerated strings.Builder

	for i, link := range links {
		if len(link) > maxCallbackMsgLen && len(names[i]) > maxCallbackMsgLen &&
			link[:maxCallbackMsgLen-5] == names[i][:maxCallbackMsgLen-5] || link == names[i] {
			enumerated.WriteString(strconv.Itoa(i+1) + ".\n" + link + "\n\n")
		} else {
			enumerated.WriteString(strconv.Itoa(i+1) + ". " + names[i] + "\n" + link + "\n\n")
		}
	}

	return enumerated.String()
}

func trimLink(url string) string {
	if len(url) > maxCallbackMsgLen {
		url = url[:maxCallbackMsgLen-5] + "..."
	}
	return url
}
