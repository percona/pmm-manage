package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type confConfig struct {
	ConfigPath         string `yaml:"config"               default:""                        desc:"configuration file location"`
	HtpasswdPath       string `yaml:"htpasswd-path"        default:"/srv/nginx/.htpasswd"    desc:"htpasswd file location"`
	ListenAddress      string `yaml:"listen-address"       default:"127.0.0.1:7777"          desc:"Address and port to listen on: [ip_address]:port"`
	PathPrefix         string `yaml:"url-prefix"           default:"/configurator"           desc:"Prefix for the internal routes of web endpoints"`
	SSHKeyPath         string `yaml:"ssh-key-path"         default:""                        desc:"authorized_keys file location"`
	SSHKeyOwner        string `yaml:"ssh-key-owner"        default:"admin"                   desc:"Owner of authorized_keys file"`
	GrafanaDBPath      string `yaml:"grafana-db-path"      default:"/srv/grafana/grafana.db" desc:"grafana database location"`
	PrometheusConfPath string `yaml:"prometheus-conf-path" default:"/etc/prometheus.yml"     desc:"prometheus configuration file location"`
	UpdateDirPath      string `yaml:"update-dir-path"      default:"/srv/update"             desc:"update directory location"`
}

func parseFlag() {
	t := reflect.TypeOf(&c).Elem()
	v := reflect.ValueOf(&c).Elem()
	// iterate over all confConfig fields
	for i := 0; i < v.NumField(); i++ {
		// get poiner to confConfig field, like &c.SSHKeyOwner
		valueAddr := v.Field(i).Addr().Interface().(*string)
		// get string with argument name, like "ssh-key-owner"
		yaml := string(t.Field(i).Tag.Get("yaml"))
		// get string with argument description, like "Owner of authorized_keys file"
		desc := string(t.Field(i).Tag.Get("desc"))
		// pass pointer, argument name and argument description to flag library
		flag.StringVar(valueAddr, yaml, "", desc)
	}

	flag.Parse()
	parseConfig()
	flag.Parse() // command line should overide config

	setDefaultValues()
	runSSHKeyChecks()
}

func setDefaultValues() {
	t := reflect.TypeOf(&c).Elem()
	v := reflect.ValueOf(&c).Elem()
	// iterate over all confConfig fields
	for i := 0; i < v.NumField(); i++ {
		// get string with current value of field (which have been read from config)
		curValue := v.Field(i).String()
		// get string with default value of field
		defValue := string(t.Field(i).Tag.Get("default"))

		if curValue == "" {
			v.Field(i).SetString(defValue)
		}
	}
}

func parseConfig() {
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
		errorStr := fmt.Sprintf("Cannot read '%s' config file: %s\n", c.ConfigPath, err)
		log.Fatal(errorStr)
	}

	err = yaml.Unmarshal(configBytes, &c)
	if err != nil {
		errorStr := fmt.Sprintf("Cannot parse '%s' config file: %s\n", c.ConfigPath, err)
		log.Fatal(errorStr)
	}
}
