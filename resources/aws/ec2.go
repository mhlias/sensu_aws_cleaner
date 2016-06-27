package aws


import (

    "github.com/aws/aws-sdk-go/service/ec2"
    
)


type EC2 struct {

  Instances *ec2.DescribeInstancesOutput
  AllInstances map[string] string

  last_err error

}


func (e *EC2) GetAllInstances( meta interface{} ) {

  instances, err := meta.(*AWSClient).ec2conn.DescribeInstances(nil)
  if err != nil {
      panic(err)
  }

  e.AllInstances = make(map[string] string)

  for idx, _ := range instances.Reservations {
    for _, inst := range instances.Reservations[idx].Instances {
      
      e.AllInstances[*inst.InstanceId] = *inst.State.Name

      if inst.PrivateIpAddress!= nil {
        e.AllInstances[*inst.PrivateIpAddress] = *inst.State.Name
      }
      
    }
  }

  if len(e.AllInstances) <1 {
    panic("No EC2 instances found in AWS, aborting as unable to function.")
  }


}


func (e *EC2) CheckInstanceState( meta interface{}, key string ) bool {

  if state, ok := e.AllInstances[key]; ok {

    if ( state == "terminated" || state == "shutting-down" ) {
      return false
    } else {
      return true
    }

  }

  return false

}




