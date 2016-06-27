## Overview

A Sensu handler in Go to clean up sensu entries when an AWS EC2 instance dies and also provides cleaning of Route53 records associated with the instance.



#### AwsCleaner

AwsCleaner is a handler that processes keepalive requests that have gone critical and checks the AWS account in the default region `eu-west-1` or the specified one by command line argument `-region` to find if the instance matching the private IP `does not exist` or is in `terminated` or `shutting-down` state and removes the instance from sensu using the local sensu-api. It also searches all hosted zones in Route53 to match the private IP and removes the record from them. In case of an A record pointing to multiple IPs it will be updated to not include the IP of the instance that doesn't exist anymore.

Requirements:
Sensu instance requires the following IAM profile access:

```
"Action": [
    "ec2:Describe*",
    "route53:ChangeResourceRecordSets",
    "route53:ListHostedZones",
    "route53:ListResourceRecordSets"
],

```