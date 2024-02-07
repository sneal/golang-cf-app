package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const serverAddressKey = "ServerAddress"

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

func index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	out(w, "<h1>Hello World from Cloud Foundry!</h1>")
	out(w, "<p>listening on: %s</p>", ctx.Value(serverAddressKey))

	v := os.Getenv("VCAP_SERVICES")
	if v != "" {
		var vcapServices VcapServices
		err := json.NewDecoder(strings.NewReader(v)).Decode(&vcapServices)
		if err != nil {
			out(w, "%s", err)
			w.WriteHeader(500)
			return
		}

		if len(vcapServices.PConfigServer) > 0 {
			p := vcapServices.PConfigServer[0]
			out(w, "<h2>p-config server:</h2>")
			out(w, "<p>name: %s</p>", p.Name)
			out(w, "<p>binding id: %s</p>", p.BindingGUID)
			out(w, "<p>binding instance id: %s</p>", p.InstanceGUID)
			out(w, "<p>instance name: %s</p>", p.InstanceName)
			out(w, "<p>label: %s</p>", p.Label)
			out(w, "<p>plan: %s</p>", p.Plan)
			//out(w, "<p>credentials: %v</p>", p.Credentials)
			out(w, "<p>syslog drain url: %v</p>", p.SyslogDrainURL)
			out(w, "<p>provider: %v</p>", p.Provider)
			out(w, "<p>tags: %v</p>", p.Tags)
			out(w, "<p>volume mounts: %v</p>", p.VolumeMounts)
		}
	}
}

func out(w http.ResponseWriter, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}
	var tcpPort string
	if tcpPort = os.Getenv("TCP_PORT"); len(tcpPort) == 0 {
		tcpPort = "32000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	ctx, cancelCtx := context.WithCancel(context.Background())
	webServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, serverAddressKey, l.Addr().String())
			return ctx
		},
	}
	tcpServer := &http.Server{
		Addr:    ":" + tcpPort,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, serverAddressKey, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := webServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("web server closed\n")
		} else if err != nil {
			fmt.Printf("error listening for web server: %s\n", err)
		}
		cancelCtx()
	}()

	go func() {
		err := tcpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("TCP server closed\n")
		} else if err != nil {
			fmt.Printf("error listening for TCP server: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}
