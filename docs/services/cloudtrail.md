# CloudTrail

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateTrail` | Core implementation of CreateTrail |
| `DescribeTrails` | Core implementation of DescribeTrails |
| `DeleteTrail` | Core implementation of DeleteTrail |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws cloudtrail create-trail --name mytrail --s3-bucket-name mybucket
aws cloudtrail describe-trails
```
