package downloader

import (
	"fmt"
	"github.com/alichz2001/http_downloader/http"
	"os"
)

type PartDownloader struct {
	ID       int
	name     string
	receiver chan int

	fromByte    int64
	toByte      int64
	currentByte int64

	request *http.HTTP

	bufferSize int
	file       *os.File
}

func (pd *PartDownloader) Run() {

	file, err := os.OpenFile(pd.getFileName(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	pd.file = file

	pd.request.Run()
	body := pd.request.GetResponseBodyWriter()

	p := make([]byte, pd.bufferSize)

	bytesCount := pd.toByte - pd.fromByte
	for i := int64(0); i < bytesCount; i += int64(pd.bufferSize) {
		//TODO this approche make padding to end of file. append some empty bytes!
		//TODO handle error and maybe save state for failure
		body.Read(p)
		//log.Print(n)
		pd.file.Write(p)
		//log.Print(m)
	}
}

func (pd *PartDownloader) getFileName() string {
	return fmt.Sprintf("%s-%d.part", pd.name, pd.ID)
}
