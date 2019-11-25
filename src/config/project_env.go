package config

import (
	"cluster-membership-go/src/common"
	"cluster-membership-go/src/model"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/magiconair/properties"
)

//Configuration type
type Configuration struct {
	Mode       string
	SeedNodes  []model.Node
	ConfigPath string

	IterationIntervalMs      int
	ConnectionTimeoutMs      int
	ReadIddleIterationFactor int
	CyclesForWaitKeepAlive   int
	MaxRumorLogSize          int
	MaxObjectSize            int
	ClientThreads            int
	ServerThreads            int

	CurrentNode model.Node
}

//NilConfig for return when needed
var NilConfig = Configuration{}

//NodeID is the command line argument name for seeds node id
var NodeID = "id"

//NodeAddress is the command line argument name for seeds node address
var NodeAddress = "address"

//NodeProtocolPort is the command line argument name for seeds node protocol port
var NodeProtocolPort = "protocol.port"

//NodeServerPort is the command line argument name for seeds node server port
var NodeServerPort = "server.port"

//AppMode is the command line argument that tells us if it is executing in test or production mode
var AppMode = "mode"
var runningMode = "RUNNING"
var testMode = "TEST"

//ConfigPath is the path for program configuration file
var ConfigPath = "config.path"

//ArgListSep is the string used to parse the command line parameters that contains list of elements
var ArgListSep = ","

//ArgStringList implements
type ArgStringList struct {
	Values []string
}

//IsBoolFlag implements the Interface Value, from package flag
func (arg *ArgStringList) IsBoolFlag() bool {
	return false
}

//Set implements the Interface Value, from package flag
func (arg *ArgStringList) Set(argS string) error {
	var array = strings.Split(argS, ",")
	if len(array) == 0 {
		return &common.BaseError{ErrorS: "Field is not valid, no value provided"}
	}
	arg.Values = array
	return nil
}

//String implements the Interface Value, from package flag
func (arg *ArgStringList) String() string {
	return common.ArrayToString(arg.Values, ArgListSep)
}

//Validate validates the command line input
func Validate(mode string, configPath string, ids ArgStringList, addresses ArgStringList,
	protocolPorts ArgStringList, serverPorts ArgStringList) (Configuration, error) {
	if mode != strings.ToUpper(runningMode) && mode != strings.ToUpper(testMode) {
		return Configuration{}, &common.ArgumentParsingError{ErrorS: fmt.Sprintf("The 'mode' variable must be either %s or %s", runningMode, testMode)}
	}
	if strings.Trim(configPath, " ") == "" {
		return Configuration{}, &common.ArgumentParsingError{ErrorS: "Empty configPath"}
	}
	if len(ids.Values) != len(addresses.Values) || len(addresses.Values) != len(protocolPorts.Values) ||
		len(protocolPorts.Values) != len(serverPorts.Values) {
		return NilConfig, &common.ArgumentParsingError{ErrorS: fmt.Sprintf("Arguments %s, %s, %s and %s must contains the same number of element", NodeID, NodeAddress, NodeProtocolPort, NodeServerPort)}
	}
	if chk := common.CheckNonEmpty(ids.Values, NodeID); chk != nil {
		return NilConfig, chk
	} else if chk := common.CheckNonEmpty(addresses.Values, NodeAddress); chk != nil {
		return NilConfig, chk
	} else if chk := common.CheckNonEmpty(protocolPorts.Values, NodeProtocolPort); chk != nil {
		return NilConfig, chk
	} else if chk := common.CheckNonEmpty(serverPorts.Values, NodeServerPort); chk != nil {
		return NilConfig, chk
	} else {
		return buildConf(mode, configPath, ids, addresses, protocolPorts, serverPorts)
	}
}

func buildConf(mode string, configPath string, ids ArgStringList, addresses ArgStringList,
	protocolPorts ArgStringList, serverPorts ArgStringList) (Configuration, error) {
	var nodes = make([]model.Node, len(ids.Values))
	for i := 0; i < len(ids.Values); i++ {
		var pPort, pError = strconv.Atoi(protocolPorts.Values[i])
		if pError != nil {
			return NilConfig, pError
		}
		var sPort, sError = strconv.Atoi(serverPorts.Values[i])
		if sError != nil {
			return NilConfig, sError
		}
		nodes[i] = model.Node{ID: ids.Values[i], Address: addresses.Values[i], ProtocolPort: pPort, ServerPort: sPort}
	}

	absConfigPath, absConfigPathErr := filepath.Abs(configPath)
	if absConfigPathErr != nil {
		return NilConfig, absConfigPathErr
	}
	var config = &Configuration{SeedNodes: nodes, Mode: mode, ConfigPath: absConfigPath}
	return *config, nil
}

//AddProgramConfig adds the configuration parameter from the config file
func AddProgramConfig(config *Configuration) error {
	prop := properties.MustLoadFile(config.ConfigPath, properties.UTF8)

	config.IterationIntervalMs = prop.GetInt("iteration.interval.ms", 3000)
	config.ConnectionTimeoutMs = prop.GetInt("connection.timeout.ms", 1000)
	config.ReadIddleIterationFactor = prop.GetInt("read.iddle.iteration.factor", 3)
	config.CyclesForWaitKeepAlive = prop.GetInt("cycles.for.wait.keep.alive", 3)
	config.MaxRumorLogSize = prop.GetInt("max.rumor.log.size", 1000000)
	config.MaxObjectSize = prop.GetInt("max.object.size", 2147483647)
	config.ClientThreads = prop.GetInt("client.threads", 3)
	config.ServerThreads = prop.GetInt("server.threads", 3)
	config.CurrentNode = model.Node{
		ID:           strings.Trim(prop.GetString(NodeID, ""), " "),
		Address:      strings.Trim(prop.GetString(NodeAddress, ""), " "),
		ProtocolPort: prop.GetInt(NodeProtocolPort, -1),
		ServerPort:   prop.GetInt(NodeServerPort, -1),
	}

	if config.CurrentNode.ID == "" {
		return &common.BaseError{ErrorS: fmt.Sprintf("%v can't have empty values", NodeID)}
	}
	if config.CurrentNode.Address == "" {
		return &common.BaseError{ErrorS: fmt.Sprintf("%v can't have empty values", NodeAddress)}
	}
	if config.CurrentNode.ProtocolPort < 0 {
		return &common.BaseError{ErrorS: fmt.Sprintf("%v should be setted", NodeProtocolPort)}
	}
	if config.CurrentNode.ServerPort < 0 {
		return &common.BaseError{ErrorS: fmt.Sprintf("%v should be setted", NodeServerPort)}
	}

	return nil
}
