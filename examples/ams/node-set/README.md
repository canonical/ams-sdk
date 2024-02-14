Set Node Property Example
=========================

Demonstrates how to set a node property using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server                | required |
| `name`    | Name of the node             | required |
| `key`     | Key of the property to set   | required |
| `value`   | Value of the property to set | required |


Available keys:

- public-address
- cpus
- cpu-allocation-rate
- memory
- memory-allocation-rate

Example:

    node-set -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=lxd1 -key=public-address -value=174.56.55.2

Output:

empty
