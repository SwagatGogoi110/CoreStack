# Testcontainers — Node.js / TypeScript

The `@CloudStack/testcontainers` package integrates CloudStack with [Testcontainers for Node.js](https://node.testcontainers.org/). It works with any test runner that supports `async`/`await` — Jest, Vitest, Mocha, and others.

## Installation

```sh
npm install --save-dev @CloudStack/testcontainers
```

```sh
# yarn
yarn add --dev @CloudStack/testcontainers

# pnpm
pnpm add -D @CloudStack/testcontainers
```

## Basic usage — Jest

```typescript
import { CloudStackContainer } from "@CloudStack/testcontainers";
import { S3Client, CreateBucketCommand, ListBucketsCommand } from "@aws-sdk/client-s3";

describe("S3", () => {
    let CloudStack: CloudStackContainer;

    beforeAll(async () => {
        CloudStack = await new CloudStackContainer().start();
    });

    afterAll(async () => {
        await CloudStack.stop();
    });

    it("should create and list a bucket", async () => {
        const s3 = new S3Client({
            endpoint: CloudStack.getEndpoint(),
            region: CloudStack.getRegion(),
            credentials: {
                accessKeyId: CloudStack.getAccessKey(),
                secretAccessKey: CloudStack.getSecretKey(),
            },
            forcePathStyle: true,
        });

        await s3.send(new CreateBucketCommand({ Bucket: "my-bucket" }));

        const { Buckets } = await s3.send(new ListBucketsCommand({}));
        expect(Buckets?.some(b => b.Name === "my-bucket")).toBe(true);
    });
});
```

## SQS example

```typescript
import { CloudStackContainer } from "@CloudStack/testcontainers";
import {
    SQSClient,
    CreateQueueCommand,
    SendMessageCommand,
    ReceiveMessageCommand,
} from "@aws-sdk/client-sqs";

describe("SQS", () => {
    let CloudStack: CloudStackContainer;
    let sqs: SQSClient;

    beforeAll(async () => {
        CloudStack = await new CloudStackContainer().start();
        sqs = new SQSClient({
            endpoint: CloudStack.getEndpoint(),
            region: CloudStack.getRegion(),
            credentials: {
                accessKeyId: CloudStack.getAccessKey(),
                secretAccessKey: CloudStack.getSecretKey(),
            },
        });
    });

    afterAll(async () => {
        await CloudStack.stop();
    });

    it("should send and receive a message", async () => {
        const { QueueUrl } = await sqs.send(
            new CreateQueueCommand({ QueueName: "orders" })
        );

        await sqs.send(
            new SendMessageCommand({
                QueueUrl,
                MessageBody: JSON.stringify({ event: "order.placed" }),
            })
        );

        const { Messages } = await sqs.send(
            new ReceiveMessageCommand({ QueueUrl, MaxNumberOfMessages: 1 })
        );

        expect(Messages).toHaveLength(1);
        expect(JSON.parse(Messages![0].Body!).event).toBe("order.placed");
    });
});
```

## DynamoDB example

```typescript
import { CloudStackContainer } from "@CloudStack/testcontainers";
import {
    DynamoDBClient,
    CreateTableCommand,
    PutItemCommand,
    GetItemCommand,
} from "@aws-sdk/client-dynamodb";

describe("DynamoDB", () => {
    let CloudStack: CloudStackContainer;
    let dynamo: DynamoDBClient;

    beforeAll(async () => {
        CloudStack = await new CloudStackContainer().start();
        dynamo = new DynamoDBClient({
            endpoint: CloudStack.getEndpoint(),
            region: CloudStack.getRegion(),
            credentials: {
                accessKeyId: CloudStack.getAccessKey(),
                secretAccessKey: CloudStack.getSecretKey(),
            },
        });
    });

    afterAll(async () => {
        await CloudStack.stop();
    });

    it("should put and get an item", async () => {
        await dynamo.send(
            new CreateTableCommand({
                TableName: "Orders",
                AttributeDefinitions: [{ AttributeName: "id", AttributeType: "S" }],
                KeySchema: [{ AttributeName: "id", KeyType: "HASH" }],
                BillingMode: "PAY_PER_REQUEST",
            })
        );

        await dynamo.send(
            new PutItemCommand({
                TableName: "Orders",
                Item: {
                    id: { S: "order-1" },
                    status: { S: "placed" },
                },
            })
        );

        const { Item } = await dynamo.send(
            new GetItemCommand({
                TableName: "Orders",
                Key: { id: { S: "order-1" } },
            })
        );

        expect(Item?.status?.S).toBe("placed");
    });
});
```

## Vitest

The same pattern works with Vitest — replace `describe`/`it`/`expect` with their Vitest equivalents (the API is identical):

```typescript
import { describe, it, expect, beforeAll, afterAll } from "vitest";
import { CloudStackContainer } from "@CloudStack/testcontainers";
import { S3Client, CreateBucketCommand, ListBucketsCommand } from "@aws-sdk/client-s3";

describe("S3", () => {
    let CloudStack: CloudStackContainer;

    beforeAll(async () => {
        CloudStack = await new CloudStackContainer().start();
    });

    afterAll(async () => {
        await CloudStack.stop();
    });

    it("should create a bucket", async () => {
        const s3 = new S3Client({
            endpoint: CloudStack.getEndpoint(),
            region: CloudStack.getRegion(),
            credentials: {
                accessKeyId: CloudStack.getAccessKey(),
                secretAccessKey: CloudStack.getSecretKey(),
            },
            forcePathStyle: true,
        });

        await s3.send(new CreateBucketCommand({ Bucket: "vitest-bucket" }));

        const { Buckets } = await s3.send(new ListBucketsCommand({}));
        expect(Buckets?.some(b => b.Name === "vitest-bucket")).toBe(true);
    });
});
```

## Reusing the container across test files

Start the container once in a global setup file and expose the endpoint via an environment variable or a shared module so individual test files don't each start their own container.

=== "Jest — globalSetup"

    ```typescript
    // jest.global-setup.ts
    import { CloudStackContainer } from "@CloudStack/testcontainers";

    let CloudStack: CloudStackContainer;

    export async function setup() {
        CloudStack = await new CloudStackContainer().start();
        process.env.CLOUDSTACK_ENDPOINT = CloudStack.getEndpoint();
    }

    export async function teardown() {
        await CloudStack?.stop();
    }
    ```

    ```json
    // jest.config.json
    {
      "globalSetup": "./jest.global-setup.ts"
    }
    ```

=== "Vitest — globalSetup"

    ```typescript
    // vitest.global-setup.ts
    import { CloudStackContainer } from "@CloudStack/testcontainers";

    let CloudStack: CloudStackContainer;

    export async function setup() {
        CloudStack = await new CloudStackContainer().start();
        process.env.CLOUDSTACK_ENDPOINT = CloudStack.getEndpoint();
    }

    export async function teardown() {
        await CloudStack?.stop();
    }
    ```

    ```typescript
    // vitest.config.ts
    import { defineConfig } from "vitest/config";

    export default defineConfig({
        test: {
            globalSetup: "./vitest.global-setup.ts",
        },
    });
    ```

## Source and changelog

[github.com/CloudStack-io/testcontainers-CloudStack-node](https://github.com/CloudStack-io/testcontainers-CloudStack-node)
