package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type VcapServices struct {
	PConfigServer []struct {
		BindingGUID string      `json:"binding_guid"`
		BindingName interface{} `json:"binding_name"`
		Credentials struct {
			CredhubRef string `json:"credhub-ref"`
		} `json:"credentials"`
		InstanceGUID   string        `json:"instance_guid"`
		InstanceName   string        `json:"instance_name"`
		Label          string        `json:"label"`
		Name           string        `json:"name"`
		Plan           string        `json:"plan"`
		Provider       interface{}   `json:"provider"`
		SyslogDrainURL interface{}   `json:"syslog_drain_url"`
		Tags           []string      `json:"tags"`
		VolumeMounts   []interface{} `json:"volume_mounts"`
	} `json:"p.config-server"`
}

// IndexHandler returns a simple message
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World from Cloud Foundry!</h1>")

	v := os.Getenv("VCAP_SERVICES")
	if v != "" {
		var vcapServices VcapServices
		err := json.NewDecoder(strings.NewReader(v)).Decode(&vcapServices)
		if err != nil {
			fmt.Fprintf(w, "%s", err)
			w.WriteHeader(500)
			return
		}

		if len(vcapServices.PConfigServer) > 0 {
			p := vcapServices.PConfigServer[0]
			fmt.Fprintf(w, "<h2>p-config server:</h2>")
			fmt.Fprintf(w, "<p>name: %s</p>", p.Name)
			fmt.Fprintf(w, "<p>binding id: %s</p>", p.BindingGUID)
			fmt.Fprintf(w, "<p>binding instance id: %s</p>", p.InstanceGUID)
			fmt.Fprintf(w, "<p>instance name: %s</p>", p.InstanceName)
			fmt.Fprintf(w, "<p>label: %s</p>", p.Label)
			fmt.Fprintf(w, "<p>plan: %s</p>", p.Plan)
			//fmt.Fprintf(w, "<p>credentials: %v</p>", p.Credentials)
			fmt.Fprintf(w, "<p>syslog drain url: %v</p>", p.SyslogDrainURL)
			fmt.Fprintf(w, "<p>provider: %v</p>", p.Provider)
			fmt.Fprintf(w, "<p>tags: %v</p>", p.Tags)
			fmt.Fprintf(w, "<p>volume mounts: %v</p>", p.VolumeMounts)
		}
	}
}

func main() {
	http.HandleFunc("/", IndexHandler)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
