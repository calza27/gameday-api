# GameDay-API

API Gateway and Lambda functions to process requests for the GameDay API.

# notes
    VPC Endpoint allowing the lambdas to hit DynamoDB and S3 (?) is required in the VPC
        See [networks repo](https://github.com/calza27/network)

# Ideas

GET Endpoints
    List the games stored in the DDB
    Get a specific game by ID
    Get a presigned URL for access to a specific file

PUT Endpoints
    Export game to CSV and PDF
        If this ends up being a long action to execute, put an SQS Queue with DLQ between the gateway and lambda 