package aws


import (
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/route53"
    "github.com/aws/aws-sdk-go/aws/credentials"

)



type AWSClient struct {
  ec2conn            *ec2.EC2
  r53conn            *route53.Route53

  region             string
}


type Config struct {
  Awsconf string
  Profile string
  Region  string
}




func (c *Config) Connect() interface{} {


  var client AWSClient
  
  awsConfig := new(aws.Config)

  if len(c.Profile)>0 {
    awsConfig = &aws.Config{
      Credentials: credentials.NewSharedCredentials(c.Awsconf, fmt.Sprintf("profile %s", c.Profile)),
      Region:      aws.String(c.Region),
      MaxRetries:  aws.Int(3),
    }

  } else {
    // use instance role
    awsConfig = &aws.Config{
      Region:      aws.String(c.Region),
    }

  }


  sess := session.New(awsConfig)

  client.ec2conn = ec2.New(sess)

  client.r53conn = route53.New(sess)


  return &client

}



