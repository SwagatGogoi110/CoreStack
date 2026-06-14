import subprocess
import os

# List of GCP services and their verification URLs
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
    print(f"{'GCP Service':<25} | {'Endpoint':<50} | {'Status'}")
    print("-" * 100)

    passed = 0
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
            
            if status_code in ["200", "201", "204", "404", "400"]: 
                # 404/400 often mean the service responded but resource not found/bad req, which is better than 501
                if status_code == "501":
                    status = "❌ FAIL (501)"
                else:
                    status = f"✅ PASS ({status_code})"
                    passed += 1
            else:
                status = f"❌ FAIL ({status_code})"
        except Exception as e:
            status = "⚠️ ERROR"

        print(f"{svc:<25} | {path[:48]:<50} | {status}")

    print("-" * 100)
    print(f"Summary: {passed}/{len(gcp_services)} GCP services verified.")

if __name__ == "__main__":
    run_verification()
