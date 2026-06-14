# SageMaker

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateNotebookInstance` | Core implementation of CreateNotebookInstance |
| `ListNotebookInstances` | Core implementation of ListNotebookInstances |
| `CreateModel` | Core implementation of CreateModel |
| `ListModels` | Core implementation of ListModels |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws sagemaker create-notebook-instance --notebook-instance-name mynotebook --instance-type ml.t2.medium
aws sagemaker list-notebook-instances
```
