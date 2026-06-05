package cloudwatch

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/cloudwatch/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type CloudWatchService struct {
	metricStore storage.Backend[string, *model.MetricDatum]
	alarmStore  storage.Backend[string, *model.MetricAlarm]
	mu          sync.RWMutex
}

func NewCloudWatchService(factory *storage.Factory) (*CloudWatchService, error) {
	metricStore, _ := storage.CreateAccountAware[*model.MetricDatum](factory, "cloudwatch", "cw-metrics.json", "wal")
	alarmStore, _ := storage.CreateAccountAware[*model.MetricAlarm](factory, "cloudwatch", "cw-alarms.json", "wal")

	return &CloudWatchService{
		metricStore: metricStore,
		alarmStore:  alarmStore,
	}, nil
}

func (s *CloudWatchService) PutMetricData(ctx context.Context, namespace string, data []*model.MetricDatum) error {
	for _, d := range data {
		d.Namespace = namespace
		if d.Timestamp.IsZero() {
			d.Timestamp = time.Now()
		}
		
		key := fmt.Sprintf("%s:%s:%s", namespace, d.MetricName, uuid.New().String())
		if err := s.metricStore.Put(ctx, key, d); err != nil {
			return err
		}
	}
	log.Printf("PutMetricData: %d datums in %s", len(data), namespace)
	return nil
}

func (s *CloudWatchService) PutMetricAlarm(ctx context.Context, alarm *model.MetricAlarm) error {
	alarm.AlarmArn = fmt.Sprintf("arn:aws:cloudwatch:us-east-1:000000000000:alarm:%s", alarm.AlarmName)
	alarm.StateValue = "OK"
	
	if err := s.alarmStore.Put(ctx, alarm.AlarmName, alarm); err != nil {
		return err
	}
	log.Printf("PutMetricAlarm: %s", alarm.AlarmName)
	return nil
}

func (s *CloudWatchService) DescribeAlarms(ctx context.Context) ([]*model.MetricAlarm, error) {
	return s.alarmStore.Scan(ctx, func(k string) bool { return true })
}
