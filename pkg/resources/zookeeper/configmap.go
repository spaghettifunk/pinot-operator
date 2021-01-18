package zookeeper

import (
	"fmt"
	"strconv"

	"github.com/hoisie/mustache"
	"github.com/spaghettifunk/pinot-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var okConfig = `
#!/bin/sh
zkServer.sh status
`

var readyConfig = `
#!/bin/sh
echo ruok | nc 127.0.0.1 ${1:-{{clientPort}}}
`

var runConfig = `
#!/bin/bash

set -a
ROOT=$(echo /apache-zookeeper-*)

ZK_USER={{user}}
ZK_LOG_LEVEL={{logLevel}}
ZK_DATA_DIR={{dataDir}}
ZK_DATA_LOG_DIR={{dataLogDir}}
ZK_CONF_DIR={{configDir}}
ZK_CLIENT_PORT={{clientPort}}
ZK_SERVER_PORT={{serverPort}}
ZK_ELECTION_PORT={{electionPort}}
ZK_TICK_TIME={{tickTime}}
ZK_INIT_LIMIT={{initLimit}}
ZK_SYNC_LIMIT={{syncLimit}}
ZK_MAX_CLIENT_CNXNS={{maxClientsConnections}}
ZK_MIN_SESSION_TIMEOUT={{minimumSessionTimeout}}
ZK_MAX_SESSION_TIMEOUT={{maxSessionTimeout}}
ZK_SNAP_RETAIN_COUNT={{snapRetailCount}}
ZK_PURGE_INTERVAL={{purgeInterval}}
ID_FILE="$ZK_DATA_DIR/{{zookeeperId}}"
ZK_CONFIG_FILE={{configFilePath}}
LOG4J_PROPERTIES={{log4jPropertiesPath}}
HOST=$(hostname)
DOMAIN=$(hostname -d)
JVMFLAGS="{{jvmOptions}}"

APPJAR=$(echo $ROOT/*jar)
CLASSPATH="${ROOT}/lib/*:${APPJAR}:${ZK_CONF_DIR}:"

if [[ $HOST =~ (.*)-([0-9]+)$ ]]; then
	NAME=${BASH_REMATCH[1]}
	ORD=${BASH_REMATCH[2]}
	MY_ID=$((ORD+1))
else
	echo "Failed to extract ordinal from hostname $HOST"
	exit 1
fi

mkdir -p $ZK_DATA_DIR
mkdir -p $ZK_DATA_LOG_DIR
echo $MY_ID >> $ID_FILE

echo "clientPort=$ZK_CLIENT_PORT" >> $ZK_CONFIG_FILE
echo "dataDir=$ZK_DATA_DIR" >> $ZK_CONFIG_FILE
echo "dataLogDir=$ZK_DATA_LOG_DIR" >> $ZK_CONFIG_FILE
echo "tickTime=$ZK_TICK_TIME" >> $ZK_CONFIG_FILE
echo "initLimit=$ZK_INIT_LIMIT" >> $ZK_CONFIG_FILE
echo "syncLimit=$ZK_SYNC_LIMIT" >> $ZK_CONFIG_FILE
echo "maxClientCnxns=$ZK_MAX_CLIENT_CNXNS" >> $ZK_CONFIG_FILE
echo "minSessionTimeout=$ZK_MIN_SESSION_TIMEOUT" >> $ZK_CONFIG_FILE
echo "maxSessionTimeout=$ZK_MAX_SESSION_TIMEOUT" >> $ZK_CONFIG_FILE
echo "autopurge.snapRetainCount=$ZK_SNAP_RETAIN_COUNT" >> $ZK_CONFIG_FILE
echo "autopurge.purgeInterval=$ZK_PURGE_INTERVAL" >> $ZK_CONFIG_FILE
echo "4lw.commands.whitelist=*" >> $ZK_CONFIG_FILE

for (( i=1; i<=$ZK_REPLICAS; i++ ))
do
	echo "server.$i=$NAME-$((i-1)).$DOMAIN:$ZK_SERVER_PORT:$ZK_ELECTION_PORT" >> $ZK_CONFIG_FILE
done

rm -f $LOG4J_PROPERTIES

echo "zookeeper.root.logger=$ZK_LOG_LEVEL, CONSOLE" >> $LOG4J_PROPERTIES
echo "zookeeper.console.threshold=$ZK_LOG_LEVEL" >> $LOG4J_PROPERTIES
echo "zookeeper.log.threshold=$ZK_LOG_LEVEL" >> $LOG4J_PROPERTIES
echo "zookeeper.log.dir=$ZK_DATA_LOG_DIR" >> $LOG4J_PROPERTIES
echo "zookeeper.log.file=zookeeper.log" >> $LOG4J_PROPERTIES
echo "zookeeper.log.maxfilesize=256MB" >> $LOG4J_PROPERTIES
echo "zookeeper.log.maxbackupindex=10" >> $LOG4J_PROPERTIES
echo "zookeeper.tracelog.dir=$ZK_DATA_LOG_DIR" >> $LOG4J_PROPERTIES
echo "zookeeper.tracelog.file=zookeeper_trace.log" >> $LOG4J_PROPERTIES
echo "log4j.rootLogger=\${zookeeper.root.logger}" >> $LOG4J_PROPERTIES
echo "log4j.appender.CONSOLE=org.apache.log4j.ConsoleAppender" >> $LOG4J_PROPERTIES
echo "log4j.appender.CONSOLE.Threshold=\${zookeeper.console.threshold}" >> $LOG4J_PROPERTIES
echo "log4j.appender.CONSOLE.layout=org.apache.log4j.PatternLayout" >> $LOG4J_PROPERTIES
echo "log4j.appender.CONSOLE.layout.ConversionPattern=%d{ISO8601} [{{zookeeperId}}:%X{{{zookeeperId}}}] - %-5p [%t:%C{1}@%L] - %m%n" >> $LOG4J_PROPERTIES

if [ -n "$JMXDISABLE" ]
then
	MAIN=org.apache.zookeeper.server.quorum.QuorumPeerMain
else
	MAIN="-Dcom.sun.management.jmxremote -Dcom.sun.management.jmxremote.port=$JMXPORT -Dcom.sun.management.jmxremote.authenticate=$JMXAUTH -Dcom.sun.management.jmxremote.ssl=$JMXSSL -Dzookeeper.jmx.log4j.disable=$JMXLOG4J org.apache.zookeeper.server.quorum.QuorumPeerMain"
fi

set -x
exec java -cp "$CLASSPATH" $JVMFLAGS $MAIN $ZK_CONFIG_FILE
`

const (
	user                  = "zookeeper"
	logLevel              = "INFO"
	dataDir               = "/data"
	dataLogDir            = "/data/log"
	configDir             = "/conf"
	configFilename        = "zoo.cfg"
	log4jFilename         = "log4j.properties"
	tickTime              = 2000
	initLimit             = 10
	syncLimit             = 5
	maxClientsConnections = 60
	snapRetailCount       = 3
	purgeInterval         = 0
	zookeeperID           = "myid"
)

func (r *Reconciler) configmap() runtime.Object {
	return &apiv1.ConfigMap{
		ObjectMeta: templates.ObjectMeta(configmapName, r.labels(), r.Config),
		Data: map[string]string{
			"ok": mustache.Render(okConfig, map[string]string{}),
			"ready": mustache.Render(readyConfig, map[string]string{
				"clientPort": strconv.Itoa(zookeeperClientPort),
			}),
			"run": mustache.Render(runConfig, map[string]string{
				"user":                  user,
				"logLevel":              logLevel,
				"dataDir":               dataDir,
				"dataLogDir":            dataLogDir,
				"configDir":             configDir,
				"initLimit":             strconv.Itoa(initLimit),
				"syncLimit":             strconv.Itoa(syncLimit),
				"tickTime":              strconv.Itoa(tickTime),
				"snapRetailCount":       strconv.Itoa(snapRetailCount),
				"purgeInterval":         strconv.Itoa(purgeInterval),
				"maxClientsConnections": strconv.Itoa(maxClientsConnections),
				"minimumSessionTimeout": strconv.Itoa(tickTime * 2),
				"maxSessionTimeout":     strconv.Itoa(tickTime * 20),
				"clientPort":            strconv.Itoa(zookeeperClientPort),
				"serverPort":            strconv.Itoa(zookeeperServerPort),
				"electionPort":          strconv.Itoa(zookeeperElectionPort),
				"configFilePath":        fmt.Sprintf("%s/%s", configDir, configFilename),
				"log4jPropertiesPath":   fmt.Sprintf("%s/%s", configDir, log4jFilename),
				"zookeeperId":           zookeeperID,
			}),
		},
	}
}
