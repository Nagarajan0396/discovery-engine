package libs

import (
	"database/sql"
	"errors"

	"github.com/accuknox/auto-policy-discovery/src/types"
)

// ================= //
// == Network Log == //
// ================= //

// LastFlowID network flow between [ startTime <= time < endTime ]
var LastFlowID int64 = 0

// ==================== //
// == Network Policy == //
// ==================== //

func GetNetworkPolicies(cfg types.ConfigDB, cluster, namespace, status, nwtype, rule string) []types.KnoxNetworkPolicy {
	results := []types.KnoxNetworkPolicy{}

	if cfg.DBDriver == "mysql" {
		docs, err := GetNetworkPoliciesFromMySQL(cfg, cluster, namespace, status, nwtype, rule)
		if err != nil {
			return results
		}
		results = docs
	} else if cfg.DBDriver == "sqlite3" {
		docs, err := GetNetworkPoliciesFromSQLite(cfg, cluster, namespace, status)
		if err != nil {
			return results
		}
		results = docs
	}

	return results
}

func GetNetworkPoliciesBySelector(cfg types.ConfigDB, cluster, namespace, status string, selector map[string]string) ([]types.KnoxNetworkPolicy, error) {
	results := []types.KnoxNetworkPolicy{}

	if cfg.DBDriver == "mysql" {
		docs, err := GetNetworkPoliciesFromMySQL(cfg, cluster, namespace, status, "", "")
		if err != nil {
			return nil, err
		}
		results = docs
	} else if cfg.DBDriver == "sqlite3" {
		docs, err := GetNetworkPoliciesFromSQLite(cfg, cluster, namespace, status)
		if err != nil {
			return nil, err
		}
		results = docs
	} else {
		return results, nil
	}

	filtered := []types.KnoxNetworkPolicy{}
	for _, policy := range results {
		matched := true
		for k, v := range selector {
			val := policy.Spec.Selector.MatchLabels[k]
			if val != v {
				matched = false
				break
			}
		}

		if matched {
			filtered = append(filtered, policy)
		}
	}

	return filtered, nil
}

func UpdateOutdatedNetworkPolicy(cfg types.ConfigDB, outdatedPolicy string, latestPolicy string) {
	if cfg.DBDriver == "mysql" {
		if err := UpdateOutdatedNetworkPolicyFromMySQL(cfg, outdatedPolicy, latestPolicy); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := UpdateOutdatedNetworkPolicyFromSQLite(cfg, outdatedPolicy, latestPolicy); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func UpdateNetworkPolicies(cfg types.ConfigDB, policies []types.KnoxNetworkPolicy) {
	for _, policy := range policies {
		UpdateNetworkPolicy(cfg, policy)
	}
}

func UpdateNetworkPolicy(cfg types.ConfigDB, policy types.KnoxNetworkPolicy) {
	if cfg.DBDriver == "mysql" {
		if err := UpdateNetworkPolicyToMySQL(cfg, policy); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := UpdateNetworkPolicyToSQLite(cfg, policy); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func InsertNetworkPolicies(cfg types.ConfigDB, policies []types.KnoxNetworkPolicy) {
	if cfg.DBDriver == "mysql" {
		if err := InsertNetworkPoliciesToMySQL(cfg, policies); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := InsertNetworkPoliciesToSQLite(cfg, policies); err != nil {
			log.Error().Msg(err.Error())
		}
	}

}

// ================ //
// == System Log == //
// ================ //

// LastSyslogID system log between [ startTime <= time < endTime ]
var LastSyslogID int64 = 0

// ================== //
// == System Alert == //
// ================== //

// LastSysAlertID system_alert between [ startTime <= time < endTime ]
var LastSysAlertID int64 = 0

// =================== //
// == System Policy == //
// =================== //

func UpdateOutdatedSystemPolicy(cfg types.ConfigDB, outdatedPolicy string, latestPolicy string) {
	if cfg.DBDriver == "mysql" {
		if err := UpdateOutdatedNetworkPolicyFromMySQL(cfg, outdatedPolicy, latestPolicy); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := UpdateOutdatedNetworkPolicyFromSQLite(cfg, outdatedPolicy, latestPolicy); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func GetSystemPolicies(cfg types.ConfigDB, namespace, status string) []types.KnoxSystemPolicy {
	results := []types.KnoxSystemPolicy{}

	if cfg.DBDriver == "mysql" {
		docs, err := GetSystemPoliciesFromMySQL(cfg, namespace, status)
		if err != nil {
			return results
		}
		results = docs
	} else if cfg.DBDriver == "sqlite3" {
		docs, err := GetSystemPoliciesFromSQLite(cfg, namespace, status)
		if err != nil {
			return results
		}
		results = docs
	}

	return results
}

func InsertSystemPolicies(cfg types.ConfigDB, policies []types.KnoxSystemPolicy) {
	if cfg.DBDriver == "mysql" {
		if err := InsertSystemPoliciesToMySQL(cfg, policies); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := InsertSystemPoliciesToSQLite(cfg, policies); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func UpdateSystemPolicy(cfg types.ConfigDB, policy types.KnoxSystemPolicy) {
	if cfg.DBDriver == "mysql" {
		if err := UpdateSystemPolicyToMySQL(cfg, policy); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := UpdateSystemPolicyToSQLite(cfg, policy); err != nil {
			log.Error().Msg(err.Error())
		}
	}

}

func GetWorkloadProcessFileSet(cfg types.ConfigDB, wpfs types.WorkloadProcessFileSet) (map[types.WorkloadProcessFileSet][]string, types.PolicyNameMap, error) {
	if cfg.DBDriver == "mysql" {
		res, pnMap, err := GetWorkloadProcessFileSetMySQL(cfg, wpfs)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		return res, pnMap, err
	} else if cfg.DBDriver == "sqlite3" {
		res, pnMap, err := GetWorkloadProcessFileSetSQLite(cfg, wpfs)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		return res, pnMap, err
	}
	return nil, nil, errors.New("no db driver")
}

func InsertWorkloadProcessFileSet(cfg types.ConfigDB, wpfs types.WorkloadProcessFileSet, fs []string) error {
	if cfg.DBDriver == "mysql" {
		return InsertWorkloadProcessFileSetMySQL(cfg, wpfs, fs)
	} else if cfg.DBDriver == "sqlite3" {
		return InsertWorkloadProcessFileSetSQLite(cfg, wpfs, fs)
	}
	return errors.New("no db driver")
}

func ClearWPFSDb(cfg types.ConfigDB, wpfs types.WorkloadProcessFileSet, duration int64) error {
	if cfg.DBDriver == "mysql" {
		return ClearWPFSDbMySQL(cfg, wpfs, duration)
	} else if cfg.DBDriver == "sqlite3" {
		return ClearWPFSDbSQLite(cfg, wpfs, duration)
	}
	return errors.New("no db driver")
}

// =========== //
// == Table == //
// =========== //

func ClearDBTables(cfg types.ConfigDB) {
	if cfg.DBDriver == "mysql" {
		if err := ClearDBTablesMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := ClearDBTablesSQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func ClearNetworkDBTable(cfg types.ConfigDB) {
	if cfg.DBDriver == "mysql" {
		if err := ClearNetworkDBTableMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func CreateTablesIfNotExist(cfg types.ConfigDB) {
	if cfg.DBDriver == "mysql" {
		if err := CreateTableNetworkPolicyMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableSystemPolicyMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableWorkLoadProcessFileSetMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableSystemLogsMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableNetworkLogsMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreatePolicyTableMySQL(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
	} else if cfg.DBDriver == "sqlite3" {
		if err := CreateTableNetworkPolicySQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableSystemPolicySQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableWorkLoadProcessFileSetSQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableSystemLogsSQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateTableNetworkLogsSQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreatePolicyTableSQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
		if err := CreateSystemSummaryTableSQLite(cfg); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

// =================== //
// == Observability == //
// =================== //
func UpdateOrInsertKubearmorLogs(cfg types.ConfigDB, kubearmorLogMap map[types.KubeArmorLog]int) error {
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		err = UpdateOrInsertKubearmorLogsMySQL(cfg, kubearmorLogMap)
	} else if cfg.DBDriver == "sqlite3" {
		err = UpdateOrInsertKubearmorLogsSQLite(cfg, kubearmorLogMap)
	}
	return err
}

func GetKubearmorLogs(cfg types.ConfigDB, filterLog types.KubeArmorLog) ([]types.KubeArmorLog, []uint32, error) {
	kubearmorLog := []types.KubeArmorLog{}
	totalCount := []uint32{}
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		kubearmorLog, totalCount, err = GetSystemLogsMySQL(cfg, filterLog)
	} else if cfg.DBDriver == "sqlite3" {
		kubearmorLog, totalCount, err = GetSystemLogsSQLite(cfg, filterLog)
	}
	return kubearmorLog, totalCount, err
}

func UpdateOrInsertCiliumLogs(cfg types.ConfigDB, ciliumLogs []types.CiliumLog) error {
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		err = UpdateOrInsertCiliumLogsMySQL(cfg, ciliumLogs)
	} else if cfg.DBDriver == "sqlite3" {
		err = UpdateOrInsertCiliumLogsSQLite(cfg, ciliumLogs)
	}
	return err
}

func GetCiliumLogs(cfg types.ConfigDB, ciliumFilter types.CiliumLog) ([]types.CiliumLog, []uint32, error) {
	ciliumLogs := []types.CiliumLog{}
	ciliumTotalCount := []uint32{}
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		ciliumLogs, ciliumTotalCount, err = GetCiliumLogsMySQL(cfg, ciliumFilter)
	} else if cfg.DBDriver == "sqlite3" {
		ciliumLogs, ciliumTotalCount, err = GetCiliumLogsSQLite(cfg, ciliumFilter)
	}
	return ciliumLogs, ciliumTotalCount, err
}

func GetPodNames(cfg types.ConfigDB, filter types.ObsPodDetail) ([]string, error) {
	res := []string{}
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		res, err = GetPodNamesMySQL(cfg, filter)
	} else if cfg.DBDriver == "sqlite3" {
		res, err = GetPodNamesSQLite(cfg, filter)
	}
	return res, err
}

// =============== //
// == Policy DB == //
// =============== //
func GetPolicyYamls(cfg types.ConfigDB, policyType string) ([]types.PolicyYaml, error) {
	var err error
	var results []types.PolicyYaml

	if cfg.DBDriver == "mysql" {
		results, err = GetPolicyYamlsMySQL(cfg, policyType)
		if err != nil {
			return nil, err
		}
	} else if cfg.DBDriver == "sqlite3" {
		results, err = GetPolicyYamlsSQLite(cfg, policyType)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func UpdateOrInsertPolicyYamls(cfg types.ConfigDB, policies []types.PolicyYaml) error {
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		err = UpdateOrInsertPolicyYamlsMySQL(cfg, policies)
	} else if cfg.DBDriver == "sqlite3" {
		err = UpdateOrInsertPolicyYamlsSQLite(cfg, policies)
	}
	return err
}

// ============= //
// == Summary == //
// ============= //
func UpsertSystemSummary(cfg types.ConfigDB, summaryMap map[types.SystemSummary]types.SysSummaryTimeCount) error {
	var err = errors.New("unknown db driver")
	if cfg.DBDriver == "mysql" {
		err = UpsertSystemSummaryMySQL(cfg, summaryMap)
	} else if cfg.DBDriver == "sqlite3" {
		err = UpsertSystemSummarySQLite(cfg, summaryMap)
	}
	return err
}

func upsertSysSummarySQL(db *sql.DB, summary types.SystemSummary, timeCount types.SysSummaryTimeCount) error {
	queryString := `cluster_name = ? and cluster_id = ? and namespace_name = ? and namespace_id = ? and container_name = ? and container_image = ? 
					and container_id = ? and podname = ? and operation = ? and labels = ? and deployment_name = ? and source = ? and destination = ? 
					and destination_namespace = ? and destination_labels = ? and type = ? and ip = ? and port = ? and protocol = ? and action = ?`

	query := "UPDATE " + TableSystemSummarySQLite + " SET count=count+?, updated_time=? WHERE " + queryString + " "

	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	result, err := updateStmt.Exec(
		timeCount.Count,
		timeCount.UpdatedTime,
		summary.ClusterName,
		summary.ClusterId,
		summary.NamespaceName,
		summary.NamespaceId,
		summary.ContainerName,
		summary.ContainerImage,
		summary.ContainerID,
		summary.PodName,
		summary.Operation,
		summary.Labels,
		summary.Deployment,
		summary.Source,
		summary.Destination,
		summary.DestNamespace,
		summary.DestLabels,
		summary.NwType,
		summary.IP,
		summary.Port,
		summary.Protocol,
		summary.Action,
	)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err == nil && rowsAffected == 0 {

		insertQueryString := `(cluster_name,cluster_id,namespace_name,namespace_id,container_name,container_image,container_id,podname,operation,labels,deployment_name,
				source,destination,destination_namespace,destination_labels,type,ip,port,protocol,action,count,updated_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

		insertQuery := "INSERT INTO " + TableSystemSummarySQLite + insertQueryString

		insertStmt, err := db.Prepare(insertQuery)
		if err != nil {
			return err
		}
		defer insertStmt.Close()

		_, err = insertStmt.Exec(
			summary.ClusterName,
			summary.ClusterId,
			summary.NamespaceName,
			summary.NamespaceId,
			summary.ContainerName,
			summary.ContainerImage,
			summary.ContainerID,
			summary.PodName,
			summary.Operation,
			summary.Labels,
			summary.Deployment,
			summary.Source,
			summary.Destination,
			summary.DestNamespace,
			summary.DestLabels,
			summary.NwType,
			summary.IP,
			summary.Port,
			summary.Protocol,
			summary.Action,
			timeCount.Count,
			timeCount.UpdatedTime)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
	}

	return nil
}
