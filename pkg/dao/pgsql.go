package dao

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"harbor-operator/config"
	"k8s.io/klog"
)

// SyncDatabases 用户创建所有的实例数据库以及添加用户权限
func SyncDatabases(instanceName string, opConfig *config.ConfigFile) error {
	conninfo := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s",
		opConfig.DB.PGSQL.Username,
		opConfig.DB.PGSQL.Password,
		opConfig.DB.PGSQL.Host,
		opConfig.DB.PGSQL.Port,
		opConfig.DB.PGSQL.SSLMode)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to init pg, witherr:%s", err.Error()))
		return err
	}
	defer db.Close()

	// 创建表的操作用不了数据库事务
	createRegistryDatabaseRaw := fmt.Sprintf(`CREATE DATABASE %s_registry ENCODING 'UTF8'`, instanceName)
	_, err = db.Exec(createRegistryDatabaseRaw)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create registry database of %s, witherr: %s", instanceName, err.Error()))
		return err
	}

	createNotaryServerDatabaseRaw := fmt.Sprintf(`CREATE DATABASE %s_notaryserver`, instanceName)
	_, err = db.Exec(createNotaryServerDatabaseRaw)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create notaryserver database of %s, witherr: %s", instanceName, err.Error()))
		return err
	}

	createNotarySignerDatabaseRaw := fmt.Sprintf(`CREATE DATABASE %s_notarysigner`, instanceName)
	_, err = db.Exec(createNotarySignerDatabaseRaw)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create notarysigner database of %s, witherr: %s", instanceName, err.Error()))
		return err
	}

	createClairDatabaseRaw := fmt.Sprintf(`CREATE DATABASE %s_clair ENCODING 'UTF8'`, instanceName)
	_, err = db.Exec(createClairDatabaseRaw)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create clair database of %s, witherr: %s", instanceName, err.Error()))
		return err
	}

	// 默认用户部署的pgsql已经创建好server用户，参考https://github.com/goharbor/harbor/blob/v2.1.3/make/photon/db/initial-notarysigner.sql
	grantNotaryserver := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s_notaryserver TO server", instanceName)
	_, err = db.Exec(grantNotaryserver)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to grant notaryserver of %s, witherr: %s", instanceName, err.Error()))
		return err
	}

	// 默认用户部署的pgsql已经创建好signer用户
	grantNotarysigner := fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s_notarysigner TO signer", instanceName)
	_, err = db.Exec(grantNotarysigner)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to grant notarysigner of %s, witherr: %s", instanceName, err.Error()))
		return err
	}

	return nil
}

// SyncTables 用于在registry库中创建表schema_migrations
func SyncTables(instanceName string, opConfig *config.ConfigFile) error {
	connRegistryinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=%s",
		opConfig.DB.PGSQL.Username,
		opConfig.DB.PGSQL.Password,
		opConfig.DB.PGSQL.Host,
		opConfig.DB.PGSQL.Port,
		fmt.Sprintf("%s_registry", instanceName),
		opConfig.DB.PGSQL.SSLMode)
	db, err := sql.Open("postgres", connRegistryinfo)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to init pg, witherr:%s", err.Error()))
		return err
	}
	defer db.Close()

	createSchemaRaw := `CREATE TABLE schema_migrations(version bigint not null primary key, dirty boolean not null)`

	_, err = db.Exec(createSchemaRaw)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to create registry schema table, witherr:%s", err.Error()))
		return err
	}
	return nil
}

// DeleteDatabase 用于清除相关数据库
func DeleteDatabases(instanceName string, opConfig *config.ConfigFile) error {
	conninfo := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=%s",
		opConfig.DB.PGSQL.Username,
		opConfig.DB.PGSQL.Password,
		opConfig.DB.PGSQL.Host,
		opConfig.DB.PGSQL.Port,
		opConfig.DB.PGSQL.SSLMode)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		klog.Error(fmt.Sprintf("fail to init pg, witherr:%s", err.Error()))
		return err
	}
	defer db.Close()

	// 断开registry库并删除
	disconnectRegistryRaw := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE datname='%s_registry' AND pid<>pg_backend_pid()", instanceName)
	_, errDisRegistry := db.Exec(disconnectRegistryRaw)
	if errDisRegistry != nil {
		klog.Error(fmt.Sprintf("fail to disconnect registry database of %s, witherr: %s", instanceName, errDisRegistry.Error()))
	}

	deleteRegistryDatabaseRaw := fmt.Sprintf(`DROP DATABASE %s_registry`, instanceName)
	_, errRegistry := db.Exec(deleteRegistryDatabaseRaw)
	if errRegistry != nil {
		klog.Error(fmt.Sprintf("fail to drop registry database of %s, witherr: %s", instanceName, errRegistry.Error()))
	}

	// 断开notaryserver库并删除
	disconnectNotaryserverRaw := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE datname='%s_notaryserver' AND pid<>pg_backend_pid()", instanceName)
	_, errDisServer := db.Exec(disconnectNotaryserverRaw)
	if errDisServer != nil {
		klog.Error(fmt.Sprintf("fail to disconnect notaryserver database of %s, witherr: %s", instanceName, errDisServer.Error()))
	}

	deleteNotaryserverDatabaseRaw := fmt.Sprintf(`DROP DATABASE %s_notaryserver`, instanceName)
	_, errNotaryserver := db.Exec(deleteNotaryserverDatabaseRaw)
	if errNotaryserver != nil {
		klog.Error(fmt.Sprintf("fail to drop notaryserver database of %s, witherr: %s", instanceName, errNotaryserver.Error()))
	}

	// 断开notarysigner库并删除
	disconnectNotarysignerRaw := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE datname='%s_notarysigner' AND pid<>pg_backend_pid()", instanceName)
	_, errDisSigner := db.Exec(disconnectNotarysignerRaw)
	if errDisServer != nil {
		klog.Error(fmt.Sprintf("fail to disconnect notarysigner database of %s, witherr: %s", instanceName, errDisSigner.Error()))
	}

	deleteNotarysigerDatabaseRaw := fmt.Sprintf(`DROP DATABASE %s_notarysigner`, instanceName)
	_, errNotarysigner := db.Exec(deleteNotarysigerDatabaseRaw)
	if errNotarysigner != nil {
		klog.Error(fmt.Sprintf("fail to drop notarysigner database of %s, witherr: %s", instanceName, errNotarysigner.Error()))
	}

	// 断开clair库并删除
	disconnectClairRaw := fmt.Sprintf("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE datname='%s_clair' AND pid<>pg_backend_pid()", instanceName)
	_, errDisClair := db.Exec(disconnectClairRaw)
	if errDisServer != nil {
		klog.Error(fmt.Sprintf("fail to disconnect clair database of %s, witherr: %s", instanceName, errDisClair.Error()))
	}

	deleteClairDatabaseRaw := fmt.Sprintf(`DROP DATABASE %s_clair`, instanceName)
	_, errClair := db.Exec(deleteClairDatabaseRaw)
	if errClair != nil {
		klog.Error(fmt.Sprintf("fail to drop clair database of %s, witherr: %s", instanceName, errClair.Error()))
	}

	if errRegistry != nil || errNotaryserver != nil || errNotarysigner != nil || errClair != nil {
		return errors.New("fail to clear pgsql databases")
	}

	return nil
}
