# GCP Integration Plan

## Goal
Integrate Google Cloud Platform (GCP) services (currently implemented in Java in `~/Developer/floci-gcp`) into the `cloudstack` Go emulator, maintaining strict architectural separation between AWS, GCP, and upcoming Azure services.

## Phase 1: Architectural Separation
To support multi-cloud emulation natively, we must isolate existing AWS services and prepare the structure for GCP and Azure.

1.  **Refactor Directory Structure:**
    *   Move all current 66 AWS services from `internal/services/<service>` to `internal/services/aws/<service>`.
    *   Create directories for `internal/services/gcp` and `internal/services/azure`.
2.  **Update Core Dispatcher:**
    *   Update import paths in `internal/core/server.go` and `internal/core/catalog_init.go`.
    *   Introduce multi-cloud request routing in `dispatcher.go`. AWS requests typically rely on `X-Amz-*` headers or specific domains. GCP requests typically rely on JSON over HTTP (`https://<service>.googleapis.com/v1/...`).

## Phase 2: Scaffolding and Tracking GCP Services
The `floci-gcp` Java project currently implements the following services:
*   `cloudfunctions`
*   `cloudrun`
*   `datastore`
*   `firestore`
*   `gcs` (Google Cloud Storage)
*   `iam` (GCP IAM)
*   `kafka` (Managed Kafka)
*   `operations` (Long Running Operations)
*   `pubsub`
*   `secretmanager`
*   `tasks` (Cloud Tasks)

1.  **Tracker Creation:** Create a Python script (`scripts/dev-tools/gcp_service_audit.py`) to track implementation progress.
2.  **Scaffolding:** Update or create a new generator script to scaffold GCP services in `internal/services/gcp/`.

## Phase 3: Porting and Expanding Implementations
We have successfully ported the initial 11 services and expanded the GCP provider to 29 total functional services. All implementations utilize the `storage.Backend` for persistence.

### ✅ Completed Services (29)
*   **Infrastructure**: Compute Engine, GKE, Cloud Load Balancing, Cloud DNS, Cloud Armor.
*   **Databases**: Cloud SQL, Spanner, Bigtable, Firestore, Datastore, BigQuery.
*   **Serverless**: Cloud Functions, Cloud Run, App Engine.
*   **DevOps**: Cloud Build, Artifact Registry, Workflows, Cloud Scheduler.
*   **Messaging**: Pub/Sub, Cloud Tasks, Kafka.
*   **Observability**: Cloud Logging, Cloud Monitoring, Cloud Trace.
*   **Security & Ops**: IAM, Secret Manager, CAS, Operations.

## Phase 4: Verification & Performance
1.  **Strict Verification:** Implemented `verify_all.py` for health checks across all 95 AWS/GCP services.
2.  **Dispatcher Optimization:** Refined multi-cloud heuristic routing for high-fidelity REST-JSON emulation.
3.  **Persistence Audit:** Ensure all services correctly utilize WAL-backed storage for production-grade reliability.

