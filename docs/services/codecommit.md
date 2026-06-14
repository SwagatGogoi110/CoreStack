# CodeCommit

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateRepository` | Core implementation of CreateRepository |
| `GetRepository` | Core implementation of GetRepository |
| `ListRepositories` | Core implementation of ListRepositories |
| `DeleteRepository` | Core implementation of DeleteRepository |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws codecommit create-repository --repository-name myrepo
aws codecommit list-repositories
```
