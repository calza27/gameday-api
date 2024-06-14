# GameDay-API

API Gateway and Lambda functions to process requests for the GameDay API.

# Ideas

GET endpoints to get the current status of aspects of the game - trigger a lambda to return data from DynamoDB 
    Current score
    Scoring events

PUT/POST endpoints to create events
    API Gateway -> SQS FIFO Queue (with DLQ) -> Lambda -> DynamoDB