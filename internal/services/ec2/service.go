package ec2

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/ec2/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

const (
	hexChars = "0123456789abcdef"
)

type Ec2Service struct {
	vpcStore      storage.Backend[string, *model.Vpc]
	instanceStore storage.Backend[string, *model.Instance]
	mu            sync.RWMutex
}

func NewEc2Service(factory *storage.Factory) (*Ec2Service, error) {
	vpcStore, _ := storage.CreateAccountAware[*model.Vpc](factory, "ec2", "ec2-vpcs.json", "wal")
	instanceStore, _ := storage.CreateAccountAware[*model.Instance](factory, "ec2", "ec2-instances.json", "wal")

	return &Ec2Service{
		vpcStore:      vpcStore,
		instanceStore: instanceStore,
	}, nil
}

func (s *Ec2Service) RunInstances(ctx context.Context, imageID, instanceType string, count int) (*model.Reservation, error) {
	reservationID := "r-" + s.randomHex(17)
	res := &model.Reservation{
		ReservationID: reservationID,
		OwnerID:       "000000000000",
		Instances:     make([]*model.Instance, 0),
	}

	for i := 0; i < count; i++ {
		instID := "i-" + s.randomHex(17)
		inst := &model.Instance{
			InstanceID:   instID,
			ImageID:      imageID,
			InstanceType: instanceType,
			State:        "running",
			LaunchTime:   time.Now(),
			PrivateIp:    fmt.Sprintf("172.31.0.%d", 10+i),
		}
		s.instanceStore.Put(ctx, instID, inst)
		res.Instances = append(res.Instances, inst)
	}

	log.Printf("RunInstances: created %d instances in reservation %s", count, reservationID)
	return res, nil
}

func (s *Ec2Service) DescribeInstances(ctx context.Context) ([]*model.Reservation, error) {
	instances, err := s.instanceStore.Scan(ctx, func(k string) bool { return true })
	if err != nil {
		return nil, err
	}
	
	// For simplicity, group all instances into one reservation
	res := &model.Reservation{
		ReservationID: "r-default",
		OwnerID:       "000000000000",
		Instances:     instances,
	}
	return []*model.Reservation{res}, nil
}

func (s *Ec2Service) randomHex(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(b)
}
