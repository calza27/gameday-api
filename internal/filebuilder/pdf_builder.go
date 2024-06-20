package filebuilder

import (
	"GameDay-API/internal/models"
	"bufio"
	"bytes"

	"github.com/jung-kurt/gofpdf"
)

type ByteWriteCloser struct {
	*bufio.Writer
}

func (bwc *ByteWriteCloser) Close() error {
	if err := bwc.Flush(); err != nil {
		return err
	}
	return nil
}

func BuildPdf(data models.GameData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, World!")

	var b bytes.Buffer
	byteWriter := bufio.NewWriter(&b)
	bwc := &ByteWriteCloser{byteWriter}
	err := pdf.OutputAndClose(bwc)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
