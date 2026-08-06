# GameDay-API
 
A serverless REST API, written in Go, for exporting, storing and retrieving GameDay match data. Match data submitted by a client is converted to a CSV export, stored in S3, and indexed in DynamoDB so it can be listed, retrieved and downloaded later via a pre-signed URL. The whole stack (API Gateway, Lambda functions, S3 bucket and DynamoDB table) is defined as AWS SAM / CloudFormation templates and deployed with the AWS SAM CLI.
 
## Repository structure
 
```
gameday-api/
├── bin/                     Shell scripts used to build and deploy the stack
│   ├── deploy.sh            Compiles each Lambda and deploys the infra, processor and gateway stacks via `sam deploy`
│   └── parse-yaml.sh        Helper that flattens a YAML file into `key=value` pairs (used to pass stack tags)
├── cf/                      AWS SAM / CloudFormation templates that define the deployed infrastructure
│   ├── infra.yaml           Core data infrastructure: S3 bucket, DynamoDB table, SSM parameters, IAM logging role
│   ├── processors.yaml      Lambda functions, their execution role/policies, log groups and SSM parameters
│   ├── gateway.yaml         The API Gateway (SAM `Serverless::Api`), access logging and Lambda invoke permissions
│   └── tags.yaml            Common resource tags applied to every deployed stack
├── cmd/                     Entry points for each Lambda function (one `main` package per function)
│   ├── export/              `export` Lambda — builds and stores a CSV export of submitted game data
│   ├── list/                `list` Lambda — lists previously exported games
│   ├── getGame/              `getGame` Lambda — fetches a single exported game record by ID
│   └── getFileUrl/           `getFileUrl` Lambda — generates a pre-signed S3 URL for a stored export file
├── docs/
│   └── openapi.yaml         OpenAPI 3.0 specification for the public API, also consumed by `gateway.yaml` to define API Gateway routes/integrations
├── internal/                Shared application code, not exposed outside the module
│   ├── aws/
│   │   ├── awsclient/       Constructors for AWS SDK v2 clients (S3, SSM, DynamoDB)
│   │   └── awsconfig/       Shared AWS SDK configuration loading
│   ├── filebuilder/         Builds export files from game data (`csv_builder.go` implemented, `pdf_builder.go` present but currently unused)
│   ├── models/              Domain types: `GameData`, `Player`, `ScoringEvent`, `QuarterTime`, `AppStorageEvent`, `Export`
│   ├── repositories/        DynamoDB access layer for reading/writing `Export` records
│   └── utils/               Shared helpers, including building API Gateway proxy responses
├── go.mod / go.sum          Go module definition and dependency lockfile
└── README.md
```
 
Each function under `cmd/` is built as an individual, statically linked Go binary (`GOOS=linux GOARCH=amd64`, custom runtime `provided.al2023`) and packaged as a separate Lambda function.
 
## Infrastructure
 
Infrastructure is defined as three AWS SAM templates in `cf/`, deployed independently and in order by `bin/deploy.sh` (`infra` → `processors` → `gateway`) using `sam deploy`, all in the `ap-southeast-2` region.
 
### `infra.yaml` — core data resources
- **S3 bucket** (`gameday-export-bucket`) — stores generated CSV/PDF export files.
- **DynamoDB table** (`GameDay_Exports`) — on-demand (`PAY_PER_REQUEST`) table keyed on `id`, storing metadata about each export (name, file references, generation date).
- **IAM role** shared by API Gateway/Lambda for writing CloudWatch logs, plus an `AWS::ApiGateway::Account` resource wiring that role in as the account-level CloudWatch role.
- **SSM parameters** publishing the bucket name/ARN, table name/ARN, and the configurable pre-signed file URL duration (`FileUrlDuration`, default `60s`) for consumption by the other stacks.
### `processors.yaml` — Lambda functions
- A shared **Lambda execution role** (`FunctionRole`) with permissions for CloudWatch Logs, VPC ENI management, S3 object read/write/delete/list on the export bucket, DynamoDB read/write/query/scan on the exports table, and SSM `GetParameter`.
- Four **Lambda functions**, all Go binaries running on the `provided.al2023` custom runtime with a 30s timeout, 128 MB memory, and an auto-published `live` alias:
  - `gameday-api-processor-export`
  - `gameday-api-processor-list`
  - `gameday-api-processor-getGame`
  - `gameday-api-processor-getFileUrl`
- A dedicated **CloudWatch log group** (365-day retention) per function.
- **SSM parameters** publishing each function's ARN so `gateway.yaml` can wire up API Gateway integrations without a hard dependency between stacks.
### `gateway.yaml` — public API
- An **`AWS::Serverless::Api`** resource (stage `api`) whose route/integration definitions are imported from `docs/openapi.yaml` via `Fn::Transform`/`AWS::Include`.
- **Access logging** to a dedicated CloudWatch log group, plus X-Ray tracing, INFO-level execution logging and full metrics/data tracing on all methods.
- **CORS** allowing `PUT`/`GET` from any origin with a `Content-Type` header.
- **Lambda invoke permissions** granting API Gateway permission to invoke the `live` alias of each of the four processor functions.
### `tags.yaml`
Common tags (`sec:datatype`, `ops:name`, `ops:origin`, `client`) applied to every deployed stack via `bin/parse-yaml.sh`.
 
## API
 
Base path (per `docs/openapi.yaml`): `/api/gameday`. All endpoints are implemented as `aws_proxy` Lambda integrations and validated against the OpenAPI schema (`x-amazon-apigateway-request-validator: basic`) before hitting the Lambda.
 
| Method | Path | Lambda | Description |
|---|---|---|---|
| `PUT` | `/gameday/export` | `export` | Accepts a `GameData` JSON payload, builds a CSV export of the match, stores it in S3, and writes an `Export` record (id, name, file reference, generation date) to DynamoDB. Returns `201` on success. PDF export is implemented in `internal/filebuilder/pdf_builder.go` but currently commented out in the handler. |
| `GET` | `/gameday/list` | `list` | Returns a JSON array of all previously exported games (`ExportedGame`: `id`, `name`, `pdfFile`, `csvFile`) from DynamoDB. |
| `GET` | `/gameday/get/{id}` | `getGame` | Returns a single exported game record by its DynamoDB `id`. |
| `GET` | `/gameday/getFileUrl/{file}` | `getFileUrl` | Checks the named file exists in S3, then returns a pre-signed download URL valid for the configured `FileUrlDuration`. Returns `404` if the file doesn't exist. |
 
Common response schema (`error`): `{ "status": <int>, "message": <string> }`, returned for `400`/`403`/`404`/`500` responses.
 
### Request/response shapes
 
- **`GameData`** (request body for `PUT /gameday/export`): match metadata (`Id`, `GameDate`, `Competition`, `TeamA`/`TeamB` names and abbreviations, `Venue`, `Level`, `Round`, summary fields) plus nested arrays for `TeamAPlayers` (`Player`: id, surname, given name, number, selected), `ScoringEvents` (per-quarter scoring detail), `QuarterTimes`, and `AppStorage` (arbitrary typed data blobs).
- **`ExportedGame`** (response body for `list`/`getGame`): `id`, `name`, `pdfFile`, `csvFile`.
## Build & deploy
 
`bin/deploy.sh <aws-profile>`:
1. Cross-compiles each function under `cmd/*` to `cmd/bin/<name>/bootstrap` for `linux/amd64`.
2. Runs `sam deploy` for `cf/infra.yaml`, then `cf/processors.yaml`, then `cf/gateway.yaml`, each as a separate CloudFormation stack (`GameDay-api-infra`, `GameDay-api-processors`, `GameDay-api-gateway`) tagged from `cf/tags.yaml`.
3. Cleans up the local `cmd/bin` build output.
Requires the `sam` and `aws` CLIs, an AWS CLI profile with access to deploy the stacks, and an existing SSM parameter `/s3/cfn-bucket/name` pointing at the S3 bucket used for SAM deployment artifacts.