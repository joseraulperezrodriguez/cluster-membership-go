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
		"--" + config.NodeServerPort, "",
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
		"--" + config.ConfigPath, "/tmp/path",
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
