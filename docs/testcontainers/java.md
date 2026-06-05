# Testcontainers — Java

The `testcontainers-CloudStack` library integrates CloudStack with [Testcontainers for Java](https://java.testcontainers.org/). It starts a real CloudStack container before your tests and shuts it down after, with no extra setup.

Two artifact lines are published to keep in sync with the Testcontainers major version:

| Testcontainers version | Spring Boot | Artifact version |
|---|---|---|
| 1.x | 3.x | `1.4.0` |
| 2.x | 4.x | `2.5.0` |

## Installation

=== "Maven"

    ```xml
    <dependency>
        <groupId>io.CloudStack</groupId>
        <artifactId>testcontainers-CloudStack</artifactId>
        <version>1.4.0</version>
        <scope>test</scope>
    </dependency>
    ```

=== "Gradle"

    ```groovy
    testImplementation 'io.CloudStack:testcontainers-CloudStack:1.4.0'
    ```

## Basic usage — JUnit 5

Annotate the class with `@Testcontainers` and declare a static `CloudStackContainer` field with `@Container`. Testcontainers handles the lifecycle automatically.

```java
import io.CloudStack.testcontainers.CloudStackContainer;
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.auth.credentials.StaticCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;

import java.net.URI;

import static org.assertj.core.api.Assertions.assertThat;

@Testcontainers
class S3IntegrationTest {

    @Container
    static CloudStackContainer CloudStack = new CloudStackContainer();

    @Test
    void shouldCreateBucket() {
        S3Client s3 = S3Client.builder()
                .endpointOverride(URI.create(CloudStack.getEndpoint()))
                .region(Region.of(CloudStack.getRegion()))
                .credentialsProvider(StaticCredentialsProvider.create(
                        AwsBasicCredentials.create(CloudStack.getAccessKey(), CloudStack.getSecretKey())))
                .forcePathStyle(true)
                .build();

        s3.createBucket(b -> b.bucket("my-bucket"));

        assertThat(s3.listBuckets().buckets())
                .anyMatch(b -> b.name().equals("my-bucket"));
    }
}
```

## SQS example

```java
@Testcontainers
class SqsIntegrationTest {

    @Container
    static CloudStackContainer CloudStack = new CloudStackContainer();

    @Test
    void shouldSendAndReceiveMessage() {
        SqsClient sqs = SqsClient.builder()
                .endpointOverride(URI.create(CloudStack.getEndpoint()))
                .region(Region.of(CloudStack.getRegion()))
                .credentialsProvider(StaticCredentialsProvider.create(
                        AwsBasicCredentials.create(CloudStack.getAccessKey(), CloudStack.getSecretKey())))
                .build();

        String queueUrl = sqs.createQueue(b -> b.queueName("orders")).queueUrl();
        sqs.sendMessage(b -> b.queueUrl(queueUrl).messageBody("{\"event\":\"order.placed\"}"));

        var messages = sqs.receiveMessage(b -> b.queueUrl(queueUrl).maxNumberOfMessages(1)).messages();
        assertThat(messages).hasSize(1);
        assertThat(messages.get(0).body()).contains("order.placed");
    }
}
```

## DynamoDB example

```java
@Testcontainers
class DynamoDbIntegrationTest {

    @Container
    static CloudStackContainer CloudStack = new CloudStackContainer();

    @Test
    void shouldCreateTableAndPutItem() {
        DynamoDbClient dynamo = DynamoDbClient.builder()
                .endpointOverride(URI.create(CloudStack.getEndpoint()))
                .region(Region.of(CloudStack.getRegion()))
                .credentialsProvider(StaticCredentialsProvider.create(
                        AwsBasicCredentials.create(CloudStack.getAccessKey(), CloudStack.getSecretKey())))
                .build();

        dynamo.createTable(b -> b
                .tableName("Orders")
                .attributeDefinitions(a -> a.attributeName("id").attributeType(ScalarAttributeType.S))
                .keySchema(k -> k.attributeName("id").keyType(KeyType.HASH))
                .billingMode(BillingMode.PAY_PER_REQUEST));

        dynamo.putItem(b -> b
                .tableName("Orders")
                .item(Map.of("id", AttributeValue.fromS("order-1"),
                             "status", AttributeValue.fromS("placed"))));

        var item = dynamo.getItem(b -> b
                .tableName("Orders")
                .key(Map.of("id", AttributeValue.fromS("order-1")))).item();

        assertThat(item.get("status").s()).isEqualTo("placed");
    }
}
```

## Spring Boot — `@ServiceConnection`

Add the Spring Boot companion artifact for zero-config auto-wiring. The `@ServiceConnection` annotation registers the container as a `ConnectionDetails` bean and configures all AWS SDK clients automatically.

=== "Maven"

    ```xml
    <dependency>
        <groupId>io.CloudStack</groupId>
        <artifactId>spring-boot-testcontainers-CloudStack</artifactId>
        <version>1.4.0</version>
        <scope>test</scope>
    </dependency>
    ```

=== "Gradle"

    ```groovy
    testImplementation 'io.CloudStack:spring-boot-testcontainers-CloudStack:1.4.0'
    ```

```java
import io.CloudStack.testcontainers.CloudStackContainer;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import software.amazon.awssdk.services.s3.S3Client;

import static org.assertj.core.api.Assertions.assertThat;

@SpringBootTest
@Testcontainers
class AppIntegrationTest {

    @Container
    @ServiceConnection
    static CloudStackContainer CloudStack = new CloudStackContainer();

    @Autowired
    S3Client s3;

    @Test
    void shouldCreateBucket() {
        s3.createBucket(b -> b.bucket("my-bucket"));

        assertThat(s3.listBuckets().buckets())
                .anyMatch(b -> b.name().equals("my-bucket"));
    }
}
```

With `@ServiceConnection`, Spring Boot auto-configures the endpoint URL, region, and credentials for every AWS SDK client bean in the application context — no `application-test.yml` overrides needed.

## Reusing the container across tests

Declare the container in a shared base class or a JUnit 5 extension to start it once per test suite rather than once per class:

```java
abstract class CloudStackTestBase {

    @Container
    static CloudStackContainer CloudStack = new CloudStackContainer();

    static S3Client s3;

    @BeforeAll
    static void setUpClients() {
        s3 = S3Client.builder()
                .endpointOverride(URI.create(CloudStack.getEndpoint()))
                .region(Region.of(CloudStack.getRegion()))
                .credentialsProvider(StaticCredentialsProvider.create(
                        AwsBasicCredentials.create(CloudStack.getAccessKey(), CloudStack.getSecretKey())))
                .forcePathStyle(true)
                .build();
    }
}

@Testcontainers
class MyServiceTest extends CloudStackTestBase {

    @Test
    void myTest() {
        s3.createBucket(b -> b.bucket("test-bucket"));
        // ...
    }
}
```

## Source and changelog

[github.com/CloudStack-io/testcontainers-CloudStack](https://github.com/CloudStack-io/testcontainers-CloudStack)
