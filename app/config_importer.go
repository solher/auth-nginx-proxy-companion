package app

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/solher/auth-nginx-proxy-companion/models"
	"gopkg.in/yaml.v2"
)

type (
	ConfigImporterPoliciesInter interface {
		Create(policy *models.Policy) (*models.Policy, error)
	}

	ConfigImporterPoliciesValidator interface {
		ValidateCreation(policy *models.Policy) error
	}

	ConfigImporterResourcesInter interface {
		Create(resource *models.Resource) (*models.Resource, error)
	}

	ConfigImporterResourcesValidator interface {
		ValidateCreation(resource *models.Resource) error
	}

	ConfigImporter struct {
		pi ConfigImporterPoliciesInter
		pv ConfigImporterPoliciesValidator
		ri ConfigImporterResourcesInter
		rv ConfigImporterResourcesValidator
	}
)

type Config struct {
	Resources []models.Resource `json:"resources" yaml:"resources"`
	Policies  []models.Policy   `json:"policies" yaml:"policies"`
}

func NewConfigImporter(
	pi ConfigImporterPoliciesInter,
	pv ConfigImporterPoliciesValidator,
	ri ConfigImporterResourcesInter,
	rv ConfigImporterResourcesValidator,
) *ConfigImporter {
	return &ConfigImporter{pi: pi, pv: pv, ri: ri, rv: rv}
}

func (ci *ConfigImporter) Import(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	config := &Config{}

	switch {
	case strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml"):
		config, err = ci.fromYAML(file)
	case strings.HasSuffix(path, ".json"):
		config, err = ci.fromJSON(file)
	}

	if err != nil {
		return err
	}

	if config.Resources != nil {
		for _, resource := range config.Resources {
			if err := ci.rv.ValidateCreation(&resource); err != nil {
				return err
			}

			if _, err := ci.ri.Create(&resource); err != nil {
				return err
			}
		}
	}

	if config.Policies != nil {
		for _, policy := range config.Policies {
			if err := ci.pv.ValidateCreation(&policy); err != nil {
				return err
			}

			if _, err := ci.pi.Create(&policy); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ci *ConfigImporter) fromJSON(conf []byte) (*Config, error) {
	config := &Config{}

	if err := json.Unmarshal(conf, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (ci *ConfigImporter) fromYAML(conf []byte) (*Config, error) {
	config := &Config{}

	if err := yaml.Unmarshal(conf, config); err != nil {
		return nil, err
	}

	return config, nil
}
