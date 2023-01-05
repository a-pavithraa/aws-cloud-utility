package service

//https://santoshk.dev/posts/2022/release-dangling-elastic-ips-using-lambda-and-go-sdk/
import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"golang.org/x/net/context"
)

type EIPDetails struct {
	AllocationId  string `json:"allocationId"`
	AssociationId string `json:"associationId"`
	InstanceId    string `json:"instanceId"`
	PublicIp      string `json:"publicIp"`
	Region        string `json:"region"`
}

func (e EIPDetails) GetRegion() string {
	return e.Region

}

type EIPDescribeAndReleaseAPI interface {
	ReleaseAddress(ctx context.Context, params *ec2.ReleaseAddressInput, optFns ...func(*ec2.Options)) (*ec2.ReleaseAddressOutput, error)
	DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(options *ec2.Options)) (*ec2.DescribeAddressesOutput, error)
}

func ReleaseEIP(region string, allocationId string) error {

	ctx := context.Background()
	ec2ServiceClient, err := NewEc2Client(ctx, region)
	if err != nil {
		return err
	}

	ReleaseAddressFilter := &ec2.ReleaseAddressInput{AllocationId: &allocationId}
	_, err = ec2ServiceClient.ReleaseAddress(ctx, ReleaseAddressFilter)

	if err != nil {
		return err
	}
	return nil

}
func GetEIPDetails(ctx context.Context, region string, client EIPDescribeAndReleaseAPI) ([]RestResponse, error) {
	var (
		eipDetails []RestResponse
	)

	IpListFilter := &ec2.DescribeAddressesInput{}

	result, err := client.DescribeAddresses(ctx, IpListFilter)
	fmt.Println(result)
	if err != nil {

		return nil, err
	}

	for _, address := range result.Addresses {
		fmt.Println("Allocation ID: " + *address.AllocationId)
		fmt.Println("Public IP: " + *address.PublicIp)

		eipDetail := EIPDetails{
			AllocationId: *address.AllocationId,
			PublicIp:     *address.PublicIp,
			Region:       region,
		}
		if address.AssociationId != nil {
			eipDetail.AssociationId = *address.AssociationId
		}
		eipDetails = append(eipDetails, eipDetail)

	}
	return eipDetails, nil

}

func GetEIP(ctx context.Context, region string) ([]RestResponse, error) {
	client, err := NewEc2Client(ctx, region)
	if err != nil {
		return nil, err
	}
	return GetEIPDetails(ctx, region, client)

}
