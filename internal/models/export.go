package models

type Export struct {
	Id      string `json:"Id"`
	Name    string `json:"Name"`
	PdfFile string `json:"PdfFile"`
	CsvFile string `json:"CsvFile"`
}
