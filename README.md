# GameDay-API

API Gateway and Lambda functions to process requests for the GameDay API.

# Ideas

GET Endpoints
    List the games stored in the DDB
    Get a specific game by ID
    Get a presigned URL for access to a specific file

PUT Endpoints
    Export game to CSV and PDF
        https://unidoc.io/post/write-pdfs-in-golang-beginners-guide/
        https://pkg.go.dev/github.com/jung-kurt/gofpdf/v2#section-readme
        https://pkg.go.dev/github.com/jung-kurt/gofpdf/v2#Pdf