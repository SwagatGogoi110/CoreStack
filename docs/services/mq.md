# Amazon MQ

**Protocol:** REST-JSON
**Management Endpoint:** `POST http://localhost:4566/`

## Supported Management Actions

| Action | Description |
|---|---|
| `CreateBroker` | Core implementation of CreateBroker |
| `DescribeBroker` | Core implementation of DescribeBroker |
| `ListBrokers` | Core implementation of ListBrokers |
| `DeleteBroker` | Core implementation of DeleteBroker |

## Examples

```bash
export AWS_ENDPOINT_URL=http://localhost:4566

aws mq create-broker --broker-name mybroker --engine-type ACTIVEMQ --deployment-mode SINGLE_INSTANCE
aws mq list-brokers
```
