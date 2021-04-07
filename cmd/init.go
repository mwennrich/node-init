package cmd

import (
	"context"
	"fmt"
	"net"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	netlink "github.com/vishvananda/netlink"
)

var (
	initNode = &cobra.Command{
		Use:   "init",
		Short: "init node networking",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initNetwork(args)
		},
	}
)

func init() {
	err := viper.BindPFlags(initNode.Flags())
	if err != nil {
		panic(err.Error())
	}
}

func initNetwork(_ []string) error {

	klog.Infoln("Starting node-init")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodeName := os.Getenv("KUBE_NODE_NAME")
	if len(nodeName) == 0 {
		return fmt.Errorf("node env is missing")
	}

	node, err := c.CoreV1().Nodes().Get(context.Background(), nodeName, v1.GetOptions{})
	if err != nil {
		return err
	}
	podCidrString := node.Spec.PodCIDR
	_, podCidr, err := net.ParseCIDR(podCidrString)
	if err != nil {
		return err
	}

	// check if route already exists
	link, err := netlink.LinkByName("lo")
	if err != nil {
		return err
	}
	routes, err := netlink.RouteList(link, netlink.FAMILY_V4)
	if err != nil {
		return err
	}
	for _, r := range routes {
		if r.Dst.String() == podCidrString {
			klog.Infof("route for %s already exists", podCidrString)
			return nil
		}
	}
	// add route
	route := netlink.Route{LinkIndex: link.Attrs().Index, Scope: netlink.SCOPE_LINK, Dst: podCidr}
	err = netlink.RouteAdd(&route)
	if err != nil {
		return err
	}
	klog.Infof("route for %s successfully added", podCidrString)
	return nil
}
