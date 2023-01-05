package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
)

type ResourceType int64
type EipStatus int64

const DefaultRegion = "us-east-1"
const (
	EC2 ResourceType = iota
	S3
	LB
	NGW
	EIP
	DYNAMODB
)
const (
	BOUND EipStatus = iota
	UNBOUND
)

type RestResponse interface {
	GetRegion() string
}

/*type AWSClientTypes interface {
	EC2DescribeInstancesAPI
}
type AWSClients[C AWSClientTypes] struct {
	Client C
}*/

func GetAwsResources(regions []string, fn func(ctx context.Context, region string) ([]RestResponse, error)) (map[string][]RestResponse, error) {

	regionWiseInstanceDetails := map[string][]RestResponse{}

	ctx := context.Background()
	resourcesChan := make(chan []RestResponse, len(regions))

	noOfRegions := int32(len(regions))
	wg, ctx := errgroup.WithContext(ctx)

	for _, region := range regions {
		//for closure
		reg := region

		wg.Go(func() error {
			defer func() {

				if atomic.AddInt32(&noOfRegions, -1) == 0 {
					close(resourcesChan)
				}
			}()
			resourceDetails, err := fn(ctx, reg)

			if err != nil {
				return err
			}
			resourcesChan <- resourceDetails
			return nil

		})

	}
	for resource := range resourcesChan {

		if len(resource) > 0 {

			reg := resource[0].GetRegion()
			regionWiseInstanceDetails[reg] = resource

		}

	}

	fmt.Println("regionWiseInstanceDetails====", regionWiseInstanceDetails)
	return regionWiseInstanceDetails, wg.Wait()

}
func NewEc2Client(ctx context.Context, region string) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %s", err)

	}
	client := ec2.NewFromConfig(cfg)
	return client, nil

}
func GetResourcesForAllRegions(resourceType ResourceType, regions []string) (map[string][]RestResponse, error) {

	var (
		response map[string][]RestResponse
		err      error
	)
	switch resourceType {
	case EC2:
		response, err = GetAwsResources(regions, GetInstanceForRegion)
	case EIP:
		response, err = GetAwsResources(regions, GetEIP)
	case LB:
		response, err = GetAwsResources(regions, GetLoadBalancerForRegion)

	case DYNAMODB:
		response, err = GetAwsResources(regions, ListDynamoDBTablesForRegion)

	}
	return response, err
}
