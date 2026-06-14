# CloudStack Services Documentation (Go Edition)

All services in CloudStack follow a consistent architectural pattern to ensure AWS SDK compatibility and persistence.

## Common Architecture

Each service is located in `internal/services/<service_name>` and contains:

1.  **Model (`/model/models.go`):** Structs representing AWS entities (e.g., `Bucket`, `User`, `Table`). These include JSON tags for persistence and API responses.
2.  **Service (`service.go`):** The business logic layer. It handles entity lifecycle, ID generation, and interacts with the storage backends.
3.  **Handler (`handler.go`):** The protocol layer. It implements either `HandleJSON` (for JSON 1.0/1.1), `HandleQuery` (for XML/Query), or `ServeHTTP` (for REST/S3/Lambda).

### Request Flow Diagram

```mermaid
graph TD
    SDK[AWS SDK / CLI] -->|HTTP Request| DISP[Core Dispatcher]
    DISP -->|Protocol Routing| HAND[Service Handler]
    HAND -->|Business Logic| SVC[Service Logic]
    SVC -->|Persistence| STORE[WalStorage / Memory]
    STORE -->|Write| DISK[(Data Directory)]
```

---

## Service Directory

| Service | Protocol | Storage | Status |
| :--- | :--- | :--- | :--- |
| **ACM** | JSON 1.1 | WAL | [x] Ported |
| **APIGATEWAY** | REST | WAL | [x] Ported |
| **APIGATEWAYV2** | JSON | WAL | [x] Ported |
| **APPCONFIG** | JSON 1.1 | WAL | [x] Ported |
| **APPRUNNER** | JSON 1.1 | WAL | [x] Ported |
| **APPSYNC** | JSON 1.1 | WAL | [x] Ported |
| **ATHENA** | JSON 1.1 | WAL | [x] Ported |
| **AUTOSCALING** | Query | WAL | [x] Ported |
| **BACKUP** | JSON 1.1 | WAL | [x] Ported |
| **BCMDATAEXPORTS** | JSON 1.1 | WAL | [x] Ported |
| **BEDROCKRUNTIME** | REST | N/A | [x] Ported |
| **CE** | JSON 1.1 | N/A | [x] Ported |
| **CLOUDFORMATION** | Query | WAL | [x] Ported |
| **CLOUDFRONT** | REST/XML | WAL | [x] Ported |
| **CLOUDSTACK** | JSON 1.1 | WAL | [x] Ported |
| **CLOUDTRAIL** | JSON 1.1 | WAL | [x] Ported |
| **CLOUDWATCH** | Query | WAL | [x] Ported |
| **CODEBUILD** | JSON 1.1 | WAL | [x] Ported |
| **CODECOMMIT** | JSON 1.1 | WAL | [x] Ported |
| **CODEDEPLOY** | JSON 1.1 | WAL | [x] Ported |
| **CODEPIPELINE** | JSON 1.1 | WAL | [x] Ported |
| **COGNITO** | JSON 1.1 | WAL | [x] Ported |
| **CONFIGSERVICE** | JSON 1.1 | WAL | [x] Ported |
| **CUR** | JSON 1.1 | WAL | [x] Ported |
| **DOCDB** | Query | WAL | [x] Ported |
| **DYNAMODB** | JSON 1.0 | WAL | [x] Ported |
| **EC2** | Query | WAL | [x] Ported |
| **ECR** | JSON 1.1 | WAL | [x] Ported |
| **ECS** | JSON 1.1 | WAL | [x] Ported |
| **EKS** | JSON 1.1 | WAL | [x] Ported |
| **ELASTICACHE** | Query | WAL | [x] Ported |
| **ELBV2** | Query | WAL | [x] Ported |
| **EMR** | JSON 1.1 | WAL | [x] Ported |
| **EVENTBRIDGE** | JSON 1.1 | WAL | [x] Ported |
| **FIREHOSE** | JSON 1.1 | WAL | [x] Ported |
| **GLUE** | JSON 1.1 | WAL | [x] Ported |
| **IAM** | Query | WAL | [x] Ported |
| **KINESIS** | JSON 1.1 | WAL | [x] Ported |
| **KMS** | JSON 1.1 | WAL | [x] Ported |
| **LAMBDA** | REST | WAL | [x] Ported |
| **MQ** | REST-JSON | WAL | [x] Ported |
| **MSK** | JSON 1.1 | WAL | [x] Ported |
| **NEPTUNE** | Query | WAL | [x] Ported |
| **OPENSEARCH** | REST/JSON | WAL | [x] Ported |
| **ORGANIZATIONS** | JSON 1.1 | WAL | [x] Ported |
| **PIPES** | JSON 1.1 | WAL | [x] Ported |
| **PRICING** | JSON 1.1 | N/A | [x] Ported |
| **RDS** | Query | WAL | [x] Ported |
| **REDSHIFT** | Query | WAL | [x] Ported |
| **RESOURCEGROUPSTAGGING** | JSON 1.1 | WAL | [x] Ported |
| **ROUTE53** | REST/XML | WAL | [x] Ported |
| **S3** | REST | File/WAL | [x] Ported |
| **SAGEMAKER** | JSON 1.1 | WAL | [x] Ported |
| **SCHEDULER** | JSON 1.1 | WAL | [x] Ported |
| **SECRETSMANAGER** | JSON 1.1 | WAL | [x] Ported |
| **SES** | Query | WAL | [x] Ported |
| **SNS** | Query | WAL | [x] Ported |
| **SQS** | Query | WAL | [x] Ported |
| **SSM** | JSON 1.1 | WAL | [x] Ported |
| **STEPFUNCTIONS** | JSON 1.0 | WAL | [x] Ported |
| **STS** | Query | WAL | [x] Ported |
| **SYNTHETICS** | JSON 1.1 | WAL | [x] Ported |
| **TEXTRACT** | JSON 1.1 | N/A | [x] Ported |
| **TRANSCRIBE** | JSON 1.1 | WAL | [x] Ported |
| **TRANSFER** | JSON 1.1 | WAL | [x] Ported |
| **WAF** | JSON 1.1 | WAL | [x] Ported |
| **XRAY** | JSON 1.1 | WAL | [x] Ported |

---

## GCP Services

| Service | Protocol | Storage | Status |
| :--- | :--- | :--- | :--- |
| **Artifact Registry** | REST-JSON | WAL | [x] Implemented |
| **App Engine** | REST-JSON | WAL | [x] Implemented |
| **BigQuery** | REST-JSON | WAL | [x] Implemented |
| **Bigtable** | REST-JSON | WAL | [x] Implemented |
| **Certificate Authority Service** | REST-JSON | WAL | [x] Implemented |
| **Cloud Armor** | REST-JSON | WAL | [x] Implemented |
| **Cloud Build** | REST-JSON | WAL | [x] Implemented |
| **Cloud DNS** | REST-JSON | WAL | [x] Implemented |
| **Cloud Functions** | REST-JSON | WAL | [x] Implemented |
| **Cloud Logging** | REST-JSON | WAL | [x] Implemented |
| **Cloud Load Balancing** | REST-JSON | WAL | [x] Implemented |
| **Cloud Monitoring** | REST-JSON | WAL | [x] Implemented |
| **Cloud Run** | REST-JSON | WAL | [x] Implemented |
| **Cloud Scheduler** | REST-JSON | WAL | [x] Implemented |
| **Cloud Spanner** | REST-JSON | WAL | [x] Implemented |
| **Cloud SQL** | REST-JSON | WAL | [x] Implemented |
| **Cloud Tasks** | REST-JSON | WAL | [x] Implemented |
| **Cloud Trace** | REST-JSON | WAL | [x] Implemented |
| **Cloud Workflows** | REST-JSON | WAL | [x] Implemented |
| **Compute Engine** | REST-JSON | WAL | [x] Implemented |
| **Datastore** | REST-JSON | WAL | [x] Implemented |
| **Firestore** | REST-JSON | WAL | [x] Implemented |
| **GCS** | REST-JSON | WAL | [x] Implemented |
| **GKE** | REST-JSON | WAL | [x] Implemented |
| **IAM** | REST-JSON | WAL | [x] Implemented |
| **Kafka** | REST-JSON | WAL | [x] Implemented |
| **Operations** | REST-JSON | WAL | [x] Implemented |
| **PubSub** | REST-JSON | WAL | [x] Implemented |
| **Secret Manager** | REST-JSON | WAL | [x] Implemented |

---

## Detailed Service Workings

### Compute Services (Lambda, EC2, ECS, EKS)
These services manage the lifecycle of virtual compute resources. Metadata is stored in WAL, while heavy artifacts (like Lambda ZIPs or ECR images) are stored in the local file system under `./data`.

### Storage Services (S3, DynamoDB)
*   **S3:** Direct mapping of bucket/key to directory/file.
*   **DynamoDB:** In-memory indexing of items with background WAL synchronization for persistence.

### Messaging Services (SQS, SNS, EventBridge, Pipes)
*   **SQS/SNS:** Uses memory-resident buffers for messages with asynchronous delivery triggers.
*   **EventBridge:** Implements pattern matching logic to route events to registered targets.

### Security Services (IAM, STS, Cognito, KMS, Secrets Manager)
*   **IAM/Cognito:** Manages multi-tenant identity partitions using the `AccountAwareBackend`.
*   **KMS:** Provides cryptographic stubs for encryption/decryption operations.
*   **STS:** Resolves temporary credentials and caller identity from Authorization headers.
