package main

import (
	"archive/zip"
	"bytes"
	"github.com/jung-kurt/gofpdf"
	"github.com/ledongthuc/pdf"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"regexp"
	"strconv"
	"syscall/js"
	"time"
)

func removeDuplicateUser(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

	var userPdfBuffer bytes.Buffer
	userPdfWriter := io.MultiWriter(&userPdfBuffer)

	var buf bytes.Buffer
	writer := io.MultiWriter(&buf)
	zipWriter := zip.NewWriter(writer)
	var userList []string
	for _, match := range re.FindAllString(list, -1) {
		userList = append(userList, match)
		log.Println(match)
	}
	log.Println("Removing duplicates")
	userList = removeDuplicateUser(userList)
	log.Println("Creating QR codes and PDF")

	newPdf := gofpdf.New("P", "mm", "A4", "")
	for _, user := range userList {
		createWriter, err := zipWriter.Create("codes/" + user + ".png")
		if err != nil {
			log.Println(err.Error())
		}
		qr, err := qrcode.Encode(user, qrcode.Medium, 256)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		_, err = createWriter.Write(qr)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		newPdf.AddPage()
		newPdf.SetFont("Arial", "B", 16)

		newPdf.CellFormat(190, 7, user, "0", 0, "CM", false, 0, "")

		newPdf.RegisterImageReader(user, "png", bytes.NewReader(qr))
		newPdf.ImageOptions(
			user, 1, 1,
			0, 0,
			false,
			gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
			0,
			"",
		)
	}

	if newPdf.Error() != nil {
		log.Println(newPdf.Error().Error())
		return js.Null()
	}

	log.Printf("pages: %d", newPdf.PageCount())
	err = newPdf.Output(userPdfWriter)
	if err != nil {
		log.Println(err.Error())
		return js.Null()
	}

	userWriter, err := zipWriter.Create("users.pdf")
	if err != nil {
		log.Println(err.Error())
	}
	_, err = userWriter.Write(userPdfBuffer.Bytes())
	if err != nil {
		log.Println(err.Error())
		return js.Null()
	}

	err = zipWriter.Close()
	if err != nil {
		log.Println(err.Error())
		return js.Null()
	}
	ret := js.Global().Get("Uint8Array").New(buf.Len())
	js.CopyBytesToJS(ret, buf.Bytes())
	return ret
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
