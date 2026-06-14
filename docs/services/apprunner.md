# AppRunner

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateService` | Core implementation of CreateService |
| `DescribeService` | Core implementation of DescribeService |
| `ListServices` | Core implementation of ListServices |
| `DeleteService` | Core implementation of DeleteService |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws apprunner create-service --service-name myservice
aws apprunner list-services
```
