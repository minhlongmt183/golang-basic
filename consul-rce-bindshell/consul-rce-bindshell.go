package consul_rce_bindshell

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	resp, err := http.Get("http://192.168.67.142:8500/v1/agent/self")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
}
