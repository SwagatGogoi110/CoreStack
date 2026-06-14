# CodePipeline

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreatePipeline` | Core implementation of CreatePipeline |
| `GetPipeline` | Core implementation of GetPipeline |
| `ListPipelines` | Core implementation of ListPipelines |
| `DeletePipeline` | Core implementation of DeletePipeline |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws codepipeline create-pipeline --pipeline '{"name": "mypipe", "roleArn": "arn:aws:iam::000000000000:role/pipe", "stages": []}'
aws codepipeline list-pipelines
```
