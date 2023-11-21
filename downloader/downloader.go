package downloader

import (
	"fmt"
	"github.com/alichz2001/http_downloader/http"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Downloader struct {
	path          string
	server        string
	uri           string
	partsCount    int
	contentLength int
	fileName      string

	partsWG sync.WaitGroup
	parts   []*PartDownloader
}

func NewDownloader(path string, partsCount int) *Downloader {
	tmp := strings.SplitN(path, "/", 2)

	server := tmp[0] + ":80"
	uri := "/" + tmp[1]

	return &Downloader{
		path:       path,
		server:     server,
		uri:        uri,
		fileName:   getFileNameFromURI(uri),
		partsCount: partsCount,
		parts:      make([]*PartDownloader, partsCount),
	}
}

func (d *Downloader) getInfo() error {
	headReq := http.NewHTTP(d.server, d.uri, "HEAD")
	defer headReq.Close()

	headReq.Run()

	if headReq.GetResponse().GetStatusCode() != "200" {
		//TODO error handeling
		return fmt.Errorf("not found")
	}

	if acceptRange, err := headReq.GetResponse().GetHeader("accept-ranges"); err != nil || acceptRange != "bytes" {
		d.partsCount = 1
		log.Printf("server does not support multipart download!")
	}

	if contentLength, err := headReq.GetResponse().GetHeader("content-length"); err == nil {
		d.contentLength, _ = strconv.Atoi(contentLength)
	} else {
		log.Fatal("server did not get content-length")
	}

	return nil
}

func (d *Downloader) Start() error {
	//TODO error handeling
	err := d.getInfo()
	if err != nil {
		return err
	}

	err = d.initParts()
	if err != nil {
		return err
	}

	err = d.runAllParts()
	if err != nil {
		return err
	}

	err = d.mergeParts()
	if err != nil {
		return err
	}

	err = d.clean()
	if err != nil {
		return err
	}

	return nil
}

func (d *Downloader) runAllParts() error {
	d.partsWG.Add(d.partsCount)

	for _, part := range d.parts {
		part := part
		go func() {
			defer d.partsWG.Done()
			part.Run()
		}()
	}

	return nil
}

func (d *Downloader) initParts() error {
	//TODO error handeling
	for i := 0; i < d.partsCount; i++ {
		pd, err := d.createPart(i)
		if err != nil {
			return err
		}
		d.parts[i] = pd
	}
	return nil
}

func (d *Downloader) createPart(partID int) (*PartDownloader, error) {
	part := &PartDownloader{
		ID:         partID,
		bufferSize: 1024,
	}

	bufferLengthDivided := int64(d.contentLength) / part.bufferSize
	perPartBufferCount := bufferLengthDivided / int64(d.partsCount)

	part.fromByte = perPartBufferCount * part.bufferSize * int64(partID)

	part.toByte = perPartBufferCount * part.bufferSize * int64(partID+1)

	if part.ID == d.partsCount-1 {
		part.toByte += bufferLengthDivided % int64(d.partsCount)
		part.toByte += int64(d.contentLength) % part.bufferSize
	}
	part.currentByte = part.fromByte

	bytesStr := fmt.Sprintf("bytes=%d-%d", part.fromByte, part.toByte)
	part.request = http.NewHTTP(d.server, d.uri, "GET").SetRequestHeader("Range", bytesStr)
	part.name = getFileNameFromURI(d.uri)
	return part, nil
}

func (d *Downloader) mergeParts() error {
	d.partsWG.Wait()

	targetFile, err := os.Create(d.fileName)
	if err != nil {
		log.Fatal(err)
	}

	for _, part := range d.parts {
		part := part
		//TODO error handeling
		//tmp, _ := os.OpenFile(part.getFileName(), 0, 0666)
		part.file.Seek(0, io.SeekStart)
		_, err := io.Copy(targetFile, part.file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Downloader) clean() error {
	for _, part := range d.parts {
		part := part

		part.file.Close()
		os.Remove(part.getFileName())
	}
	return nil
}
