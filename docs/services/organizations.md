# Organizations

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateOrganization` | Core implementation of CreateOrganization |
| `DescribeOrganization` | Core implementation of DescribeOrganization |
| `ListAccounts` | Core implementation of ListAccounts |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws organizations create-organization --feature-set ALL
aws organizations describe-organization
```
