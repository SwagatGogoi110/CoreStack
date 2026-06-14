package autoscaling

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hectorvent/cloudstack/internal/services/aws/autoscaling/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type AutoScalingService struct {
	groupStore storage.Backend[string, *model.AutoScalingGroup]
	lcStore    storage.Backend[string, *model.LaunchConfiguration]
}

func NewAutoScalingService(factory *storage.Factory) (*AutoScalingService, error) {
	groupStore, _ := storage.CreateAccountAware[*model.AutoScalingGroup](factory, "autoscaling", "asg-groups.json", "wal")
	lcStore, _ := storage.CreateAccountAware[*model.LaunchConfiguration](factory, "autoscaling", "asg-lcs.json", "wal")

	return &AutoScalingService{
		groupStore: groupStore,
		lcStore:    lcStore,
	}, nil
}

func (s *AutoScalingService) CreateAutoScalingGroup(ctx context.Context, g *model.AutoScalingGroup) error {
	if _, ok, _ := s.groupStore.Get(ctx, g.AutoScalingGroupName); ok {
		return fmt.Errorf("AlreadyExists: Auto Scaling group already exists")
	}

	g.AutoScalingGroupArn = fmt.Sprintf("arn:aws:autoscaling:us-east-1:000000000000:autoScalingGroup:%s", g.AutoScalingGroupName)
	g.CreatedTime = time.Now()

	if err := s.groupStore.Put(ctx, g.AutoScalingGroupName, g); err != nil {
		return err
	}

	log.Printf("Created Auto Scaling group: %s", g.AutoScalingGroupName)
	return nil
}

func (s *AutoScalingService) DescribeAutoScalingGroups(ctx context.Context, names []string) ([]*model.AutoScalingGroup, error) {
	groups, err := s.groupStore.Scan(ctx, func(k string) bool {
		if len(names) == 0 {
			return true
		}
		for _, n := range names {
			if n == k {
				return true
			}
		}
		return false
	})
	return groups, err
}

func (s *AutoScalingService) CreateLaunchConfiguration(ctx context.Context, lc *model.LaunchConfiguration) error {
	if _, ok, _ := s.lcStore.Get(ctx, lc.LaunchConfigurationName); ok {
		return fmt.Errorf("AlreadyExists: Launch configuration already exists")
	}

	lc.LaunchConfigurationArn = fmt.Sprintf("arn:aws:autoscaling:us-east-1:000000000000:launchConfiguration:%s", lc.LaunchConfigurationName)
	lc.CreatedTime = time.Now()

	if err := s.lcStore.Put(ctx, lc.LaunchConfigurationName, lc); err != nil {
		return err
	}

	log.Printf("Created Launch Configuration: %s", lc.LaunchConfigurationName)
	return nil
}
