# Redshift

**Protocol:** Query (XML)
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateCluster` | Core implementation of CreateCluster |
| `DescribeClusters` | Core implementation of DescribeClusters |
| `DeleteCluster` | Core implementation of DeleteCluster |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws redshift create-cluster --cluster-identifier mycluster --node-type dc2.large --master-username admin --db-name mydb
aws redshift describe-clusters --cluster-identifier mycluster
```
