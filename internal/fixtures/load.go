package fixtures

import (
	"io/ioutil"
	"text/template"
)

var (
	ControlPlane       *template.Template
	EastWestGateway    *template.Template
	ExposeService      *template.Template
	NameSpace          *template.Template
	RafayRemoteSecrete *template.Template
	HelloWorld         *template.Template
)

// Load loads fixtures
func Load() (err error) {
	ControlPlane, err = load("controlplane.yaml")
	if err != nil {
		return err
	}
	EastWestGateway, err = load("eastwest-gateway.yaml")
	if err != nil {
		return err
	}
	ExposeService, err = load("expose-service.yaml")
	if err != nil {
		return err
	}
	NameSpace, err = load("namespace-template.yaml")
	if err != nil {
		return err
	}
	RafayRemoteSecrete, err = load("rafayremote-secret.yaml")
	if err != nil {
		return err
	}
	HelloWorld, err = load("helloworld.yaml")
	if err != nil {
		return err
	}
	return nil
}

func load(fileName string) (*template.Template, error) {
	f, err := Fixtures.Open(fileName)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(fileName).Parse(string(b))
	if err != nil {
		return nil, err
	}

	return tmpl, err
}
