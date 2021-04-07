# node-init

***ALPHA version***
Reads node.spec.podCIDR and sets a route to dev lo0

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
