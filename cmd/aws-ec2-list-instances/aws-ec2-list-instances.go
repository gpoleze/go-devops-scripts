package main

import (
	"flag"
	"github.com/gpoleze/go-devops-scripts/pkg/aws/ec2"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"reflect"
)

func readFlags() (*string, *string) {
	var region string
	var profile string
	flag.StringVar(&region, "region", "", "AWS region")
	flag.StringVar(&region, "r", "", "AWS region (shorthand)")

	flag.StringVar(&profile, "profile", "", "AWS profile")
	flag.StringVar(&profile, "p", "", "AWS profile (shorthand)")

	flag.Parse()
	return &region, &profile
}

func main() {

	instances := ec2.DescribeInstances(readFlags())

	var header table.Row

	for _, field := range reflect.VisibleFields(reflect.TypeOf(instances[0])) {
		header = append(header, field.Name)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(header)

	for _, instance := range instances {
		t.AppendRow(table.Row{
			instance.Name,
			instance.Id,
			instance.Type,
			instance.State,
			instance.Ami,
			instance.LaunchTime,
			instance.PrivateIp,
			instance.PublicIp,
		})
	}

	t.Render()
}
