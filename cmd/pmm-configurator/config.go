package main

import (
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
