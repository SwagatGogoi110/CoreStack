# Testcontainers — Python

The `testcontainers-CloudStack` package integrates CloudStack with [Testcontainers for Python](https://testcontainers-python.readthedocs.io/). It works as a context manager and integrates naturally with pytest fixtures.

## Installation

```sh
pip install testcontainers-CloudStack
```

```sh
# poetry
poetry add --group dev testcontainers-CloudStack

# uv
uv add --dev testcontainers-CloudStack
```

## Basic usage — context manager

```python
import boto3
from testcontainers_CloudStack import CloudStackContainer


def test_s3_create_bucket():
    with CloudStackContainer() as CloudStack:
        s3 = boto3.client(
            "s3",
            endpoint_url=CloudStack.get_endpoint(),
            region_name=CloudStack.get_region(),
            aws_access_key_id=CloudStack.get_access_key(),
            aws_secret_access_key=CloudStack.get_secret_key(),
        )

        s3.create_bucket(Bucket="my-bucket")

        buckets = [b["Name"] for b in s3.list_buckets()["Buckets"]]
        assert "my-bucket" in buckets
```

## Pytest fixture

Use a session-scoped fixture so the container starts once and is shared across all tests in the suite.

```python
import pytest
import boto3
from testcontainers_CloudStack import CloudStackContainer


@pytest.fixture(scope="session")
def CloudStack():
    with CloudStackContainer() as container:
        yield container


@pytest.fixture(scope="session")
def s3_client(CloudStack):
    return boto3.client(
        "s3",
        endpoint_url=CloudStack.get_endpoint(),
        region_name=CloudStack.get_region(),
        aws_access_key_id=CloudStack.get_access_key(),
        aws_secret_access_key=CloudStack.get_secret_key(),
    )


def test_create_bucket(s3_client):
    s3_client.create_bucket(Bucket="my-bucket")
    buckets = [b["Name"] for b in s3_client.list_buckets()["Buckets"]]
    assert "my-bucket" in buckets


def test_upload_object(s3_client):
    s3_client.create_bucket(Bucket="uploads")
    s3_client.put_object(Bucket="uploads", Key="hello.txt", Body=b"hello CloudStack")
    body = s3_client.get_object(Bucket="uploads", Key="hello.txt")["Body"].read()
    assert body == b"hello CloudStack"
```

## SQS example

```python
import pytest
import boto3
import json
from testcontainers_CloudStack import CloudStackContainer


@pytest.fixture(scope="session")
def CloudStack():
    with CloudStackContainer() as container:
        yield container


@pytest.fixture(scope="session")
def sqs_client(CloudStack):
    return boto3.client(
        "sqs",
        endpoint_url=CloudStack.get_endpoint(),
        region_name=CloudStack.get_region(),
        aws_access_key_id=CloudStack.get_access_key(),
        aws_secret_access_key=CloudStack.get_secret_key(),
    )


def test_send_and_receive_message(sqs_client):
    queue_url = sqs_client.create_queue(QueueName="orders")["QueueUrl"]

    sqs_client.send_message(
        QueueUrl=queue_url,
        MessageBody=json.dumps({"event": "order.placed"}),
    )

    response = sqs_client.receive_message(QueueUrl=queue_url, MaxNumberOfMessages=1)
    messages = response.get("Messages", [])

    assert len(messages) == 1
    assert json.loads(messages[0]["Body"])["event"] == "order.placed"
```

## DynamoDB example

```python
import pytest
import boto3
from testcontainers_CloudStack import CloudStackContainer


@pytest.fixture(scope="session")
def CloudStack():
    with CloudStackContainer() as container:
        yield container


@pytest.fixture(scope="session")
def dynamo_client(CloudStack):
    return boto3.client(
        "dynamodb",
        endpoint_url=CloudStack.get_endpoint(),
        region_name=CloudStack.get_region(),
        aws_access_key_id=CloudStack.get_access_key(),
        aws_secret_access_key=CloudStack.get_secret_key(),
    )


def test_put_and_get_item(dynamo_client):
    dynamo_client.create_table(
        TableName="Orders",
        AttributeDefinitions=[{"AttributeName": "id", "AttributeType": "S"}],
        KeySchema=[{"AttributeName": "id", "KeyType": "HASH"}],
        BillingMode="PAY_PER_REQUEST",
    )

    dynamo_client.put_item(
        TableName="Orders",
        Item={"id": {"S": "order-1"}, "status": {"S": "placed"}},
    )

    item = dynamo_client.get_item(
        TableName="Orders",
        Key={"id": {"S": "order-1"}},
    )["Item"]

    assert item["status"]["S"] == "placed"
```

## Secrets Manager example

```python
def test_create_and_get_secret(CloudStack):
    sm = boto3.client(
        "secretsmanager",
        endpoint_url=CloudStack.get_endpoint(),
        region_name=CloudStack.get_region(),
        aws_access_key_id=CloudStack.get_access_key(),
        aws_secret_access_key=CloudStack.get_secret_key(),
    )

    sm.create_secret(Name="db/password", SecretString="supersecret")
    value = sm.get_secret_value(SecretId="db/password")["SecretString"]
    assert value == "supersecret"
```

## conftest.py pattern

Place shared fixtures in `conftest.py` so every test module picks them up automatically:

```python
# conftest.py
import pytest
import boto3
from testcontainers_CloudStack import CloudStackContainer


@pytest.fixture(scope="session")
def CloudStack():
    with CloudStackContainer() as container:
        yield container


@pytest.fixture(scope="session")
def aws_clients(CloudStack):
    kwargs = dict(
        endpoint_url=CloudStack.get_endpoint(),
        region_name=CloudStack.get_region(),
        aws_access_key_id=CloudStack.get_access_key(),
        aws_secret_access_key=CloudStack.get_secret_key(),
    )
    return {
        "s3": boto3.client("s3", **kwargs),
        "sqs": boto3.client("sqs", **kwargs),
        "dynamodb": boto3.client("dynamodb", **kwargs),
        "secretsmanager": boto3.client("secretsmanager", **kwargs),
    }
```

```python
# test_my_service.py
def test_something(aws_clients):
    s3 = aws_clients["s3"]
    s3.create_bucket(Bucket="test")
    # ...
```

## Source and changelog

[github.com/CloudStack-io/testcontainers-CloudStack-python](https://github.com/CloudStack-io/testcontainers-CloudStack-python)
