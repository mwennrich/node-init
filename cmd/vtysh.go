package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"

	"k8s.io/klog/v2"
)

type (
	BGPSummary struct {
		Ipv4Unicast Ipv4Unicast `json:"ipv4Unicast"`
	}
	Ipv4Unicast struct {
		RouterID    string          `json:"routerId"`
		Peers       map[string]Peer `json:"peers"`
		FailedPeers int             `json:"failedPeers"`
		TotalPeers  int             `json:"totalPeers"`
	}
	Peer struct {
		Hostname                   string `json:"hostname"`
		RemoteAs                   int64  `json:"remoteAs"`
		LocalAs                    int64  `json:"localAs"`
		Version                    int    `json:"version"`
		MsgRcvd                    int    `json:"msgRcvd"`
		MsgSent                    int    `json:"msgSent"`
		TableVersion               int    `json:"tableVersion"`
		Outq                       int    `json:"outq"`
		Inq                        int    `json:"inq"`
		PeerUptime                 string `json:"peerUptime"`
		PeerUptimeMsec             int    `json:"peerUptimeMsec"`
		PeerUptimeEstablishedEpoch int    `json:"peerUptimeEstablishedEpoch"`
		PfxRcd                     int    `json:"pfxRcd"`
		PfxSnt                     int    `json:"pfxSnt"`
		State                      string `json:"state"`
		PeerState                  string `json:"peerState"`
		ConnectionsEstablished     int    `json:"connectionsEstablished"`
		ConnectionsDropped         int    `json:"connectionsDropped"`
		IDType                     string `json:"idType"`
	}
)

func repairFailedBGPSession() error {
	failedInterfaces, err := checkBGPSessions()
	if err != nil {
		return err
	}

	for _, name := range failedInterfaces {
		klog.Infof("bgp session on interface %q is broken, try to repair", name)
		err := ethtool(name)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkBGPSessions() ([]string, error) {
	socketPath, err := lookupSocketPath("bgpd")
	if err != nil {
		return nil, err
	}
	output, err := runCmd(socketPath, "show bgp ipv4 summary json")
	if err != nil {
		return nil, err
	}

	var bgpSummary BGPSummary
	err = json.Unmarshal(output, &bgpSummary)
	if err != nil {
		return nil, err
	}

	if bgpSummary.Ipv4Unicast.FailedPeers == 0 {
		return nil, nil
	}

	failedInterfaces := []string{}
	for name, peer := range bgpSummary.Ipv4Unicast.Peers {
		// FIXME check with real failing interface bgp session output
		if peer.State != "Established" {
			failedInterfaces = append(failedInterfaces, name)
		}
	}

	return failedInterfaces, nil
}

func lookupSocketPath(daemon string) (string, error) {
	switch daemon {
	case
		"babeld",
		"bfdd",
		"bgpd",
		"eigrpd",
		"fabricd",
		"isisd",
		"ldpd",
		"nhrpd",
		"ospf6d",
		"ospfd",
		"pbrd",
		"pimd",
		"ripd",
		"ripngd",
		"sharpd",
		"staticd",
		"vrrpd",
		"zebra":
		return fmt.Sprintf("/var/run/frr/%s.vty", daemon), nil
	}
	return "", fmt.Errorf("unknown daemon %s", daemon)
}

func runCmd(socketPath string, cmd string) ([]byte, error) {
	socket, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	cmd = cmd + "\x00"
	_, err = socket.Write([]byte(cmd))
	if err != nil {
		return nil, err
	}

	output, err := bufio.NewReader(socket).ReadBytes('\x00')
	if err != nil {
		return nil, err
	}

	return output[:len(output)-1], nil
}

func ethtool(iface string) error {
	for _, onoff := range []string{"on", "off"} {
		c := fmt.Sprintf("ethtool --pause %s rx %s", iface, onoff)
		cmd := exec.Command("bash", "-c", c)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error executing %q error:%w", c, err)
		}
	}
	return nil
}
