# PMM Manage

* Website: https://www.percona.com/doc/percona-monitoring-and-management/index.html
* Forum: https://www.percona.com/forums/questions-discussions/percona-monitoring-and-management/

PMM Manage is a tool for configuring options inside Percona Monitoring and Management (PMM) Server.

PMM Manage provides several key features:
* add/list/modify/remove web users
* add/list Pubic Key for SSH user access

## Building
```
export GOPATH=$(pwd)
go get github.com/Percona-Lab/pmm-manage/cmd/pmm-configurator
ls -la bin/pmm-configurator
```
