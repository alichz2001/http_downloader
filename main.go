package main

import (
	"github.com/alichz2001/http_downloader/downloader"
	"log"
)

func main() {

	d := downloader.NewDownloader("localhost.com/Fight.Club.1999.Bluray.720p.Farsi.Dubbed.mkv", 2)

	err := d.Start()

	log.Printf("%#v", err)

}
