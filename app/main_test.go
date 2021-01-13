package main

import (
	"cluster-membership-go/common"
	"cluster-membership-go/config"
	"cluster-membership-go/model"
	"flag"
	"os"
	"strconv"
	"testing"
)

var mode *string
var configPath *string
var ids *config.ArgStringList
var addresses *config.ArgStringList
var protocolPorts *config.ArgStringList
var serverPorts *config.ArgStringList

func TestMain(t *testing.T) {
	mode, configPath, ids, addresses, protocolPorts, serverPorts = createFlags()
	t.Run("args basic should work ok", mainBasicTestOnGoodParameter)
	t.Run("args with 2 elements should work ok", main2ElementTestOnGoodParameter)
	t.Run("args with empty parameter should fail", basicTestOnBadParameter)
	t.Run("args with wrong int format should fail", errorOnIntParse)
	t.Run("args with different list size should fail", errorOnDifferentListSize)
	t.Run("args basic with resources file for current node should be okay", goodReadingPropertiesFile)

}

func mainBasicTestOnGoodParameter(t *testing.T) {
	os.Args = []string{"-",
		"--" + config.AppMode, "RUNNING",
		"--" + config.ConfigPath, "path/to/file",

		"--" + config.NodeID, "id1",
		"--" + config.NodeAddress, "ad1",
		"--" + config.NodeProtocolPort, "20",
		"--" + config.NodeServerPort, "21",
	}
	flag.Parse()

	conf, configError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError != nil {
		t.Error(&common.TraceableError{ErrorS: "Unexpected error due to the provided test case", RawError: configError})
	}

	var expectedConfig = config.Configuration{
		Mode:      "RUNNING",
		SeedNodes: []model.Node{{ID: "id1", Address: "ad1", ProtocolPort: 20, ServerPort: 21}},
	}
	if conf.Mode != expectedConfig.Mode || !common.EqualsArrays(conf.SeedNodes, expectedConfig.SeedNodes) {
		t.Error(&common.BaseError{ErrorS: "mainTestOnGoodParameter Failed, some values don't match the expected values"})
		t.Logf("Expected config %v distinct from result config %v\n", expectedConfig, conf)
	}
}

func main2ElementTestOnGoodParameter(t *testing.T) {
	os.Args = []string{"-",
		"--" + config.AppMode, "RUNNING",
		"--" + config.ConfigPath, "path/to/file",
		"--" + config.NodeID, "id1,id2",
		"--" + config.NodeAddress, "ad1,ad2",
		"--" + config.NodeProtocolPort, "20,21",
		"--" + config.NodeServerPort, "22,23",
	}
	flag.Parse()

	conf, configError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError != nil {
		t.Error(&common.TraceableError{ErrorS: "Unexpected error due to the provided test case", RawError: configError})
	}

	var expectedConfig = config.Configuration{
		Mode: "RUNNING",
		SeedNodes: []model.Node{
			{ID: "id1", Address: "ad1", ProtocolPort: 20, ServerPort: 22},
			{ID: "id2", Address: "ad2", ProtocolPort: 21, ServerPort: 23},
		},
	}
	if conf.Mode != expectedConfig.Mode || !common.EqualsArrays(conf.SeedNodes, expectedConfig.SeedNodes) {
		t.Error(&common.BaseError{ErrorS: "main2ElementTestOnGoodParameter Failed, some values don't match the expected values"})
		t.Logf("Expected config %v distinct from result config %v\n", expectedConfig, conf)
	}
}

func basicTestOnBadParameter(t *testing.T) {
	os.Args = []string{"-",
		"--" + config.AppMode, "RUNNING",
		"--" + config.ConfigPath, "path/to/file",
		"--" + config.NodeID, "id1,id2",
		"--" + config.NodeAddress, "ad1,ad2",
		"--" + config.NodeProtocolPort, "20,21",
		"--" + config.NodeServerPort, "2456",
	}
	flag.Parse()

	_, configError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError == nil {
		t.Error(&common.BaseError{ErrorS: "Not error detected in bad config"})
	} else {
		_, ok := configError.(*common.ArgumentParsingError)
		if !ok {
			t.Error(&common.BaseError{ErrorS: "Error detected don't match the expected error type"})
		}
	}
}

func errorOnIntParse(t *testing.T) {
	os.Args = []string{"-",
		"--" + config.AppMode, "RUNNING",
		"--" + config.ConfigPath, "path/to/file",
		"--" + config.NodeID, "id1,id2",
		"--" + config.NodeAddress, "ad1,ad2",
		"--" + config.NodeProtocolPort, "20,21",
		"--" + config.NodeServerPort, "22,ii23d",
	}
	flag.Parse()

	_, configError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError == nil {
		t.Error(&common.BaseError{ErrorS: "Not error detected in bad config"})
	} else {
		_, ok := configError.(*strconv.NumError)
		if !ok {
			t.Error(&common.BaseError{ErrorS: "Error detected don't match the expected error type"})
		}
	}
}

func errorOnDifferentListSize(t *testing.T) {
	os.Args = []string{"-",
		"--" + config.AppMode, "RUNNING",
		"--" + config.ConfigPath, "path/to/some",
		"--" + config.NodeID, "id1",
		"--" + config.NodeAddress, "ad1,ad2",
		"--" + config.NodeProtocolPort, "20,21",
		"--" + config.NodeServerPort, "22,23",
	}
	flag.Parse()

	_, configError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError == nil {
		t.Error(&common.BaseError{ErrorS: "Not error detected in bad config"})
	} else {
		_, ok := configError.(*common.ArgumentParsingError)
		if !ok {
			t.Error(&common.BaseError{ErrorS: "Error detected don't match the expected error type"})
		}
	}
}

func goodReadingPropertiesFile(t *testing.T) {
	os.Args = []string{"-",
		"--" + config.AppMode, "RUNNING",
		"--" + config.ConfigPath, "../resources/app.properties",

		"--" + config.NodeID, "id1",
		"--" + config.NodeAddress, "ad1",
		"--" + config.NodeProtocolPort, "20",
		"--" + config.NodeServerPort, "21",
	}
	flag.Parse()

	conf, configError := config.Validate(*mode, *configPath, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError != nil {
		t.Error(&common.TraceableError{ErrorS: "Unexpected error due to the provided test case", RawError: configError})
	}

	var expectedConfig = config.Configuration{
		Mode:      "RUNNING",
		SeedNodes: []model.Node{{ID: "id1", Address: "ad1", ProtocolPort: 20, ServerPort: 21}},
	}
	if conf.Mode != expectedConfig.Mode || !common.EqualsArrays(conf.SeedNodes, expectedConfig.SeedNodes) {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, some values don't match the expected values"})
		t.Logf("Expected config %v distinct from result config %v\n", expectedConfig, conf)
	}

	configError = config.AddProgramConfig(&conf)

	if conf.IterationIntervalMs != 3000 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, IterationIntervalMs not match"})
	}

	if conf.ConnectionTimeoutMs != 1000 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, ConnectionTimeoutMs not match"})
	}

	if conf.ReadIddleIterationFactor != 3 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, ReadIddleIterationFactor not match"})
	}

	if conf.CyclesForWaitKeepAlive != 3 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, CyclesForWaitKeepAlive not match"})
	}

	if conf.MaxRumorLogSize != 1000000 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, MaxRumorLogSize not match"})
	}

	if conf.MaxObjectSize != 2147483647 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, MaxObjectSize not match"})
	}

	if conf.ClientThreads != 3 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, ClientThreads not match"})
	}

	if conf.ServerThreads != 3 {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, ServerThreads not match"})
	}

	var expected = model.Node{ID: "A", Address: "localhost", ProtocolPort: 7001, ServerPort: 6001}
	if expected != conf.CurrentNode {
		t.Error(&common.BaseError{ErrorS: "goodReadingPropertiesFile Failed, CurrentNode don't match"})
	}

}
