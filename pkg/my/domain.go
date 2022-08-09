package my

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const domainFile = "deployers.json"

type Domain struct {
	URL     string `json:"dnsName"`
	Address string `json:"address"`
}

type Domains struct {
	file    string
	domains []Domain
}

func (d *Domains) List() []Domain {
	return d.domains
}

func (d *Domains) Save() error {
	cnt, err := json.Marshal(d.domains)
	if err != nil {
		return err
	}

	err = os.WriteFile("deployers.json", cnt, 0o600) // u+r, u+w
	if err != nil {
		return err
	}

	return nil
}

func (d *Domains) Add(dom Domain) {
	d.domains = append(d.domains, dom)
}

func NewDomains() (*Domains, error) {
	file, err := ioutil.ReadFile(domainFile)
	if err != nil {
		return nil, err
	}

	dom := Domains{file: "deployers.json", domains: []Domain{}}

	err = json.Unmarshal(file, &dom.domains)
	if err != nil {
		return nil, err
	}

	return &dom, nil
}
