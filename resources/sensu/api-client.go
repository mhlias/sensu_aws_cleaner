package sensu

import (
  "fmt"
  "log"
  "encoding/json"

  "github.com/parnurzeal/gorequest"

)


type Data struct {

  Host string
  Port int

  Instances []Instance
  AllEvents []Event

}

type Instance struct {

  Name string
  Address string
  Subscriptions []string
  Timestamp int64

}

type CheckData struct {

  Name string
  Command string
  Subscribers []string
  Interval int
  Issued int64
  Executed int64
  Output string
  Status int
  Duration float32
  History []string

}


type Event struct {

  Id string
  Client Instance
  Check CheckData
  Occurences int
  Action string


}



var request = gorequest.New()


func (d *Data) GetAllClients() {

  
  resp, body, errs := request.Get(fmt.Sprintf("http://%s:%d/clients", d.Host, d.Port)).EndBytes()
  if errs != nil {
    panic(errs)
  }
  
  if resp.StatusCode == 200 {

    json_err := json.Unmarshal(body, &d.Instances)

    if json_err!=nil {
      panic(json_err)
    }

  } else {
    panic("Error response received by Sensu-API")
  }

}


func (d *Data) GetAllEvents(event []byte) {

  if len(event) > 0 {

    json_err := json.Unmarshal(event, &d.AllEvents)

    if json_err!=nil {
      panic(json_err)
    }

  } else {
  
    resp, body, errs := request.Get(fmt.Sprintf("http://%s:%d/events", d.Host, d.Port)).EndBytes()
    if errs != nil {
      panic(errs)
    }
    
    if resp.StatusCode == 200 {

      json_err := json.Unmarshal(body, &d.AllEvents)

      if json_err!=nil {
        panic(json_err)
      }

    } else {
      panic("Error response received by Sensu-API")
    }

  }

}



func (d *Data) CheckRemoveClient(name string) bool {

  resp, body, errs := request.Get(fmt.Sprintf("http://%s:%d/events/%s/keepalive", d.Host, d.Port, name)).EndBytes()

  if errs != nil {
    panic(errs)
  }

  if(resp.StatusCode == 200) {

    client := new(Event)

    json_err := json.Unmarshal(body, &client)

    if json_err!=nil {
      log.Printf("[INFO] No keepalive events for instance: %s", name)
    } else {

      if client.Check.Status>0 {
        log.Printf("[INFO] Instance %s not available any more. Removing from Sensu.\n", name)
        
        resp2, _, errs2 := request.Delete(fmt.Sprintf("http://%s:%d/clients/%s", d.Host, d.Port, name)).End()
        if errs2 != nil {
          panic(errs2)
        }
        if resp2.StatusCode == 202 {
          log.Printf("[INFO] Sucessfully removed instance %s from Sensu.\n", name)
          return true
        } else {
          log.Printf("[ERROR] Failed to remove instance %s from Sensu.\n", name)
          return false
        }

      }

    }

  }

  return false

} 




