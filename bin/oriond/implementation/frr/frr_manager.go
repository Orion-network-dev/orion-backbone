package frr

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	frrTmplPath     = flag.String("frr-config-template", "/etc/oriond/frr.conf.tmpl", "Configuration template for configuring FRR")
	frrReloadBinary = flag.String("frr-reload-script", "/usr/lib/frr/frr-reload.py", "Path to the frr-reload.py utility executable")
)

// Struct containing the information regarding someone ASN on the network.
type group struct {
	ASN uint32
}

// Struct containing a peer in the orion network, this is typically
// on-per-connection.
type Peer struct {
	Address string
	ASN     uint32
}

// Contest representing th entire state of the frr config file.
type tmplContext struct {
	ASN     uint32
	OrionId uint32
	Groups  []group
	Peers   []Peer
}

// Struct used to interact and issue the FRR configuration file updates.
type FrrConfigManager struct {
	Peers    map[uint32]*Peer
	selfASN  uint32
	OrionId  uint32
	template *template.Template
}

// Function used to load the template files.
func loadTmpl() (*template.Template, error) {
	content, err := os.ReadFile(*frrTmplPath)
	if err != nil {
		return nil, err
	}
	return template.New(*frrTmplPath).Parse(string(content))
}

// Function used to create a config-manager instance.
func NewFrrConfigManager(ASN uint32, OrionId uint32) (*FrrConfigManager, error) {
	tmpl, err := loadTmpl()
	if err != nil {
		return nil, err
	}
	config := &FrrConfigManager{
		Peers:    map[uint32]*Peer{},
		selfASN:  ASN,
		OrionId:  OrionId,
		template: tmpl,
	}

	if err = config.Update(); err != nil {
		return nil, err
	}

	return config, nil
}

// This functions takes the simplified BGP configuration stored in memory
// and renders a new configuration for the `frr` daemon used to implement
// `bgpd` which is the BGP daemon used through the Orion network.
func (c *FrrConfigManager) Update() error {
	log.Debug().Msg("updating the FRR configuration")
	peers := []Peer{}

	for _, value := range c.Peers {
		if value != nil {
			peers = append(peers, *value)
		}
	}

	// We initialize some basic information regarding the BGP sessions
	// configured for FRR's bgpd.
	context := tmplContext{
		ASN:     c.selfASN,
		OrionId: c.OrionId,
		Peers:   peers,
		Groups:  []group{},
	}

	// From the list of peers we infer the list of bgp groups to connect to.
	for _, peer := range context.Peers {
		if peer.ASN != 0 {
			log.Debug().Uint32("asn", peer.ASN).Msg("new group computed")
			context.Groups = append(context.Groups, group{
				ASN: peer.ASN,
			})
		}
	}

	// In order to call the frr-reload.py script, we must render the configuration
	// to a file, we choosed to use a temporary file for this purpose.
	tempConfig, err := os.CreateTemp("/tmp", "orion-conf-update*.conf")
	if err != nil {
		log.Error().Err(err).Msg("failed to create the new frr configuration temporary file")
		return err
	}
	// Since all temporary files are in the /tmp directory, we can simply infer the full path.
	absolutePath := tempConfig.Name()

	defer func() {
		log.Debug().Msg("cleaning up temporary files")
		// We do not delete the temporary files when running in debug mode.
		if zerolog.GlobalLevel() != zerolog.DebugLevel {
			os.Remove(absolutePath)
		}
		tempConfig.Close()
	}()

	// We execute the template with the current context.
	err = c.template.Execute(tempConfig, context)
	if err != nil {
		log.Debug().Err(err).Msg("failed to template the config file")
		return err
	}

	log.Debug().Str("config-file", absolutePath).Msg("running frr-reload.py")
	execReload := exec.Command(*frrReloadBinary, "--reload", absolutePath)
	// For debug purposes, we might want to link the frr-reload.py script stdout and stderr to
	// the current process ones.
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		execReload.Stdout = os.Stdout
		execReload.Stderr = os.Stderr
	}

	if err := execReload.Run(); err != nil {
		// In case of an errornous error code.
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("the frr-reload.py script exited with the code %s", exitError)
		}
	}

	return nil
}
