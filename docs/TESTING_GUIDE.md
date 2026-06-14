# CloudStack Integration Testing Guide

This guide provides instructions on how to verify the compatibility of each of the 65+ supported AWS services in CloudStack using the AWS CLI and the built-in BATS test suite.

## Global Setup

Before testing, ensure CloudStack is running and your environment is configured:

```bash
# Start CloudStack
go run ./cmd/cloudstack/main.go

# Configure AWS CLI environment
export AWS_ENDPOINT_URL=http://localhost:4566
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
```

## Running the Automated Test Suite

We use [BATS (Bash Automated Testing System)](https://github.com/bats-core/bats-core) for integration testing.

### Run All Tests
```bash
./compatibility-tests/lib/bats-core/bin/bats compatibility-tests/sdk-test-awscli/test/
```

### Run a Specific Service Test
```bash
./compatibility-tests/lib/bats-core/bin/bats compatibility-tests/sdk-test-awscli/test/<service_name>.bats
```

## Manual Verification & Examples

Below is a list of all supported services, the command to verify them, and an example of a successful outcome.

| Service | CLI Command | Example Verification | Expected Outcome |
| :--- | :--- | :--- | :--- |
| ACM | `aws acm` | `aws acm <action>` | API response |
| APIGATEWAY | `aws apigateway` | `aws apigateway get-rest-apis` | Empty list or resource data |
| APIGATEWAYV2 | `aws apigatewayv2` | `aws apigatewayv2 get-apis` | Empty list or resource data |
| APPCONFIG | `aws appconfig` | `aws appconfig list-<resource>s` | Empty list or resource data |
| APPRUNNER | `aws apprunner` | `aws apprunner <action>` | API response |
| APPSYNC | `aws appsync` | `aws appsync list-<resource>s` | Empty list or resource data |
| ATHENA | `aws athena` | `aws athena list-<resource>s` | Empty list or resource data |
| AUTOSCALING | `aws autoscaling` | `aws autoscaling describe-<resource>s` | Resource metadata or identity info |
| BACKUP | `aws backup` | `aws backup list-<resource>s` | Empty list or resource data |
| BCMDATAEXPORTS | `aws bcmdataexports` | `aws bcmdataexports <action>` | API response |
| BEDROCKRUNTIME | `aws bedrock-runtime` | `aws bedrock-runtime <action>` | API response |
| CE | `aws ce` | `aws ce <action>` | API response |
| CLOUDFORMATION | `aws cloudformation` | `aws cloudformation describe-<resource>s` | Resource metadata or identity info |
| CLOUDFRONT | `aws cloudfront` | `aws cloudfront describe-<resource>s` | Resource metadata or identity info |
| CLOUDSTACK | `aws cloudstack` | `aws cloudstack <action>` | API response |
| CLOUDTRAIL | `aws cloudtrail` | `aws cloudtrail describe-<resource>s` | Resource metadata or identity info |
| CLOUDWATCH | `aws cloudwatch` | `aws cloudwatch describe-<resource>s` | Resource metadata or identity info |
| CODEBUILD | `aws codebuild` | `aws codebuild list-<resource>s` | Empty list or resource data |
| CODECOMMIT | `aws codecommit` | `aws codecommit list-<resource>s` | Empty list or resource data |
| CODEDEPLOY | `aws codedeploy` | `aws codedeploy list-<resource>s` | Empty list or resource data |
| CODEPIPELINE | `aws codepipeline` | `aws codepipeline list-<resource>s` | Empty list or resource data |
| COGNITO | `aws cognito` | `aws cognito <action>` | API response |
| CONFIGSERVICE | `aws configservice` | `aws configservice describe-<resource>s` | Resource metadata or identity info |
| CUR | `aws cur` | `aws cur describe-<resource>s` | Resource metadata or identity info |
| DOCDB | `aws docdb` | `aws docdb describe-<resource>s` | Resource metadata or identity info |
| DYNAMODB | `aws dynamodb` | `aws dynamodb list-tables` | List of tables |
| EC2 | `aws ec2` | `aws ec2 describe-<resource>s` | Resource metadata or identity info |
| ECR | `aws ecr` | `aws ecr list-<resource>s` | Empty list or resource data |
| ECS | `aws ecs` | `aws ecs list-<resource>s` | Empty list or resource data |
| EKS | `aws eks` | `aws eks list-<resource>s` | Empty list or resource data |
| ELASTICACHE | `aws elasticache` | `aws elasticache describe-<resource>s` | Resource metadata or identity info |
| ELBV2 | `aws elbv2` | `aws elbv2 describe-<resource>s` | Resource metadata or identity info |
| EMR | `aws emr` | `aws emr <action>` | API response |
| EVENTBRIDGE | `aws events` | `aws events <action>` | API response |
| FIREHOSE | `aws firehose` | `aws firehose list-<resource>s` | Empty list or resource data |
| GLUE | `aws glue` | `aws glue list-<resource>s` | Empty list or resource data |
| IAM | `aws iam` | `aws iam list-users` | List of users |
| KINESIS | `aws kinesis` | `aws kinesis list-<resource>s` | Empty list or resource data |
| KMS | `aws kms` | `aws kms <action>` | API response |
| LAMBDA | `aws lambda` | `aws lambda list-functions` | List of functions |
| MQ | `aws mq` | `aws mq list-brokers` | Empty list or resource data |
| MSK | `aws kafka` | `aws kafka <action>` | API response |
| NEPTUNE | `aws neptune` | `aws neptune <action>` | API response |
| OPENSEARCH | `aws opensearch` | `aws opensearch list-<resource>s` | Empty list or resource data |
| ORGANIZATIONS | `aws organizations` | `aws organizations describe-organization` | Resource metadata or identity info |
| PIPES | `aws pipes` | `aws pipes <action>` | API response |
| PRICING | `aws pricing` | `aws pricing <action>` | API response |
| RDS | `aws rds` | `aws rds describe-<resource>s` | Resource metadata or identity info |
| REDSHIFT | `aws redshift` | `aws redshift describe-clusters` | List of clusters |
| RESOURCEGROUPSTAGGING | `aws resourcegroupstaggingapi` | `aws resourcegroupstaggingapi <action>` | API response |
| ROUTE53 | `aws route53` | `aws route53 describe-<resource>s` | Resource metadata or identity info |
| S3 | `aws s3` | `aws s3 ls` | Empty list or list of buckets |
| SAGEMAKER | `aws sagemaker` | `aws sagemaker list-<resource>s` | Empty list or resource data |
| SCHEDULER | `aws scheduler` | `aws scheduler list-<resource>s` | Empty list or resource data |
| SECRETSMANAGER | `aws secretsmanager` | `aws secretsmanager list-<resource>s` | Empty list or resource data |
| SES | `aws ses` | `aws ses <action>` | API response |
| SNS | `aws sns` | `aws sns list-topics` | List of topic ARNs |
| SQS | `aws sqs` | `aws sqs list-queues` | List of queue URLs |
| SSM | `aws ssm` | `aws ssm list-<resource>s` | Empty list or resource data |
| STEPFUNCTIONS | `aws stepfunctions` | `aws stepfunctions list-<resource>s` | Empty list or resource data |
| STS | `aws sts` | `aws sts get-caller-identity` | Resource metadata or identity info |
| SYNTHETICS | `aws synthetics` | `aws synthetics <action>` | API response |
| TEXTRACT | `aws textract` | `aws textract <action>` | API response |
| TRANSCRIBE | `aws transcribe` | `aws transcribe <action>` | API response |
| TRANSFER | `aws transfer` | `aws transfer list-<resource>s` | Empty list or resource data |
| WAF | `aws waf` | `aws waf <action>` | API response |
| XRAY | `aws xray` | `aws xray <action>` | API response |
