# WAF

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateWebACL` | Core implementation of CreateWebACL |
| `GetWebACL` | Core implementation of GetWebACL |
| `ListWebACLs` | Core implementation of ListWebACLs |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws waf create-web-acl --name myacl --metric-name myacl --default-action Type=ALLOW
aws waf list-web-acls
```
