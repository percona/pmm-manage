# PMM Manage

[![Build Status](https://travis-ci.org/percona/pmm-manage.svg?branch=master)](https://travis-ci.org/percona/pmm-manage)
[![Go Report Card](https://goreportcard.com/badge/github.com/percona/pmm-manage)](https://goreportcard.com/report/github.com/percona/pmm-manage)
[![CLA assistant](https://cla-assistant.io/readme/badge/percona/pmm-manage)](https://cla-assistant.io/percona/pmm-manage)

* Website: https://www.percona.com/doc/percona-monitoring-and-management/index.html
* Forum: https://www.percona.com/forums/questions-discussions/percona-monitoring-and-management/

PMM Manage is a tool for configuring options inside Percona Monitoring and Management (PMM) Server.

PMM Manage provides several key features:
* add/list/modify/remove web users
* add/list Pubic Key for SSH user access

## Building
```
export GOPATH=$(pwd)
go get -u github.com/percona/pmm-manage/cmd/pmm-configurator
ls -la bin/pmm-configurator
```
