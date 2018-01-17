package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

// PMMDBInstance is autodiscovered RDS instance with all intrested fields for PMM
type PMMDBInstance struct {
	DBInstanceIdentifier string `json:"DBInstanceIdentifier"` // it should be same in mysqld_exporter and cloudwatch_exporter
	EndpointAddress      string `json:"EndpointAddress"`      // DNS hostname, can point to diffenet hosts in case maintenance, restoration, or AZ disaster
	MasterUsername       string `json:"MasterUsername"`       // user can use separate user for PMM
	Region               string `json:"Region"`               // important for Enhanced Monitoring
	MonitoringInterval   int64  `json:"MonitoringInterval"`   // important for Enhanced Monitoring
}

func getRDS(w http.ResponseWriter, req *http.Request) {
	var result []PMMDBInstance
	for _, region := range endpoints.AwsPartition().Services()[endpoints.RdsServiceID].Regions() {
		session := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(region.ID()),
		}))
		svc := rds.New(session)

		DBInstancesOutput, err := svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{})
		if err != nil {
			returnError(w, req, http.StatusInternalServerError, "Cannot find RDS instances", err)
			return
		}

		for _, DBInstance := range DBInstancesOutput.DBInstances {
			result = append(result, PMMDBInstance{
				DBInstanceIdentifier: *DBInstance.DBInstanceIdentifier,
				EndpointAddress:      *DBInstance.Endpoint.Address,
				MasterUsername:       *DBInstance.MasterUsername,
				Region:               region.ID(),
				MonitoringInterval:   *DBInstance.MonitoringInterval,
			})
		}
	}
	json.NewEncoder(w).Encode(result) // nolint: errcheck
}
