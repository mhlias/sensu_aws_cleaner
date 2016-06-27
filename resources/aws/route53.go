package aws


import (
    "log"
    
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/route53"
    
)



func RemoveRecords( meta interface{}, value string ) {

  resp, err := meta.(*AWSClient).r53conn.ListHostedZones(nil)


  if err != nil {
    panic(err)
  }

  for _, zone := range resp.HostedZones {

    params := &route53.ListResourceRecordSetsInput{
      HostedZoneId: zone.Id,
    }
    resp2, err2 := meta.(*AWSClient).r53conn.ListResourceRecordSets(params)

    if err2 != nil {
      panic(err)
    }


    for _, r := range resp2.ResourceRecordSets {

      if len(r.ResourceRecords) > 1 {
        for idx,v := range r.ResourceRecords {
          if *v.Value == value {

            r.ResourceRecords = append(r.ResourceRecords[:idx], r.ResourceRecords[idx+1:]...)

            params2 := &route53.ChangeResourceRecordSetsInput{
              ChangeBatch: &route53.ChangeBatch{
                Changes: []*route53.Change{
                  {
                    Action: aws.String("UPSERT"),
                    ResourceRecordSet: r,
                  },
                },
              },
              HostedZoneId: zone.Id,
            }

            _, err3 := meta.(*AWSClient).r53conn.ChangeResourceRecordSets(params2)

            if err3 == nil {
              log.Println("[INFO] RecordSet Deleted!")
            } else {
              panic(err3)
            }

            break
          }
        }

      } else if len(r.ResourceRecords) == 1 {
        if *r.ResourceRecords[0].Value == value  {
          
          params2 := &route53.ChangeResourceRecordSetsInput{
            ChangeBatch: &route53.ChangeBatch{
              Changes: []*route53.Change{
                {
                  Action: aws.String("DELETE"),
                  ResourceRecordSet: r,
                },
              },
            },
            HostedZoneId: zone.Id,
          }

          _, err3 := meta.(*AWSClient).r53conn.ChangeResourceRecordSets(params2)

          if err3 == nil {
            log.Println("[INFO] RecordSet Deleted!")
          } else {
            panic(err3)
          }
        }
      }

    }


  }

}




