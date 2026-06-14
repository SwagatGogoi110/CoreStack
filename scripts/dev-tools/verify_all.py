import subprocess
import os
import sys

# List of services and their verification commands
services = [
    ("acm", "list-certificates"),
    ("apigateway", "get-rest-apis"),
    ("apigatewayv2", "get-apis"),
    ("appconfig", "list-applications"),
    ("apprunner", "list-services"),
    ("appsync", "list-graphql-apis"),
    ("athena", "list-work-groups"),
    ("autoscaling", "describe-auto-scaling-groups"),
    ("backup", "list-backup-plans"),
    ("bcm-data-exports", "list-exports"),
    ("bedrock-runtime", "list-async-invokes"),
    ("ce", "get-cost-and-usage --time-period Start=2024-01-01,End=2024-01-02 --granularity DAILY --metrics UsageQuantity"),
    ("cloudformation", "describe-stacks"),
    ("cloudfront", "list-distributions"),
    ("cloudtrail", "describe-trails"),
    ("cloudwatch", "list-metrics"),
    ("codebuild", "list-projects"),
    ("codecommit", "list-repositories"),
    ("deploy", "list-applications"),
    ("codepipeline", "list-pipelines"),
    ("cognito-idp", "list-user-pools --max-results 1"),
    ("configservice", "describe-config-rules"),
    ("cur", "describe-report-definitions"),
    ("docdb", "describe-db-clusters"),
    ("dynamodb", "list-tables"),
    ("ec2", "describe-instances"),
    ("ecr", "describe-repositories"),
    ("ecs", "list-clusters"),
    ("eks", "list-clusters"),
    ("elasticache", "describe-cache-clusters"),
    ("elbv2", "describe-load-balancers"),
    ("emr", "list-clusters"),
    ("events", "list-event-buses"),
    ("firehose", "list-delivery-streams"),
    ("glue", "get-databases"),
    ("iam", "list-users"),
    ("kinesis", "list-streams"),
    ("kms", "list-keys"),
    ("lambda", "list-functions"),
    ("mq", "list-brokers"),
    ("kafka", "list-clusters"),
    ("neptune", "describe-db-instances"),
    ("opensearch", "list-domain-names"),
    ("organizations", "describe-organization"),
    ("pipes", "list-pipes"),
    ("pricing", "describe-services --service-code AmazonEC2"),
    ("rds", "describe-db-instances"),
    ("redshift", "describe-clusters"),
    ("resourcegroupstaggingapi", "get-resources"),
    ("route53", "list-hosted-zones"),
    ("s3", "ls"),
    ("sagemaker", "list-notebook-instances"),
    ("scheduler", "list-schedules"),
    ("secretsmanager", "list-secrets"),
    ("ses", "list-identities"),
    ("sns", "list-topics"),
    ("sqs", "list-queues"),
    ("ssm", "list-commands"),
    ("stepfunctions", "list-state-machines"),
    ("sts", "get-caller-identity"),
    ("synthetics", "describe-canaries"),
    ("textract", "list-adapters"),
    ("transcribe", "list-transcription-jobs"),
    ("transfer", "list-servers"),
    ("waf", "list-web-acls"),
    ("xray", "get-trace-summaries --start-time 1 --end-time 2")
]

def run_verification():
    os.environ["AWS_ACCESS_KEY_ID"] = "test"
    os.environ["AWS_SECRET_ACCESS_KEY"] = "test"
    os.environ["AWS_DEFAULT_REGION"] = "us-east-1"
    os.environ["AWS_ENDPOINT_URL"] = "http://localhost:4566"

    print(f"{'Service':<25} | {'Command':<40} | {'Status':<10} | {'Error/Notes'}")
    print("-" * 100)

    results = []

    for svc, cmd in services:
        try:
            # Using --endpoint-url explicitly and --region
            result = subprocess.run(
                f"aws --endpoint-url http://localhost:4566 --region us-east-1 {svc} {cmd}",
                shell=True, capture_output=True, text=True, timeout=10
            )
            if result.returncode == 0:
                status = "✅ PASS"
                note = ""
            else:
                status = "❌ FAIL"
                note = result.stderr.split('\n')[0][:50]
                if "UnknownOperationException" in result.stderr:
                    note = "UnknownOperationException (Stub)"
                elif "RepositoryNotFoundException" in result.stderr:
                    status = "✅ PASS"
                    note = "(Repo not found - API OK)"
        except Exception as e:
            status = "⚠️ ERROR"
            note = str(e)

        print(f"{svc:<25} | {cmd[:40]:<40} | {status:<10} | {note}")
        results.append((svc, status))

    passed = len([r for r in results if "PASS" in r[1]])
    print("-" * 100)
    print(f"Summary: {passed}/{len(services)} services responded correctly.")

if __name__ == "__main__":
    run_verification()
