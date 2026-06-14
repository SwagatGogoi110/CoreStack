# DocumentDB

**Protocol:** Query (XML)
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateDBCluster` | Core implementation of CreateDBCluster |
| `DescribeDBClusters` | Core implementation of DescribeDBClusters |
| `DeleteDBCluster` | Core implementation of DeleteDBCluster |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws docdb create-db-cluster --db-cluster-identifier mydbcluster --engine docdb
aws docdb describe-db-clusters
```
