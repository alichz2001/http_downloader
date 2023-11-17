package downloader

import "strings"

func getFileNameFromURI(uri string) string {
	splittedURI := strings.Split(uri, "/")
	return splittedURI[len(splittedURI)-1]
}
