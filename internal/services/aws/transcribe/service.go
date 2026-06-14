package transcribe

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/transcribe/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type TranscribeService struct {
	jobStore storage.Backend[string, *model.TranscriptionJob]
}

func NewTranscribeService(factory *storage.Factory) (*TranscribeService, error) {
	jobStore, _ := storage.CreateAccountAware[*model.TranscriptionJob](factory, "transcribe", "tr-jobs.json", "wal")

	return &TranscribeService{
		jobStore: jobStore,
	}, nil
}

func (s *TranscribeService) StartTranscriptionJob(ctx context.Context, name, language, format string) (*model.TranscriptionJob, error) {
	if _, ok, _ := s.jobStore.Get(ctx, name); ok {
		return nil, fmt.Errorf("ConflictException: Job already exists")
	}

	job := &model.TranscriptionJob{
		TranscriptionJobName:   name,
		TranscriptionJobStatus: "COMPLETED",
		LanguageCode:           language,
		MediaFormat:            format,
		Transcript: &model.Transcript{
			TranscriptFileUri: fmt.Sprintf("s3://CloudStack-transcribe/%s.json", name),
		},
		CreationTime: time.Now().Unix(),
	}

	if err := s.jobStore.Put(ctx, name, job); err != nil {
		return nil, err
	}

	log.Printf("Started Transcribe job: %s", name)
	return job, nil
}

func (s *TranscribeService) GetTranscriptionJob(ctx context.Context, name string) (*model.TranscriptionJob, error) {
	job, ok, err := s.jobStore.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("NotFoundException: Job not found")
	}
	return job, nil
}
