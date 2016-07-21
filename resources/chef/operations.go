package chef

import (
	"log"
	"os"
	"fmt"
	"net/url"
	"strings"
	"bufio"
	"path/filepath"
	"io/ioutil"

	"github.com/mattn/go-shellwords"

	"gopkg.in/chef.v0"
)

type Chef struct {
	endpoint string
	user     string
	key      string
}



func (ch *Chef) Find_instance(ip string) string {

	node_name := ""

	if !ch.read_config() {
		log.Fatal("[ERROR] Failed to read Chef config. Aborting.\n")
	}
	
	
	client, err := chef.NewClient(&chef.Config{
    Name: ch.user,
    Key:  ch.key,
    BaseURL: ch.endpoint,
	})

	if err != nil {
		fmt.Println("Issue setting up Chef API client:", err)
	}

	query, err := client.Search.NewQuery("node", fmt.Sprintf("ipaddress:%s", ip))
	
	if err != nil {
		log.Fatal("Error building query ", err)
	}

	
	res, err := query.Do(client)
	if err != nil {
		log.Fatal("Error running query ", err)
	}

	for _,row := range res.Rows {
	  for field, data := range row.(map[string]interface{}) {
      if field == "name" {
      	node_name = data.(string)
      	break
      }
	  }
	}


	return node_name



}


func (ch *Chef) Remove_instance(name string) bool {

	if len(name) < 1 {
		log.Println("Chef node name was not provided.")
		return false
	}

	if len(ch.endpoint) <= 0 || len(ch.user) <= 0 || len(ch.key) <= 0 {

		if !ch.read_config() {
			log.Fatal("[ERROR] Failed to read Chef config. Aborting.\n")
			return false
		}

	}

	client, err := chef.NewClient(&chef.Config{
    Name: ch.user,
    Key:  ch.key,
    BaseURL: ch.endpoint,
	})

	if err != nil {
		log.Println("Issue setting up Chef API client:", err)
		return false
	}

	del_err := client.Nodes.Delete(name)

	if del_err != nil {
		log.Println("Failed to delete instance from Chef server: ", err)
		return false
	}

	return true

}

func (ch *Chef) read_config(config_file ...string) bool {

	// lifted from marpaia/chef-golang with additions/imporovements
	knifeFiles := []string{}

	if len(config_file) > 0 {
		for _, v := range config_file {
			knifeFiles = append(knifeFiles, v)
		}
	}

	knifeFiles = append(knifeFiles, ".chef/knife.rb")

	homedir := os.Getenv("HOME")
	if homedir != "" {
		knifeFiles = append(knifeFiles, filepath.Join(homedir, ".chef/knife.rb"))
	}

	knifeFiles = append(knifeFiles, "/etc/chef/client.rb")

	var knifeFile string
	for _, file := range knifeFiles {
		if _, err := os.Stat(file); err == nil {
			knifeFile = file
			break
		}
	}

	if knifeFile == "" {
		log.Println("No Chef configuration file could be found.")
		return false
	}

	file, err := os.Open(knifeFile)
	defer file.Close()
	if err != nil {
		log.Println("Failed to open chef config files.")
		return false
	}

	scanner := bufio.NewScanner(file)

	chefHost := ""
	chefPort := ""
	chefOrg  := ""
	chefPath := ""
	chefUrl  := &url.URL{}

	for scanner.Scan() {
		split, _ := shellwords.Parse(scanner.Text())
		if len(split) == 2 {
			switch split[0] {
			case "node_name":
				ch.user = parse(split[1])
			case "client_key":
				key, err := ioutil.ReadFile(parse(split[1]))
				if err != nil {
					log.Println("Could not load private client key.")
					return false
				}
				ch.key = string(key)
			case "chef_server_url":
				parsedUrl := parse(split[1])

				var url_err error
				
				chefUrl, url_err = url.Parse(parsedUrl)
				if url_err != nil {
					log.Println("Invalid Chef Host URL.")
					return false
				}
				hostPath := strings.Split(chefUrl.Path, "/")
				chefPath = hostPath[1]
				if len(hostPath) == 3 && hostPath[1] == "organizations" {
					chefOrg = hostPath[2]
				}
				hostPort := strings.Split(chefUrl.Host, ":")
				if len(hostPort) == 2 {
					chefHost = hostPort[0]
					chefPort = hostPort[1]
				} else if len(hostPort) == 1 {
					chefHost = hostPort[0]
					switch chefUrl.Scheme {
					case "http":
						chefPort = "80"
					case "https":
						chefPort = "443" 
					}
				} else {
					log.Println("Invalid Chef server host format.")
					return false
				}
			}
		}
	}

	ch.endpoint = fmt.Sprintf("%s://%s:%s/%s/%s/", chefUrl.Scheme, chefHost, chefPort, chefPath, chefOrg)


	if len(ch.key) <= 0 {
		log.Println("Missing Chef client key.")
		return false
	}

	if len(ch.user) <= 0 {
		log.Println("Missing Chef User_id.")
		return false
	}

	if len(ch.endpoint) <= 0 {
		log.Println("Missing Chef Server api endpoint.")
		return false
	}

	return true

}


func parse(text string) string {

	trimmed := strings.Trim(text, "\"")

	out := ""

	if strings.Contains(trimmed, "#ENV") {

		tmp := strings.Split(trimmed, "/")

		for _, part := range tmp {
			if strings.Contains(part, "#ENV") {
				substr := strings.Split(part, "'")
				envvar := strings.Split(substr[0], "'")
				newpart := os.Getenv(envvar[0])
				out += fmt.Sprintf("%s/", newpart)
			} else {
				out += fmt.Sprintf("%s/", part)
			}
		}


	} else {
		out = trimmed
	}

	return out

}



