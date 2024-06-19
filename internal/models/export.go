package models

type Export struct {
	Id      string `json:"Id" dynamodbav:"id"`
	Name    string `json:"Name" dynamodbav:"name"`
	PdfFile string `json:"PdfFile" dynamodbav:"pdfFile"`
	CsvFile string `json:"CsvFile" dynamodbav:"csvFile"`
}
