# CloudFormation Service

## Overview
The CloudFormation service emulates the creation and management of infrastructure stacks.

## Structure
`internal/services/cloudformation/`
- `model/models.go`: Defines `Stack` and `Parameter`.
- `service.go`: Manages stack lifecycle and state transitions.
- `handler.go`: Implements the `AwsQuery` protocol.

## Working Logic
1.  **Stack Management:** Currently supports basic `CreateStack` and `DescribeStacks` operations.
2.  **Persistence:** Stack state is stored in `WalStorage`.

## Architecture
```mermaid
sequenceDiagram
    participant CLI as AWS CLI
    participant Hand as CFHandler
    participant Svc as CFService

    CLI->>Hand: CreateStack
    Hand->>Svc: CreateStack(name, template)
    Svc->>Svc: Parse Template (Stub)
    Svc-->>Hand: Stack ID
    Hand-->>CLI: 200 OK
```
