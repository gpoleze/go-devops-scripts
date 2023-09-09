package ec2

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"time"
)

type MyInstanceInfo struct {
	Name       string
	Id         string
	Type       string
	State      string
	Ami        string
	LaunchTime time.Time
	PrivateIp  string
	PublicIp   string
}

// EC2DescribeInstancesAPI defines the interface for the DescribeInstances function.
// We use this interface to test the function using a mocked service.
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// GetInstances retrieves information about your Amazon Elastic Compute Cloud (Amazon EC2) instances.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a DescribeInstancesOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to DescribeInstances.
func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func filterTagByKey(tags []types.Tag, keyName string) string {
	value := ""
	for _, tag := range tags {
		if *tag.Key == keyName {
			value = *tag.Value
		}
	}
	return value
}

func DescribeInstances(region *string, profile *string) []MyInstanceInfo {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*region), config.WithSharedConfigProfile(*profile))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{}

	result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return nil
	}

	var myInstances = []MyInstanceInfo{}

	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			var myInstance = MyInstanceInfo{
				Name:       filterTagByKey(i.Tags, "Name"),
				Id:         *i.InstanceId,
				Type:       string(i.InstanceType),
				State:      string(i.State.Name),
				Ami:        *i.ImageId,
				LaunchTime: *i.LaunchTime,
			}
			if i.PrivateIpAddress != nil {
				myInstance.PrivateIp = *i.PrivateIpAddress
			}
			if i.PublicIpAddress != nil {
				myInstance.PublicIp = *i.PublicIpAddress
			}
			myInstances = append(myInstances, myInstance)
		}
	}

	return myInstances
}
