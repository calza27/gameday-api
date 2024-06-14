# GameDay-API

API Gateway and Lambda functions to process requests for the GameDay API.

# Ideas

GET Endpoints
    List the games stored in the DDB
    Get a specific game by ID
    Get a URL to access a given file

PUT Endpoints
    Export game to CSV and PDF
    If this ends up being a long action to execute, put an SQS Queue with DLQ between the gateway and lambda 