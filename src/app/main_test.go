package main

import (
	"cluster-membership-go/src/common"
	"cluster-membership-go/src/config"
	"cluster-membership-go/src/model"
	"flag"
	"os"
	"testing"
)

var mode, ids, addresses, protocolPorts, serverPorts = createFlags()

func TestMain(t *testing.T) {
	t.Run("args working okay", mainBasicTestOnGoodParameter)
	t.Run("args with 2 elements working okay", main2ElementTestOnGoodParameter)
}

func mainBasicTestOnGoodParameter(t *testing.T) {
	os.Args = []string{"-",
		"--mode", "RUNNING",
		"--" + config.NodeID, "id1",
		"--" + config.NodeAddress, "ad1",
		"--" + config.NodeProtocolPort, "20",
		"--" + config.NodeServerPort, "21",
	}
	//var mode, ids, addresses, protocolPorts, serverPorts = createFlags()
	flag.Parse()

	conf, configError := config.Validate(*mode, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError != nil {
		t.Error(&configError)
	}

	var expectedConfig = config.Configuration{
		Mode:  "RUNNING",
		Nodes: []model.Node{{ID: "id1", Address: "ad1", ProtocolPort: 20, ServerPort: 21}},
	}
	if conf.Mode != expectedConfig.Mode || !common.EqualsArrays(conf.Nodes, expectedConfig.Nodes) {
		t.Error(&common.BasicError{ErrorS: "mainTestOnGoodParameter Failed, some values don't match the expected values"})
		t.Logf("Expected config %v distinct from result config %v\n", expectedConfig, conf)
	}
}

func main2ElementTestOnGoodParameter(t *testing.T) {
	os.Args = []string{"-",
		"--mode", "RUNNING",
		"--" + config.NodeID, "id1,id2",
		"--" + config.NodeAddress, "ad1,ad2",
		"--" + config.NodeProtocolPort, "20,21",
		"--" + config.NodeServerPort, "22,23",
	}
	//var mode, ids, addresses, protocolPorts, serverPorts = createFlags()
	flag.Parse()

	conf, configError := config.Validate(*mode, *ids, *addresses, *protocolPorts, *serverPorts)

	if configError != nil {
		t.Error(&configError)
	}

	var expectedConfig = config.Configuration{
		Mode: "RUNNING",
		Nodes: []model.Node{
			{ID: "id1", Address: "ad1", ProtocolPort: 20, ServerPort: 22},
			{ID: "id2", Address: "ad2", ProtocolPort: 21, ServerPort: 23},
		},
	}
	if conf.Mode != expectedConfig.Mode || !common.EqualsArrays(conf.Nodes, expectedConfig.Nodes) {
		t.Error(&common.BasicError{ErrorS: "main2ElementTestOnGoodParameter Failed, some values don't match the expected values"})
		t.Logf("Expected config %v distinct from result config %v\n", expectedConfig, conf)
	}
}
