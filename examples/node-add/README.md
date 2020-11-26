Add Node Example
================

Demonstrates how to add a new node using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server                | required |
| `name`    | Name of the node                     | required |
| `address` | Address where the node is accessible | required |
| `trust-password` | Password for accessing to LXD server | optional |
| `storage-device` | The pool or dataset used for storage | optional |
| `network-bridge-mtu` | Network largest packet size      | optional |


Example:

    node-add -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=node1 -address=192.168.0.8

Output:

    {
		"nodes": [
			"lxd1"
		]
	}
