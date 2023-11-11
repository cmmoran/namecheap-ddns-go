package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
)

func main() {
	configPathPtr := flag.String("config", "", "Path to the namecheap-ddns-go config")

	flag.Parse()

	if len(*configPathPtr) == 0 {
		if cfg, ok := os.LookupEnv("NAMECHEAP_DDNS_CONFIG"); ok {
			configPathPtr = &cfg
		} else {
			panic("NAMECHEAP_DDNS_CONFIG is not set")
		}
	}

	hostConfigs, err := readConfig(*configPathPtr)
	if err != nil {
		panic(errors.Wrap(err, "failed to get config"))
	}

	for _, config := range hostConfigs.Configs {

		for _, subdomain := range config.Subdomains {

			url := fmt.Sprintf("https://dynamicdns.park-your-domain.com/update?domain=%s&host=%s&password=%s",
				config.Domain,
				subdomain,
				config.Token)

			if len(config.IP) > 0 {
				url = fmt.Sprintf("%s&ip=%s", url, config.IP)
			}

			resp, err := http.Get(url)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to send DNS update request %v\n", err)
			}

			if resp.StatusCode >= 400 {
				body, _ := io.ReadAll(resp.Body)
				_, _ = fmt.Fprintf(os.Stderr, "got error response from DNS server: status_code %d, status: %s, response_body: %s\n", resp.StatusCode, resp.Status, body)
			} else {
				body, _ := io.ReadAll(resp.Body)
				dec := xml.NewDecoder(bytes.NewReader(body))
				dec.CharsetReader = identReader
				decresp := &Response{}
				if err = dec.Decode(decresp); err != nil {
					panic("could not parse response")
				}
				_, _ = fmt.Printf("%s.%s IP address updated to: %s\n", subdomain, config.Domain, decresp.IP)
			}
		}
	}
}

func identReader(_ string, input io.Reader) (io.Reader, error) {
	return input, nil
}

type Response struct {
	IP string `xml:"IP"`
}

type HostConfigs struct {
	Configs []Config `yaml:"configs"`
}

type Config struct {
	Domain     string   `yaml:"domain"`
	Subdomains []string `yaml:"subdomains"`
	Token      string   `yaml:"token"`
	IP         string   `yaml:"ip,omitempty"`
	LogLevel   string   `yaml:"log_level,omitempty"`
}

func readConfig(path string) (*HostConfigs, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	result := &HostConfigs{
		Configs: make([]Config, 0),
	}

	if err = yaml.Unmarshal(configBytes, result); err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	return result, nil
}
