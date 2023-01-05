package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"testing"
)

type MockEc2Client struct {
	Ec2Input  *ec2.DescribeInstancesInput
	Ec2Output *ec2.DescribeInstancesOutput
}

func (client MockEc2Client) DescribeInstances(ctx context.Context,
	params *ec2.DescribeInstancesInput,
	optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return client.Ec2Output, nil
}
func TestDescribeEc2(t *testing.T) {
	ctx := context.Background()
	_, err := GetInstances(ctx, "us-east-1", MockEc2Client{
		Ec2Input: &ec2.DescribeInstancesInput{},
		Ec2Output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{{
						InstanceId: aws.String("test"),
					}},
				},
			}},
	})
	if err != nil {
		t.Fatalf("Describe Ec2 error error: %s", err)
	}
}
