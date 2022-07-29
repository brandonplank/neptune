package main

import (
	"archive/zip"
	"bytes"
	"github.com/ledongthuc/pdf"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"regexp"
	"strconv"
	"syscall/js"
	"time"
)

func DoesSliceHaveName(name string, slice []string) bool {
	for _, obj := range slice {
		if obj == name {
			return true
		}
	}
	return false
}

func ProcessFile(this js.Value, args []js.Value) interface{} {
	l := args[1].Int()
	pdfBytes := make([]byte, l)
	b := js.CopyBytesToGo(pdfBytes, args[0])
	log.Printf("recived %d bytes", b)
	pdfReader := bytes.NewReader(pdfBytes)
	r, err := pdf.NewReaderEncrypted(pdfReader, pdfReader.Size(), func() string { return "" })
	if err != nil {
		log.Println(err.Error())
		return js.Null()
	}
	totalPage := r.NumPage()
	log.Println("Page numbers " + strconv.Itoa(totalPage))
	var list string
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				list += word.S + "\n"
			}
		}
	}
	var re = regexp.MustCompile(`(?m)^([a-zA-Z\-]+)\s*,\s*([a-zA-Z]+)(\s+([a-zA-Z]+))?$`)
	var stored []string

	var buf bytes.Buffer
	writer := io.MultiWriter(&buf)
	zipWriter := zip.NewWriter(writer)
	for _, match := range re.FindAllString(list, -1) {
		if !DoesSliceHaveName(match, stored) {
			log.Println(match)
			createWriter, err := zipWriter.Create(match + ".png")
			if err != nil {
				log.Println(err.Error())
			}
			qr, err := qrcode.Encode(match, qrcode.Medium, 256)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			_, err = createWriter.Write(qr)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	defer zipWriter.Close()
	js.CopyBytesToJS(js.Global().Get("ReturnedBytes"), buf.Bytes())
	log.Println(buf.Len())
	return "success"
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	log.SetPrefix("[wasm] ")
	log.Println("Module by Brandon Plank, Copyright " + time.Now().Format("2006"))
	log.Println("Contact me at brplank@moreheadstate.edu or brandon@brandonplank.org")
	log.Println("Setting up PDF processor")
	js.Global().Set("ProcessPDF", js.FuncOf(ProcessFile))
	<-make(chan bool)
}
