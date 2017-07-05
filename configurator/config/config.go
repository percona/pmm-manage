package config

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// PMMConfig implements struct with all configuration params in one place
type PMMConfig struct {
	ConfigPath           string              `yaml:"config"                 default:""                        desc:"configuration file location"`
	HtpasswdPath         string              `yaml:"htpasswd-path"          default:"/srv/nginx/.htpasswd"    desc:"htpasswd file location"`
	ListenAddress        string              `yaml:"listen-address"         default:"127.0.0.1:7777"          desc:"Address and port to listen on: [ip_address]:port"`
	PathPrefix           string              `yaml:"url-prefix"             default:"/configurator"           desc:"Prefix for the internal routes of web endpoints"`
	SSHKeyPath           string              `yaml:"ssh-key-path"           default:""                        desc:"authorized_keys file location"`
	SSHKeyOwner          string              `yaml:"ssh-key-owner"          default:"admin"                   desc:"Owner of authorized_keys file"`
	GrafanaDBPath        string              `yaml:"grafana-db-path"        default:"/srv/grafana/grafana.db" desc:"grafana database location"`
	PrometheusConfPath   string              `yaml:"prometheus-conf-path"   default:"/etc/prometheus.yml"     desc:"prometheus configuration file location"`
	UpdateDirPath        string              `yaml:"update-dir-path"        default:"/srv/update"             desc:"update directory location"`
	LogFilePath          string              `yaml:"log-file"               default:"/var/log/pmm-manage.log" desc:"log file location"`
	SkipPrometheusReload string              `yaml:"skip-prometheus-reload" default:"false"                   desc:"log file location"`
	Configuration        map[string]string   `yaml:"configuration"          default:""                        desc:""`
	Users                []map[string]string `yaml:"users"                  default:""                        desc:""`
}

// ParseConfig implements function which read command line arguments, configuration file and set default values
func ParseConfig() (c PMMConfig) {
	t := reflect.TypeOf(&c).Elem()
	v := reflect.ValueOf(&c).Elem()
	// iterate over all confConfig fields
	for i := 0; i < v.NumField(); i++ {
		// get string with argument description, like "Owner of authorized_keys file"
		descTag := t.Field(i).Tag.Get("desc")
		// skip config file only fields
		if descTag != "" {
			// get poiner to confConfig field, like &c.SSHKeyOwner
			valueAddr := v.Field(i).Addr().Interface().(*string)
			// get string with argument name, like "ssh-key-owner"
			yamlTag := t.Field(i).Tag.Get("yaml")
			// pass pointer, argument name and argument description to flag library
			flag.StringVar(valueAddr, yamlTag, "", descTag)
		}
	}

	flag.Parse()
	c.parseConfig()
	flag.Parse() // command line should overide config
	c.setDefaultValues()
	c.setLogger()
	c.validateValues()

	return c
}

func (c *PMMConfig) setDefaultValues() {
	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	// iterate over all confConfig fields
	for i := 0; i < v.NumField(); i++ {
		// get string with current value of field (which have been read from config)
		curValue := v.Field(i).String()
		// get string with default value of field
		defValue := t.Field(i).Tag.Get("default")

		if curValue == "" {
			v.Field(i).SetString(defValue)
		}
	}
}

func (c *PMMConfig) parseConfig() {
	// parseConfig() runs before setDefaultValues(), so it is needed to set default manually
	if c.ConfigPath == "" {
		c.ConfigPath = "/srv/update/pmm-manage.yml"
	}

	configBytes, err := ioutil.ReadFile(c.ConfigPath)
	if os.IsNotExist(err) {
		// ignore config file is not exists
		return
	}
	if err != nil {
		log.WithFields(log.Fields{
			"file":  c.ConfigPath,
			"error": err,
		}).Error("Cannot read config file")
		return
	}

	err = yaml.Unmarshal(configBytes, &c)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  c.ConfigPath,
			"error": err,
		}).Error("Cannot parse config file")
		return
	}
}

// Save dump configuration values to configuration file
func (c *PMMConfig) Save() error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Cannot encode configuration")
		return err
	}

	if err = ioutil.WriteFile(c.ConfigPath, bytes, 0644); err != nil {
		log.WithFields(log.Fields{
			"file":  c.ConfigPath,
			"error": err,
		}).Error("Cannot save configuration file")
		return err
	}
	return nil
}

func (c *PMMConfig) setLogger() {
	log.SetFormatter(&log.TextFormatter{DisableColors: true})
	if logFile, err := os.OpenFile(c.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil { // nolint: gas
		log.WithFields(log.Fields{
			"file":  c.LogFilePath,
			"error": err,
		}).Error("Failed to log to file, using default stderr")
	} else {
		log.SetOutput(io.MultiWriter(logFile, os.Stderr))
	}
}

func (c *PMMConfig) validateValues() {
	if len(c.PathPrefix) > 0 && !strings.HasPrefix(c.PathPrefix, "/") {
		c.PathPrefix = "/" + c.PathPrefix
		log.Warning("Prefix has been changed to " + c.PathPrefix)
	}
}
