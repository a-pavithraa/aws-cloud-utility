package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2InstanceDetails struct {
	Id        string `json:"instanceId"`
	State     string `json:"instanceState"`
	Type      string `json:"instanceType"`
	PublicIp  string `json:"publicIp"`
	PrivateIp string `json:"privateIp"`
	VpcId     string `json:"vpcId"`
	Region    string `json:"region"`
}

func (e EC2InstanceDetails) GetRegion() string {
	return e.Region

}

// EC2DescribeInstancesAPI defines the interface for the DescribeInstances function.
// We use this interface to test the function using a mocked service.
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func GetInstances(ctx context.Context, region string, client EC2DescribeInstancesAPI) ([]RestResponse, error) {
	var ec2InstanceDetailsList []RestResponse

	input := &ec2.DescribeInstancesInput{}

	result, err := client.DescribeInstances(ctx, input)

	if err != nil {
		return nil, err
	}

	for _, r := range result.Reservations {
		fmt.Println("Reservation==", r)

		for _, i := range r.Instances {
			//Checking private Ip since terminated instances take some time to be cleaned up
			if i.PrivateIpAddress != nil {
				instanceDetails := EC2InstanceDetails{
					Id:        *i.InstanceId,
					State:     string(i.State.Name),
					Type:      string(i.InstanceType),
					PrivateIp: *i.PrivateIpAddress,
					VpcId:     *i.VpcId,
					Region:    region,
				}
				if i.PublicIpAddress != nil {
					instanceDetails.PublicIp = *i.PublicIpAddress
				}

				ec2InstanceDetailsList = append(ec2InstanceDetailsList, instanceDetails)
				fmt.Println(instanceDetails)
			}

		}

	}
	return ec2InstanceDetailsList, nil

}
func GetInstanceForRegion(ctx context.Context, region string) ([]RestResponse, error) {
	client, err := NewEc2Client(ctx, region)
	if err != nil {
		return nil, err
	}
	return GetInstances(ctx, region, client)

}

func StopInstance(region string, instanceId string) error {
	ctx := context.Background()
	client, err := NewEc2Client(ctx, region)
	if err != nil {
		return err

	}

	stopInstanceInput := ec2.StopInstancesInput{InstanceIds: []string{instanceId}}
	client.StopInstances(ctx, &stopInstanceInput)
	return nil
}

func StartInstance(region string, instanceId string) error {
	ctx := context.Background()
	client, err := NewEc2Client(ctx, region)
	if err != nil {
		return err

	}

	startInstanceInput := ec2.StartInstancesInput{InstanceIds: []string{instanceId}}
	client.StartInstances(ctx, &startInstanceInput)
	return nil

}
