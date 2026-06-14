# EMR

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `RunJobFlow` | Core implementation of RunJobFlow |
| `DescribeCluster` | Core implementation of DescribeCluster |
| `ListClusters` | Core implementation of ListClusters |
| `TerminateJobFlows` | Core implementation of TerminateJobFlows |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws emr run-job-flow --name "My Cluster" --release-label emr-6.2.0
aws emr list-clusters
```
