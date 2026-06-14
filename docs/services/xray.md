# X-Ray

**Protocol:** JSON 1.1
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `PutTraceSegments` | Core implementation of PutTraceSegments |
| `GetTraceSummaries` | Core implementation of GetTraceSummaries |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws xray put-trace-segments --trace-segment-documents '{"id": "1", "name": "myservice"}'
aws xray get-trace-summaries --start-time 1 --end-time 2
```
