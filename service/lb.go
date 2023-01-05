package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type LoadBalancerResponse struct {
	Name    string `json:"name"`
	Arn     string `json:"arn"`
	State   string `json:"state"`
	Region  string `json:"region"`
	DnsName string `json:"dnsName"`
}

type LoadBalancerDescribeAndDeleteClient interface {
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(options *elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error)
	DeleteLoadBalancer(ctx context.Context, params *elasticloadbalancingv2.DeleteLoadBalancerInput, optFns ...func(options *elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DeleteLoadBalancerOutput, error)
}

func (lb LoadBalancerResponse) GetRegion() string {
	return lb.Region
}

func NewLBClient(ctx context.Context, region string) (*elasticloadbalancingv2.Client, error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %s", err)
	}
	loadBalancingClient := elasticloadbalancingv2.NewFromConfig(cfg)
	return loadBalancingClient, nil
}
func GetLoadBalancerDetails(ctx context.Context, region string, client LoadBalancerDescribeAndDeleteClient) ([]RestResponse, error) {
	var lbDetailsList []RestResponse
	resp, err := client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{})

	if err != nil {
		return nil, err
	}
	for _, lb := range resp.LoadBalancers {
		loadBalancerResponse := LoadBalancerResponse{
			Name:    *lb.LoadBalancerName,
			Arn:     *lb.LoadBalancerArn,
			DnsName: *lb.DNSName,
			Region:  region,
			State:   string(lb.State.Code),
		}
		lbDetailsList = append(lbDetailsList, loadBalancerResponse)
		fmt.Println("name===>", *lb.LoadBalancerName, "dns name==", *lb.DNSName, "", lb.State.Code, lb.LoadBalancerArn)

	}
	return lbDetailsList, nil
}

func GetLoadBalancerForRegion(ctx context.Context, region string) ([]RestResponse, error) {
	client, err := NewLBClient(ctx, region)
	if err != nil {
		return nil, err
	}
	return GetLoadBalancerDetails(ctx, region, client)

}

func DeleteLoadBalancer(region string, arn string) error {
	ctx := context.Background()
	lbClient, err := NewLBClient(ctx, region)
	if err != nil {
		return err
	}
	lbClient.DeleteLoadBalancer(ctx, &elasticloadbalancingv2.DeleteLoadBalancerInput{
		LoadBalancerArn: &arn,
	})
	return nil
}
