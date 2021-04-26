# node-init

metal-stack uses frr and BGP for routing-to-the-host. If a calico interface / pod IP gets added or removed to the kernel routing table, frr announces or removes the according BGP route.

There is an open issue with frr, that sometimes - in rare cases - frr misses these route-changes: <https://github.com/FRRouting/frr/issues/7299> and the BGP routing tables differs from kernel routing table, leading to unreachable pods.

Since all pod IPs on a node are from the same podCIDR anyway, there is no need for host-routes for every pod. It's sufficient enough to simply route the whole node.podCIDR to that node.

`node-init` reads node.spec.podCIDR and sets a route to dev lo0. This results in a BGP route for this podCIDR to that node and therefore mitigates the issue with missing podIP routes.

## Usage

```text
Usage:
  node-init [command]

Available Commands:
  help        Help about any command
  init        init node networking

Flags:
  -h, --help   help for node-init

Use "node-init [command] --help" for more information about a command.
```

## Example

```bash
kubectl apply -f deploy/node-init.yaml
```
