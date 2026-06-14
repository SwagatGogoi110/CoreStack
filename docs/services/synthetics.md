# Synthetics

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateCanary` | Core implementation of CreateCanary |
| `GetCanary` | Core implementation of GetCanary |
| `DescribeCanaries` | Core implementation of DescribeCanaries |
| `DeleteCanary` | Core implementation of DeleteCanary |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws synthetics create-canary --name mycanary --artifact-s3-location s3://mybucket --execution-role-arn arn:aws:iam::000000000000:role/canary
aws synthetics describe-canaries
```
