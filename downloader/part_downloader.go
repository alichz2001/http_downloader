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

	readedBytes int64
	writedBytes int64

	request *http.HTTP

	bufferSize int64
	file       *os.File
}

func (pd *PartDownloader) getBytesCount() int64 {
	return pd.toByte - pd.fromByte + 1
}

func (pd *PartDownloader) Run() {

	file, err := os.OpenFile(pd.getFileName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	pd.file = file

	pd.request.Run()
	body := pd.request.GetResponseBodyWriter()

	p := make([]byte, pd.bufferSize)

	bytesCount := pd.toByte - pd.fromByte
	for i := int64(0); i < bytesCount-pd.bufferSize; i += pd.bufferSize {
		//TODO this approche make padding to end of file. append some empty bytes!
		//TODO handle error and maybe save state for failure
		n, _ := body.Read(p)
		pd.readedBytes += int64(n)

		//log.Print(n)
		m, _ := pd.file.Write(p)
		pd.writedBytes += int64(m)
		//log.Print(m)
	}
	remainBytes := pd.getBytesCount() - pd.readedBytes
	if remainBytes > 0 {
		c := make([]byte, remainBytes)
		n, _ := body.Read(c)
		pd.readedBytes += int64(n)
		m, _ := pd.file.Write(c)
		pd.writedBytes += int64(m)
	}
}

func (pd *PartDownloader) getFileName() string {
	return fmt.Sprintf("%s-%d.part", pd.name, pd.ID)
}
