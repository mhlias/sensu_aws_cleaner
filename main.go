package main


import (
  "flag"
  "fmt"

  "github.com/mhlias/sensu_aws_cleaner/resources/aws"
  "github.com/mhlias/sensu_aws_cleaner/resources/sensu"
  "github.com/mhlias/sensu_aws_cleaner/resources/chef"


)



func main() {


  hostPtr := flag.String("host", "localhost", "Sensu API host address")
  portPtr := flag.Int("port", 4567, "Sensu API port")
  regionPtr := flag.String("region", "", "AWS Region to look for the instances in")
  removeChefPtr := flag.Bool("remove-chef", false, "Set to remove instance from managed chef too.")
  
  

  
  flag.Parse()

  if len(*regionPtr) <= 0 {
    fmt.Println("Error AWS Region is required.")
    return
  }

  event := flag.Arg(0)

  config := &aws.Config{ Region: *regionPtr }

  awsclient := config.Connect()

  ec2 := new(aws.EC2)

  ec2.GetAllInstances(awsclient)

  sensu := &sensu.Data{ Host: *hostPtr, Port: *portPtr }

  sensu.GetAllEvents([]byte(event))

  for _, event := range sensu.AllEvents {
    if event.Check.Name == "keepalive" && event.Check.Status >0 {
      if (!ec2.CheckInstanceState(awsclient, event.Client.Address)) {
        if sensu.CheckRemoveClient(event.Client.Name) {
          aws.RemoveRecords(awsclient, event.Client.Address)
          if *removeChefPtr {
            chef := &chef.Chef{}
            chef_node := chef.Find_instance(event.Client.Address)
            if len(chef_node) > 0 {
              chef.Remove_instance(chef_node)
            }
          }
        }
      }
    }
  } 

  










}