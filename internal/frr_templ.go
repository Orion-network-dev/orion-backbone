package internal

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type group struct {
	ASN int64
}

type Peer struct {
	Address string
	ASN     int64
}

type tmplContext struct {
	ASN     int64
	OrionId int64
	Groups  []group
	Peers   []Peer
}

var (
	tmplPath = flag.String("frr-config-template", "/etc/oriond/frr.conf.tmpl", "the configuration template file")
)

func loadTmpl() *template.Template {
	tmpl, err := template.New(*tmplPath).ParseFiles(*tmplPath)
	if err != nil {
		panic(err)
	}
	return tmpl
}

type FrrConfigManager struct {
	Peers    []Peer
	ASN      int64
	OrionId  int64
	template *template.Template
}

func NewFrrConfigManager(ASN int64, OrionId int64) (*FrrConfigManager, error) {
	config := &FrrConfigManager{
		Peers:    make([]Peer, 0),
		ASN:      ASN,
		OrionId:  OrionId,
		template: loadTmpl(),
	}
	err := config.Update()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *FrrConfigManager) Update() error {
	context := tmplContext{}
	context.ASN = c.ASN
	context.OrionId = c.OrionId
	context.Peers = c.Peers
	groups := make(map[int64]group, len(c.Peers))

	// Might be useful when applying policies, when we give the control to the user
	for _, peer := range context.Peers {
		groups[peer.ASN] = group{
			ASN: peer.ASN,
		}
	}

	tempConfig, err := os.CreateTemp("/tmp", "orion-conf-update*.conf")
	if err != nil {
		return err
	}
	defer tempConfig.Close()

	err = c.template.Execute(tempConfig, context)
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(tempConfig.Name())
	if err != nil {
		return err
	}

	// Apply the configuration patch
	execReload := exec.Command("/usr/lib/frr/frr-reload.py", "-reload", abs)
	if err := execReload.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); !ok {
			return fmt.Errorf("reload process returned failure %s", exitError)
		}
	}

	return nil
}
