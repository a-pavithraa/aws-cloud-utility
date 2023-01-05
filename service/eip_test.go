package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"testing"
)

type MockEIPClient struct {
	eipReleaseInput          *ec2.ReleaseAddressInput
	eipReleaseOutput         *ec2.ReleaseAddressOutput
	eipDescribeAddressInput  *ec2.DescribeAddressesInput
	eipDescribeAddressOutput *ec2.DescribeAddressesOutput
}

func (client MockEIPClient) ReleaseAddress(ctx context.Context, params *ec2.ReleaseAddressInput, optFns ...func(*ec2.Options)) (*ec2.ReleaseAddressOutput, error) {
	return client.eipReleaseOutput, nil
}

func (client MockEIPClient) DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(options *ec2.Options)) (*ec2.DescribeAddressesOutput, error) {
	return client.eipDescribeAddressOutput, nil
}
func TestDescribeEip(t *testing.T) {
	ctx := context.Background()
	_, err := GetEIPDetails(ctx, "us-east-1", MockEIPClient{
		eipDescribeAddressInput: &ec2.DescribeAddressesInput{},
		eipDescribeAddressOutput: &ec2.DescribeAddressesOutput{
			Addresses: []types.Address{},
		},
	})
	if err != nil {
		t.Fatalf("Describe Eip error : %s", err)
	}
}
