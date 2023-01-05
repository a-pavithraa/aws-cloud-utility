package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3BucketResponse struct {
	Name         string
	CreationDate time.Time
	Region       string
	Policy       string
}
type S3ObjectResponse struct {
	Name         string
	Owner        string
	LastModified time.Time
	Size         int64
}

func (bucket S3BucketResponse) GetRegion() string {
	return bucket.Region
}

type ListAndUploadApiClient interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(options *s3.Options)) (*s3.ListBucketsOutput, error)
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(options *s3.Options)) (*s3.ListObjectsOutput, error)
	DeleteObjects(ctx context.Context, params *s3.DeleteObjectsInput, optFns ...func(options *s3.Options)) (*s3.DeleteObjectsOutput, error)
	GetBucketPolicy(ctx context.Context, params *s3.GetBucketPolicyInput, optFns ...func(options *s3.Options)) (*s3.GetBucketPolicyOutput, error)
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(options *s3.Options)) (*s3.GetBucketLocationOutput, error)
}

func NewS3Client(ctx context.Context, region string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %s", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region

	})

	return s3Client, nil
}
func ListObjects(ctx context.Context, bucketName string, s3Client ListAndUploadApiClient) ([]S3ObjectResponse, error) {

	listObjectsOutput, err := s3Client.ListObjects(ctx, &s3.ListObjectsInput{Bucket: &bucketName})
	var (
		s3Objects []S3ObjectResponse
	)

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %s", err)
	}
	objects := listObjectsOutput.Contents
	for _, obj := range objects {
		s3Object := S3ObjectResponse{
			Name:         *obj.Key,
			LastModified: *obj.LastModified,
			Size:         obj.Size,
		}
		s3Objects = append(s3Objects, s3Object)
	}
	return s3Objects, nil
}

func ListAndGroupBucketsByRegion(ctx context.Context, client ListAndUploadApiClient) (map[string][]RestResponse, error) {
	fmt.Println("Inside ListAndGroupBucketsByRegion")
	allBuckets, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	bucketRegionMappings := map[string][]RestResponse{}

	if err != nil {
		return nil, fmt.Errorf("error in listing buckets,  %s", err)
	}
	if allBuckets != nil {

		for _, bucket := range allBuckets.Buckets {

			bucketDetail := S3BucketResponse{
				Name:         *bucket.Name,
				CreationDate: *bucket.CreationDate,
			}
			bucketLocation, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
				Bucket: bucket.Name,
			})
			if err != nil {
				return nil, fmt.Errorf("error in getting location,  %s", err)

			}

			if len(bucketLocation.LocationConstraint) > 0 {
				bucketDetail.Region = string(bucketLocation.LocationConstraint)
			} else {
				bucketDetail.Region = DefaultRegion
			}

			bucketList, _ := bucketRegionMappings[bucketDetail.Region]
			bucketList = append(bucketList, bucketDetail)
			bucketRegionMappings[bucketDetail.Region] = bucketList

		}
	}
	log.Println("bucketRegionMappings==", bucketRegionMappings)
	return bucketRegionMappings, nil

}

func ListAllBuckets() (map[string][]RestResponse, error) {
	ctx := context.Background()
	client, err := NewS3Client(ctx, DefaultRegion)
	if err != nil {
		return nil, err
	}
	return ListAndGroupBucketsByRegion(ctx, client)
}
