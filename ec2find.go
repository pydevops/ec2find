package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "EC2 Instance Finder"
	app.Usage = "Use this app to find the IP addresses of the ec2 instance by searching a given tag name"
	app.Version = "0.1.0"

	var profileName string = "default"
	var instanceStateName string = "running"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "profile,p",
			Value:       "default",
			Usage:       "aws `profile`",
			Destination: &profileName,
		},
		cli.StringFlag{
			Name:        "status,s",
			Value:       "running",
			Usage:       "`instance-state-name`",
			Destination: &instanceStateName,
		},
		cli.BoolFlag{
			Name:  "deploy_group,d",
			Usage: "search by deploy_group tag instead",
		},
		cli.BoolFlag{
			Name:  "login,l",
			Usage: "Prompt to log-in (ssh) after list is returned",
		},
	}
	app.Action = func(ctx *cli.Context) error {

		//https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
		// uses shared config fields
		sess, err := session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           profileName,
		})

		if err != nil {
			log.Fatalf("failed to create session %v\n", err)
		}

		// another way of
		// svc := ec2.New(sess, &aws.Config{
		// 	Region:      aws.String(awsRegion),
		// 	Credentials: credentials.NewSharedCredentials("", profileName),
		// })

		var name string = "tag:Name"
		if ctx.Bool("deploy_group") {
			name = "tag:DeployGroup"
		}

		nameFilter := ctx.Args().First()
		svc := ec2.New(sess)
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String(name),
					Values: []*string{
						aws.String(strings.Join([]string{"*", nameFilter, "*"}, "")),
					},
				},
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String(instanceStateName),
					},
				},
			},
		}

		// send API request for AWS ec2
		resp, err := svc.DescribeInstances(params)
		if err != nil {
			fmt.Println("error list instances: ", err.Error())
			log.Fatal(err.Error())
		}

		// output
		ip_list := make([]string, 0)
		for index, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {

				name := "None"
				for _, keys := range instance.Tags {
					if *keys.Key == "Name" {
						name = url.QueryEscape(*keys.Value)
						break
					}
				}
				instanceID := *instance.InstanceId
				instanceType := *instance.InstanceType
				privateIP := *instance.PrivateIpAddress
				ip_list = append(ip_list, privateIP)
				fmt.Printf("%d %s %s  %s  %s\n", index+1, name, instanceID, instanceType, privateIP)
			}
		}

		if ctx.Bool("login") {
			var selection string
			var index int
			for {
				fmt.Printf("Enter number:")
				_, err := fmt.Scanf("%s", &selection)
				index, err = strconv.Atoi(selection)
				if err != nil {
					fmt.Printf("%q is NOT a number \n", selection)
				}
				if index > 0 && index <= len(ip_list) {
					break
				}
			}
			ip := ip_list[index-1]
			sshCmd := exec.Command("ssh", ip)
			sshCmd.Stdin = os.Stdin
			sshCmd.Stdout = os.Stdout
			sshCmd.Stderr = os.Stderr
			err := sshCmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	}

	app.Run(os.Args)

}
