import subprocess
import os
import sys

# AWS Services and their verification commands
aws_services = [
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

# GCP Services and their verification URLs
gcp_services = [
    ("cloudfunctions", "v1/projects/test/locations/us-central1/functions"),
    ("cloudrun", "v1/projects/test/locations/us-central1/services"),
    ("datastore", "v1/projects/test:lookup"),
    ("firestore", "v1/projects/test/databases/(default)/documents/test-collection"),
    ("gcs", "storage/v1/b"),
    ("iam", "v1/projects/test/serviceAccounts"),
    ("kafka", "v1/projects/test/locations/us-central1/clusters"),
    ("operations", "v1/projects/test/locations/us-central1/operations"),
    ("pubsub", "v1/projects/test/topics"),
    ("secretmanager", "v1/projects/test/secrets"),
    ("tasks", "v2/projects/test/locations/us-central1/queues"),
    ("bigquery", "bigquery/v2/projects/test/datasets"),
    ("cloudsql", "sql/v1/projects/test/instances"),
    ("spanner", "v1/projects/test/instances"),
    ("bigtable", "v2/projects/test/instances"),
    ("artifactregistry", "v1/projects/test/locations/us-central1/repositories"),
    ("compute", "compute/v1/projects/test/zones/us-central1-a/instances"),
    ("gke", "v1/projects/test/locations/us-central1/clusters"),
    ("cloudbuild", "v1/projects/test/locations/us-central1/builds"),
    ("appengine", "v1/apps/test"),
    ("workflows", "v1/projects/test/locations/us-central1/workflows"),
    ("cloudscheduler", "v1/projects/test/locations/us-central1/jobs"),
    ("dns", "dns/v1/projects/test/managedZones"),
    ("armor", "compute/v1/projects/test/global/securityPolicies"),
    ("loadbalancing", "compute/v1/projects/test/global/forwardingRules"),
    ("cas", "v1/projects/test/locations/us-central1/caPools"),
    ("logging", "v1/projects/test/logs"),
    ("monitoring", "v3/projects/test/timeSeries"),
    ("trace", "v2/projects/test/traces:batchWrite")
]

def run_verification():
    os.environ["AWS_ACCESS_KEY_ID"] = "test"
    os.environ["AWS_SECRET_ACCESS_KEY"] = "test"
    os.environ["AWS_DEFAULT_REGION"] = "us-east-1"
    os.environ["AWS_ENDPOINT_URL"] = "http://localhost:4566"

    print("=" * 100)
    print(f"{'AWS Service':<25} | {'Command':<40} | {'Status':<10} | {'Error/Notes'}")
    print("-" * 100)

    aws_passed = 0
    for svc, cmd in aws_services:
        try:
            result = subprocess.run(
                f"aws --endpoint-url http://localhost:4566 --region us-east-1 {svc} {cmd}",
                shell=True, capture_output=True, text=True, timeout=10
            )
            if result.returncode == 0:
                status = "✅ PASS"
                note = ""
                aws_passed += 1
            else:
                status = "❌ FAIL"
                note = result.stderr.split('\n')[0][:50]
                if "UnknownOperationException" in result.stderr:
                    note = "UnknownOperationException (Stub)"
                elif "RepositoryNotFoundException" in result.stderr:
                    status = "✅ PASS"
                    note = "(Repo not found - API OK)"
                    aws_passed += 1
        except Exception as e:
            status = "⚠️ ERROR"
            note = str(e)

        print(f"{svc:<25} | {cmd[:40]:<40} | {status:<10} | {note}")

    print("\n" + "=" * 100)
    print(f"{'GCP Service':<25} | {'Endpoint':<50} | {'Status'}")
    print("-" * 100)

    gcp_passed = 0
    for svc, path in gcp_services:
        url = f"http://localhost:4566/{path}"
        method = "GET"
        if ":write" in path or ":lookup" in path or ":batchWrite" in path or ":executeSql" in path:
            method = "POST"
        
        cmd = f"curl -s -o /dev/null -w '%{{http_code}}' -X {method} '{url}'"
        if method == "POST":
            cmd += " -d '{}' -H 'Content-Type: application/json'"
            
        try:
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=5)
            status_code = result.stdout.strip()
            
            if status_code in ["200", "201", "204"]: 
                status = f"✅ PASS ({status_code})"
                gcp_passed += 1
            else:
                status = f"❌ FAIL ({status_code})"
        except Exception as e:
            status = "⚠️ ERROR"

        print(f"{svc:<25} | {path[:48]:<50} | {status}")

    print("=" * 100)
    print(f"Summary AWS: {aws_passed}/{len(aws_services)} passed.")
    print(f"Summary GCP: {gcp_passed}/{len(gcp_services)} passed.")
    print(f"Total: {aws_passed + gcp_passed}/{len(aws_services) + len(gcp_services)} services verified.")

if __name__ == "__main__":
    run_verification()
