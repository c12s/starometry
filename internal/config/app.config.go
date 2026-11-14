package config

import (
	"os"
	"strings"
)

type AppConfig struct {
	serverPort       string
	nodeExporterAddr string
	cAdvisorAddr     string
	natsAddr         string
	grpcPort         string
	nodeID           string
}

func NewAppConfigFromEnv() *AppConfig {
	return &AppConfig{
		serverPort:       os.Getenv("APP_PORT"),
		nodeExporterAddr: hostPort("NODE_EXPORTER_URL", "NODE_EXPORTER_PORT"),
		cAdvisorAddr:     hostPort("CADVISOR_URL", "CADVISOR_PORT"),
		natsAddr:         natsAddress(),
		grpcPort:         os.Getenv("GRPC_PORT"),
		nodeID:           "",
	}
}

func (ap *AppConfig) GetServerPort() string {
	return ap.serverPort
}

func (ap *AppConfig) GetNodeExporterAddress() string {
	return ap.nodeExporterAddr
}

func (ap *AppConfig) GetCAdvisorAddress() string {
	return ap.cAdvisorAddr
}

func (ap *AppConfig) GetNatsAddress() string {
	return ap.natsAddr
}

func (ap *AppConfig) GetGRPCPort() string {
	return ap.grpcPort
}

func (ap *AppConfig) GetNodeID() string {
	return ap.nodeID
}

func (ap *AppConfig) SetNodeID(nodeID string) {
	nodeID = strings.TrimSpace(nodeID)
	if nodeID == "" {
		return
	}
	ap.nodeID = nodeID
}

// hostPort joins a host and a port variable; an unset host means the source
// is not configured (bootc nodes run node-exporter but no cAdvisor).
func hostPort(hostVar, portVar string) string {
	host := os.Getenv(hostVar)
	if host == "" {
		return ""
	}
	return host + ":" + os.Getenv(portVar)
}

// natsAddress prefers NATS_ADDRESS (host:port, the key star and the bootc
// image's node.env use) over the NATS_URL/NATS_PORT pair.
func natsAddress() string {
	if addr := strings.TrimSpace(os.Getenv("NATS_ADDRESS")); addr != "" {
		return addr
	}
	return os.Getenv("NATS_URL") + ":" + os.Getenv("NATS_PORT")
}
