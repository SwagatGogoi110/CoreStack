package lambda

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/aws/lambda/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type LambdaService struct {
	functionStore storage.Backend[string, *model.LambdaFunction]
	codeRoot      string
}

func NewLambdaService(factory *storage.Factory, codeRoot string) (*LambdaService, error) {
	functionStore, _ := storage.CreateAccountAware[*model.LambdaFunction](factory, "lambda", "lambda-functions.json", "wal")

	s := &LambdaService{
		functionStore: functionStore,
		codeRoot:      codeRoot,
	}

	if err := os.MkdirAll(codeRoot, 0755); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *LambdaService) CreateFunction(ctx context.Context, fn *model.LambdaFunction, zipCode []byte) (*model.LambdaFunction, error) {
	if _, ok, _ := s.functionStore.Get(ctx, fn.FunctionName); ok {
		return nil, fmt.Errorf("ResourceConflictException: Function already exists")
	}

	fn.FunctionArn = fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", fn.FunctionName)
	fn.LastModified = time.Now()
	fn.RevisionID = uuid.New().String()

	// Store code
	if len(zipCode) > 0 {
		codePath := filepath.Join(s.codeRoot, fn.FunctionName+".zip")
		if err := os.WriteFile(codePath, zipCode, 0644); err != nil {
			return nil, err
		}
		fn.CodeLocalPath = codePath
		fn.CodeSizeBytes = int64(len(zipCode))
	}

	if err := s.functionStore.Put(ctx, fn.FunctionName, fn); err != nil {
		return nil, err
	}

	log.Printf("Created Lambda function: %s", fn.FunctionName)
	return fn, nil
}

func (s *LambdaService) GetFunction(ctx context.Context, name string) (*model.LambdaFunction, error) {
	fn, ok, err := s.functionStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("ResourceNotFoundException: Function not found")
	}
	return fn, nil
}

func (s *LambdaService) ListFunctions(ctx context.Context) ([]*model.LambdaFunction, error) {
	return s.functionStore.Scan(ctx, func(k string) bool { return true })
}

func (s *LambdaService) Invoke(ctx context.Context, name string, payload []byte) (*model.InvokeResult, error) {
	// TODO: Implement Docker-based execution
	return nil, fmt.Errorf("Invoke not yet implemented in Go version")
}
