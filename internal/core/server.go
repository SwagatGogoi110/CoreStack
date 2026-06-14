package core

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hectorvent/cloudstack/internal/core/common"
	"github.com/hectorvent/cloudstack/internal/services/aws/acm"
	"github.com/hectorvent/cloudstack/internal/services/aws/apigateway"
	"github.com/hectorvent/cloudstack/internal/services/aws/apigatewayv2"
	"github.com/hectorvent/cloudstack/internal/services/aws/appconfig"
	"github.com/hectorvent/cloudstack/internal/services/aws/appsync"
	"github.com/hectorvent/cloudstack/internal/services/aws/athena"
	"github.com/hectorvent/cloudstack/internal/services/aws/autoscaling"
	"github.com/hectorvent/cloudstack/internal/services/aws/backup"
	"github.com/hectorvent/cloudstack/internal/services/aws/bcmdataexports"
	"github.com/hectorvent/cloudstack/internal/services/aws/bedrockruntime"
	"github.com/hectorvent/cloudstack/internal/services/aws/ce"
	"github.com/hectorvent/cloudstack/internal/services/aws/cloudformation"
	"github.com/hectorvent/cloudstack/internal/services/aws/cloudfront"
	"github.com/hectorvent/cloudstack/internal/services/aws/cloudtrail"
	"github.com/hectorvent/cloudstack/internal/services/aws/cloudwatch"
	"github.com/hectorvent/cloudstack/internal/services/aws/codebuild"
	"github.com/hectorvent/cloudstack/internal/services/aws/codecommit"
	"github.com/hectorvent/cloudstack/internal/services/aws/codedeploy"
	"github.com/hectorvent/cloudstack/internal/services/aws/codepipeline"
	"github.com/hectorvent/cloudstack/internal/services/aws/cognito"
	"github.com/hectorvent/cloudstack/internal/services/aws/configservice"
	"github.com/hectorvent/cloudstack/internal/services/aws/cur"
	"github.com/hectorvent/cloudstack/internal/services/aws/docdb"
	"github.com/hectorvent/cloudstack/internal/services/aws/dynamodb"
	"github.com/hectorvent/cloudstack/internal/services/aws/ec2"
	"github.com/hectorvent/cloudstack/internal/services/aws/ecr"
	"github.com/hectorvent/cloudstack/internal/services/aws/ecs"
	"github.com/hectorvent/cloudstack/internal/services/aws/eks"
	"github.com/hectorvent/cloudstack/internal/services/aws/elasticache"
	"github.com/hectorvent/cloudstack/internal/services/aws/elbv2"
	"github.com/hectorvent/cloudstack/internal/services/aws/emr"
	"github.com/hectorvent/cloudstack/internal/services/aws/eventbridge"
	"github.com/hectorvent/cloudstack/internal/services/aws/firehose"
	"github.com/hectorvent/cloudstack/internal/services/aws/glue"
	"github.com/hectorvent/cloudstack/internal/services/aws/iam"
	"github.com/hectorvent/cloudstack/internal/services/aws/kinesis"
	"github.com/hectorvent/cloudstack/internal/services/aws/kms"
	"github.com/hectorvent/cloudstack/internal/services/aws/lambda"
	"github.com/hectorvent/cloudstack/internal/services/aws/mq"
	"github.com/hectorvent/cloudstack/internal/services/aws/msk"
	"github.com/hectorvent/cloudstack/internal/services/aws/neptune"
	"github.com/hectorvent/cloudstack/internal/services/aws/opensearch"
	"github.com/hectorvent/cloudstack/internal/services/aws/organizations"
	"github.com/hectorvent/cloudstack/internal/services/aws/pipes"
	"github.com/hectorvent/cloudstack/internal/services/aws/pricing"
	"github.com/hectorvent/cloudstack/internal/services/aws/rds"
	"github.com/hectorvent/cloudstack/internal/services/aws/redshift"
	"github.com/hectorvent/cloudstack/internal/services/aws/resourcegroupstagging"
	"github.com/hectorvent/cloudstack/internal/services/aws/route53"
	"github.com/hectorvent/cloudstack/internal/services/aws/s3"
	"github.com/hectorvent/cloudstack/internal/services/aws/sagemaker"
	"github.com/hectorvent/cloudstack/internal/services/aws/scheduler"
	"github.com/hectorvent/cloudstack/internal/services/aws/secretsmanager"
	"github.com/hectorvent/cloudstack/internal/services/aws/ses"
	"github.com/hectorvent/cloudstack/internal/services/aws/sns"
	"github.com/hectorvent/cloudstack/internal/services/aws/sqs"
	"github.com/hectorvent/cloudstack/internal/services/aws/ssm"
	"github.com/hectorvent/cloudstack/internal/services/aws/stepfunctions"
	"github.com/hectorvent/cloudstack/internal/services/aws/sts"
	"github.com/hectorvent/cloudstack/internal/services/aws/synthetics"
	"github.com/hectorvent/cloudstack/internal/services/aws/textract"
	"github.com/hectorvent/cloudstack/internal/services/aws/transcribe"
	"github.com/hectorvent/cloudstack/internal/services/aws/transfer"
	"github.com/hectorvent/cloudstack/internal/services/aws/waf"
	"github.com/hectorvent/cloudstack/internal/services/aws/xray"
	"github.com/hectorvent/cloudstack/internal/services/aws/apprunner"
	"github.com/hectorvent/cloudstack/internal/services/aws/cloudstack"
	"github.com/hectorvent/cloudstack/internal/storage"

	// GCP Services
	gcp_cf "github.com/hectorvent/cloudstack/internal/services/gcp/cloudfunctions"
	gcp_cr "github.com/hectorvent/cloudstack/internal/services/gcp/cloudrun"
	gcp_ds "github.com/hectorvent/cloudstack/internal/services/gcp/datastore"
	gcp_fs "github.com/hectorvent/cloudstack/internal/services/gcp/firestore"
	gcp_gcs "github.com/hectorvent/cloudstack/internal/services/gcp/gcs"
	gcp_iam "github.com/hectorvent/cloudstack/internal/services/gcp/iam"
	gcp_kafka "github.com/hectorvent/cloudstack/internal/services/gcp/kafka"
	gcp_ops "github.com/hectorvent/cloudstack/internal/services/gcp/operations"
	gcp_pubsub "github.com/hectorvent/cloudstack/internal/services/gcp/pubsub"
	gcp_sm "github.com/hectorvent/cloudstack/internal/services/gcp/secretmanager"
	gcp_tasks "github.com/hectorvent/cloudstack/internal/services/gcp/tasks"
	gcp_bq "github.com/hectorvent/cloudstack/internal/services/gcp/bigquery"
	gcp_sql "github.com/hectorvent/cloudstack/internal/services/gcp/cloudsql"
	gcp_spanner "github.com/hectorvent/cloudstack/internal/services/gcp/spanner"
	gcp_bt "github.com/hectorvent/cloudstack/internal/services/gcp/bigtable"
	gcp_ar "github.com/hectorvent/cloudstack/internal/services/gcp/artifactregistry"
	gcp_ce "github.com/hectorvent/cloudstack/internal/services/gcp/compute"
	gcp_gke "github.com/hectorvent/cloudstack/internal/services/gcp/container"
	gcp_cb "github.com/hectorvent/cloudstack/internal/services/gcp/cloudbuild"
	gcp_ae "github.com/hectorvent/cloudstack/internal/services/gcp/appengine"
	gcp_wf "github.com/hectorvent/cloudstack/internal/services/gcp/workflows"
	gcp_cs "github.com/hectorvent/cloudstack/internal/services/gcp/cloudscheduler"
	gcp_dns "github.com/hectorvent/cloudstack/internal/services/gcp/dns"
	gcp_armor "github.com/hectorvent/cloudstack/internal/services/gcp/armor"
	gcp_lb "github.com/hectorvent/cloudstack/internal/services/gcp/loadbalancing"
	gcp_cas "github.com/hectorvent/cloudstack/internal/services/gcp/privateca"
	gcp_logging "github.com/hectorvent/cloudstack/internal/services/gcp/logging"
	gcp_monitoring "github.com/hectorvent/cloudstack/internal/services/gcp/monitoring"
	gcp_trace "github.com/hectorvent/cloudstack/internal/services/gcp/trace"
)

// Server represents the core CloudStack HTTP server, routing to various AWS services.
type Server struct {
	mux            *http.ServeMux
	resolver       *Resolver
	catalog        *common.Catalog
	handlers       map[string]common.ServiceHandler
	storage        *storage.Factory
	s3             *s3.S3Handler
	lambda         *lambda.LambdaHandler
	apigateway     *apigateway.ApiGatewayHandler
	bedrockruntime *bedrockruntime.BedrockRuntimeHandler
	cloudfront     *cloudfront.CloudFrontHandler
	opensearch     *opensearch.OpenSearchJsonHandler
	route53        *route53.Route53Handler
	mq             *mq.MqHandler
	duck           *cloudstack.DuckManager
	gcpHandlers    map[string]http.Handler
}

// NewServer initializes and returns a new Server instance.
func NewServer() *Server {
	mux := http.NewServeMux()

	storageFactory := storage.NewFactory("./data", "000000000000", 5*time.Second)
	duckManager := cloudstack.NewDuckManager()

	s := &Server{
		mux:         mux,
		resolver:    NewResolver("", ""), // TODO: Load from config
		catalog:     initCatalog(),
		handlers:    make(map[string]common.ServiceHandler),
		gcpHandlers: make(map[string]http.Handler),
		storage:     storageFactory,
		duck:        duckManager,
	}


	// Initialize IAM
	iamService, err := iam.NewIamService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize IAM: %v", err)
	}
	s.handlers["iam"] = iam.NewIamQueryHandler(iamService)
	s.handlers["sts"] = sts.NewStsQueryHandler(iamService)

	// Initialize Cognito
	cognitoService, err := cognito.NewCognitoService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Cognito: %v", err)
	}
	s.handlers["cognito-idp"] = cognito.NewCognitoJsonHandler(cognitoService)

	// Initialize DynamoDB
	dynamoDbService, err := dynamodb.NewDynamoDbService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB: %v", err)
	}
	s.handlers["dynamodb"] = dynamodb.NewDynamoDbHandler(dynamoDbService)

	// Initialize S3
	s3Service, err := s3.NewS3Service(storageFactory, "./data/s3", false)
	if err != nil {
		log.Fatalf("Failed to initialize S3: %v", err)
	}
	s.s3 = s3.NewS3Handler(s3Service)

	// Initialize Lambda
	lambdaService, err := lambda.NewLambdaService(storageFactory, "./data/lambda")
	if err != nil {
		log.Fatalf("Failed to initialize Lambda: %v", err)
	}
	s.lambda = lambda.NewLambdaHandler(lambdaService)

	// Initialize API Gateway
	apigwService, err := apigateway.NewApiGatewayService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize API Gateway: %v", err)
	}
	s.apigateway = apigateway.NewApiGatewayHandler(apigwService)

	// Initialize API Gateway V2
	apigw2Service, err := apigatewayv2.NewApiGatewayV2Service(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize API Gateway V2: %v", err)
	}
	s.handlers["apigatewayv2"] = apigatewayv2.NewApiGatewayV2JsonHandler(apigw2Service)

	// Initialize AppConfig
	appconfigService, err := appconfig.NewAppConfigService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize AppConfig: %v", err)
	}
	s.handlers["appconfig"] = appconfig.NewAppConfigJsonHandler(appconfigService)

	// Initialize AppSync
	appsyncService, err := appsync.NewAppSyncService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize AppSync: %v", err)
	}
	s.handlers["appsync"] = appsync.NewAppSyncJsonHandler(appsyncService)

	// Initialize Athena
	athenaService, err := athena.NewAthenaService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Athena: %v", err)
	}
	s.handlers["athena"] = athena.NewAthenaJsonHandler(athenaService)

	// Initialize Autoscaling
	asgService, err := autoscaling.NewAutoScalingService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Autoscaling: %v", err)
	}
	s.handlers["autoscaling"] = autoscaling.NewAutoScalingQueryHandler(asgService)

	// Initialize Backup
	backupService, err := backup.NewBackupService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Backup: %v", err)
	}
	s.handlers["backup"] = backup.NewBackupJsonHandler(backupService)

	// Initialize BCM Data Exports
	bcmService, err := bcmdataexports.NewBcmDataExportsService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize BCM Data Exports: %v", err)
	}
	s.handlers["bcmdataexports"] = bcmdataexports.NewBcmDataExportsJsonHandler(bcmService)

	// Initialize Bedrock Runtime
	s.bedrockruntime = bedrockruntime.NewBedrockRuntimeHandler(bedrockruntime.NewBedrockRuntimeService())
	s.handlers["bedrock-runtime"] = s.bedrockruntime

	// Initialize Cost Explorer
	s.handlers["ce"] = ce.NewCostExplorerJsonHandler(ce.NewCostExplorerService())

	// Initialize CloudFormation
	cfService, err := cloudformation.NewCloudFormationService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize CloudFormation: %v", err)
	}
	s.handlers["cloudformation"] = cloudformation.NewCloudFormationQueryHandler(cfService)

	// Initialize CloudFront
	cfDistService, err := cloudfront.NewCloudFrontService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize CloudFront: %v", err)
	}
	s.cloudfront = cloudfront.NewCloudFrontHandler(cfDistService)

	// Initialize CloudWatch
	cwService, err := cloudwatch.NewCloudWatchService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize CloudWatch: %v", err)
	}
	s.handlers["monitoring"] = cloudwatch.NewCloudWatchQueryHandler(cwService)

	// Initialize CodeBuild
	cbService, err := codebuild.NewCodeBuildService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize CodeBuild: %v", err)
	}
	s.handlers["codebuild"] = codebuild.NewCodeBuildJsonHandler(cbService)

	// Initialize CodeDeploy
	cdService, err := codedeploy.NewCodeDeployService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize CodeDeploy: %v", err)
	}
	s.handlers["codedeploy"] = codedeploy.NewCodeDeployJsonHandler(cdService)

	// Initialize Config
	cfgService, err := configservice.NewConfigService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Config: %v", err)
	}
	s.handlers["config"] = configservice.NewConfigServiceJsonHandler(cfgService)

	// Initialize CUR
	curService, err := cur.NewCurService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize CUR: %v", err)
	}
	s.handlers["cur"] = cur.NewCurJsonHandler(curService)

	// Initialize EC2
	ec2Service, err := ec2.NewEc2Service(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize EC2: %v", err)
	}
	s.handlers["ec2"] = ec2.NewEc2QueryHandler(ec2Service)

	// Initialize ECR
	ecrService, err := ecr.NewEcrService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize ECR: %v", err)
	}
	s.handlers["ecr"] = ecr.NewEcrJsonHandler(ecrService)

	// Initialize ECS
	ecsService, err := ecs.NewEcsService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize ECS: %v", err)
	}
	s.handlers["ecs"] = ecs.NewEcsJsonHandler(ecsService)

	// Initialize EKS
	eksService, err := eks.NewEksService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize EKS: %v", err)
	}
	s.handlers["eks"] = eks.NewEksJsonHandler(eksService)

	// Initialize ElastiCache
	ecService, err := elasticache.NewElastiCacheService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize ElastiCache: %v", err)
	}
	s.handlers["elasticache"] = elasticache.NewElastiCacheQueryHandler(ecService)

	// Initialize ELBv2
	elbService, err := elbv2.NewElbV2Service(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize ELBv2: %v", err)
	}
	s.handlers["elbv2"] = elbv2.NewElbV2QueryHandler(elbService)

	// Initialize Firehose
	fhService, err := firehose.NewFirehoseService(storageFactory, s3Service)
	if err != nil {
		log.Fatalf("Failed to initialize Firehose: %v", err)
	}
	s.handlers["firehose"] = firehose.NewFirehoseJsonHandler(fhService)

	// Initialize Glue
	glueService, err := glue.NewGlueService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Glue: %v", err)
	}
	s.handlers["glue"] = glue.NewGlueJsonHandler(glueService)

	// Initialize SQS
	sqsService, err := sqs.NewSqsService(storageFactory, "http://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to initialize SQS: %v", err)
	}
	s.handlers["sqs"] = sqs.NewSqsQueryHandler(sqsService)

	// Initialize SNS
	snsService, err := sns.NewSnsService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize SNS: %v", err)
	}
	s.handlers["sns"] = sns.NewSnsQueryHandler(snsService)

	// Initialize EventBridge
	ebService, err := eventbridge.NewEventBridgeService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize EventBridge: %v", err)
	}
	s.handlers["events"] = eventbridge.NewEventBridgeHandler(ebService)

	// Initialize Kinesis
	kinesisService, err := kinesis.NewKinesisService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Kinesis: %v", err)
	}
	s.handlers["kinesis"] = kinesis.NewKinesisJsonHandler(kinesisService)

	// Initialize KMS
	kmsService, err := kms.NewKmsService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize KMS: %v", err)
	}
	s.handlers["kms"] = kms.NewKmsJsonHandler(kmsService)

	// Initialize MSK
	mskService, err := msk.NewMskService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize MSK: %v", err)
	}
	s.handlers["msk"] = msk.NewMskJsonHandler(mskService)

	// Initialize Neptune
	neptuneService, err := neptune.NewNeptuneService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Neptune: %v", err)
	}
	s.handlers["neptune"] = neptune.NewNeptuneQueryHandler(neptuneService)

	// Initialize OpenSearch
	osService, err := opensearch.NewOpenSearchService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize OpenSearch: %v", err)
	}
	s.opensearch = opensearch.NewOpenSearchJsonHandler(osService)
	s.handlers["opensearch"] = s.opensearch

	// Initialize Pipes
	pipesService, err := pipes.NewPipesService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Pipes: %v", err)
	}
	s.handlers["pipes"] = pipes.NewPipesJsonHandler(pipesService)

	// Initialize Pricing
	s.handlers["pricing"] = pricing.NewPricingJsonHandler(pricing.NewPricingService())

	// Initialize RDS
	rdsService, err := rds.NewRdsService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize RDS: %v", err)
	}
	s.handlers["rds"] = rds.NewRdsQueryHandler(rdsService)

	// Initialize ResourceGroupsTagging
	rgtService, err := resourcegroupstagging.NewResourceGroupsTaggingService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize ResourceGroupsTagging: %v", err)
	}
	s.handlers["resourcegroupstagging"] = resourcegroupstagging.NewResourceGroupsTaggingJsonHandler(rgtService)

	// Initialize Route53
	r53Service, err := route53.NewRoute53Service(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Route53: %v", err)
	}
	s.route53 = route53.NewRoute53Handler(r53Service)

	// Initialize Scheduler
	schService, err := scheduler.NewSchedulerService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Scheduler: %v", err)
	}
	s.handlers["scheduler"] = scheduler.NewSchedulerJsonHandler(schService)

	// Initialize SecretsManager
	smService, err := secretsmanager.NewSecretsManagerService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize SecretsManager: %v", err)
	}
	s.handlers["secretsmanager"] = secretsmanager.NewSecretsManagerJsonHandler(smService)

	// Initialize SES
	sesService, err := ses.NewSesService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize SES: %v", err)
	}
	s.handlers["email"] = ses.NewSesQueryHandler(sesService)

	// Initialize SSM
	ssmService, err := ssm.NewSsmService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize SSM: %v", err)
	}
	s.handlers["ssm"] = ssm.NewSsmJsonHandler(ssmService)

	// Initialize StepFunctions
	sfnService, err := stepfunctions.NewStepFunctionsService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize StepFunctions: %v", err)
	}
	s.handlers["stepfunctions"] = stepfunctions.NewStepFunctionsJsonHandler(sfnService)

	// Initialize Textract
	s.handlers["textract"] = textract.NewTextractJsonHandler(textract.NewTextractService())

	// Initialize Transcribe
	trService, err := transcribe.NewTranscribeService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Transcribe: %v", err)
	}
	s.handlers["transcribe"] = transcribe.NewTranscribeJsonHandler(trService)

	// Initialize Transfer
	transferService, err := transfer.NewTransferService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize Transfer: %v", err)
	}
	s.handlers["transfer"] = transfer.NewTransferJsonHandler(transferService)

	// Initialize ACM
	acmService, err := acm.NewAcmService(storageFactory)
	if err != nil {
		log.Fatalf("Failed to initialize ACM: %v", err)
	}
	s.handlers["acm"] = acm.NewAcmJsonHandler(acmService)

	// Initialize new services
	redshiftService, _ := redshift.NewRedshiftService(storageFactory)
	s.handlers["redshift"] = redshift.NewRedshiftQueryHandler(redshiftService)

	emrService, _ := emr.NewEmrService(storageFactory)
	s.handlers["emr"] = emr.NewEmrJsonHandler(emrService)

	mqService, _ := mq.NewMqService(storageFactory)
	s.mq = mq.NewMqHandler(mqService)
	s.handlers["mq"] = s.mq

	docdbService, _ := docdb.NewDocdbService(storageFactory)
	s.handlers["docdb"] = docdb.NewDocdbQueryHandler(docdbService)

	cloudtrailService, _ := cloudtrail.NewCloudtrailService(storageFactory)
	s.handlers["cloudtrail"] = cloudtrail.NewCloudtrailJsonHandler(cloudtrailService)

	codecommitService, _ := codecommit.NewCodecommitService(storageFactory)
	s.handlers["codecommit"] = codecommit.NewCodecommitJsonHandler(codecommitService)

	codepipelineService, _ := codepipeline.NewCodepipelineService(storageFactory)
	s.handlers["codepipeline"] = codepipeline.NewCodepipelineJsonHandler(codepipelineService)

	xrayService, _ := xray.NewXrayService(storageFactory)
	s.handlers["xray"] = xray.NewXrayJsonHandler(xrayService)

	wafService, _ := waf.NewWafService(storageFactory)
	s.handlers["waf"] = waf.NewWafJsonHandler(wafService)

	organizationsService, _ := organizations.NewOrganizationsService(storageFactory)
	s.handlers["organizations"] = organizations.NewOrganizationsJsonHandler(organizationsService)

	syntheticsService, _ := synthetics.NewSyntheticsService(storageFactory)
	s.handlers["synthetics"] = synthetics.NewSyntheticsJsonHandler(syntheticsService)

	sagemakerService, _ := sagemaker.NewSagemakerService(storageFactory)
	s.handlers["sagemaker"] = sagemaker.NewSagemakerJsonHandler(sagemakerService)

	apprunnerService, _ := apprunner.NewApprunnerService(storageFactory)
	s.handlers["apprunner"] = apprunner.NewApprunnerJsonHandler(apprunnerService)

	// Initialize GCP Services
	cfSvc, _ := gcp_cf.NewCloudFunctionsService(storageFactory)
	s.gcpHandlers["cloudfunctions"] = gcp_cf.NewCloudFunctionsHandler(cfSvc)
	crSvc, _ := gcp_cr.NewCloudRunService(storageFactory)
	s.gcpHandlers["cloudrun"] = gcp_cr.NewCloudRunHandler(crSvc)
	dsSvc, _ := gcp_ds.NewDatastoreService(storageFactory)
	s.gcpHandlers["datastore"] = gcp_ds.NewDatastoreHandler(dsSvc)
	fsSvc, _ := gcp_fs.NewFirestoreService(storageFactory)
	s.gcpHandlers["firestore"] = gcp_fs.NewFirestoreHandler(fsSvc)
	gcsSvc, _ := gcp_gcs.NewGcsService(storageFactory)
	s.gcpHandlers["gcs"] = gcp_gcs.NewGcsHandler(gcsSvc)
	iamGcpSvc, _ := gcp_iam.NewIamGcpService(storageFactory)
	s.gcpHandlers["iam"] = gcp_iam.NewIamGcpHandler(iamGcpSvc)
	kafkaSvc, _ := gcp_kafka.NewKafkaService(storageFactory)
	s.gcpHandlers["kafka"] = gcp_kafka.NewKafkaHandler(kafkaSvc)
	opsSvc, _ := gcp_ops.NewOperationsService(storageFactory)
	s.gcpHandlers["operations"] = gcp_ops.NewOperationsHandler(opsSvc)
	pubsubSvc, _ := gcp_pubsub.NewPubSubService(storageFactory)
	s.gcpHandlers["pubsub"] = gcp_pubsub.NewPubSubHandler(pubsubSvc)
	smGcpSvc, _ := gcp_sm.NewSecretManagerGcpService(storageFactory)
	s.gcpHandlers["secretmanager"] = gcp_sm.NewSecretManagerGcpHandler(smGcpSvc)
	tasksSvc, _ := gcp_tasks.NewTasksService(storageFactory)
	s.gcpHandlers["tasks"] = gcp_tasks.NewTasksHandler(tasksSvc)
	bqSvc, _ := gcp_bq.NewBigQueryService(storageFactory, duckManager)
	s.gcpHandlers["bigquery"] = gcp_bq.NewBigQueryHandler(bqSvc)
	sqlSvc, _ := gcp_sql.NewCloudSqlService(storageFactory)
	s.gcpHandlers["cloudsql"] = gcp_sql.NewCloudSqlHandler(sqlSvc)
	spannerSvc, _ := gcp_spanner.NewSpannerService(storageFactory)
	s.gcpHandlers["spanner"] = gcp_spanner.NewSpannerHandler(spannerSvc)
	btSvc, _ := gcp_bt.NewBigtableService(storageFactory)
	s.gcpHandlers["bigtable"] = gcp_bt.NewBigtableHandler(btSvc)
	arSvc, _ := gcp_ar.NewArtifactRegistryService(storageFactory)
	s.gcpHandlers["artifactregistry"] = gcp_ar.NewArtifactRegistryHandler(arSvc)
	ceSvc, _ := gcp_ce.NewComputeEngineService(storageFactory)
	s.gcpHandlers["compute"] = gcp_ce.NewComputeEngineHandler(ceSvc)
	gkeSvc, _ := gcp_gke.NewGkeService(storageFactory)
	s.gcpHandlers["gke"] = gcp_gke.NewGkeHandler(gkeSvc)
	cbSvc, _ := gcp_cb.NewCloudBuildService(storageFactory)
	s.gcpHandlers["cloudbuild"] = gcp_cb.NewCloudBuildHandler(cbSvc)
	aeSvc, _ := gcp_ae.NewAppEngineService(storageFactory)
	s.gcpHandlers["appengine"] = gcp_ae.NewAppEngineHandler(aeSvc)
	wfSvc, _ := gcp_wf.NewWorkflowsService(storageFactory)
	s.gcpHandlers["workflows"] = gcp_wf.NewWorkflowsHandler(wfSvc)
	csSvc, _ := gcp_cs.NewCloudSchedulerService(storageFactory)
	s.gcpHandlers["cloudscheduler"] = gcp_cs.NewCloudSchedulerHandler(csSvc)
	dnsSvc, _ := gcp_dns.NewCloudDnsService(storageFactory)
	s.gcpHandlers["dns"] = gcp_dns.NewCloudDnsHandler(dnsSvc)
	armorSvc, _ := gcp_armor.NewCloudArmorService(storageFactory)
	s.gcpHandlers["armor"] = gcp_armor.NewCloudArmorHandler(armorSvc)
	lbSvc, _ := gcp_lb.NewLoadBalancingService(storageFactory)
	s.gcpHandlers["loadbalancing"] = gcp_lb.NewLoadBalancingHandler(lbSvc)
	casSvc, _ := gcp_cas.NewCasService(storageFactory)
	s.gcpHandlers["cas"] = gcp_cas.NewCasHandler(casSvc)
	loggingSvc, _ := gcp_logging.NewCloudLoggingService(storageFactory)
	s.gcpHandlers["logging"] = gcp_logging.NewCloudLoggingHandler(loggingSvc)
	monitoringSvc, _ := gcp_monitoring.NewCloudMonitoringService(storageFactory)
	s.gcpHandlers["monitoring"] = gcp_monitoring.NewCloudMonitoringHandler(monitoringSvc)
	traceSvc, _ := gcp_trace.NewCloudTraceService(storageFactory)
	s.gcpHandlers["trace"] = gcp_trace.NewCloudTraceHandler(traceSvc)

	// Register basic routes
	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	// Health check / ping
	s.mux.HandleFunc("/_cloudstack/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok", "version": "go-migration"}`)
	})
	
	// Catch-all for AWS API requests
	s.mux.Handle("/", RequestIDMiddleware(GlobalCorsFilterMiddleware(s.RequestContextMiddleware(http.HandlerFunc(s.handleAWSRequest)))))
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
