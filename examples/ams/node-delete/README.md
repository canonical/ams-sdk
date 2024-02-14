Delete Node Example
===================

Demonstrates how to delete a node using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server                | required |
| `name`    | Name of the node                     | required |


Example:

    node-delete -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=lxd1

Output:

empty
