Show Node Example
=================

Demonstrates how to show the information of a node using AMS SDK

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

    node-show -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=lxd1

Output:
    {
        "name": "lxd0",
        "status": "online",
        "network": {
            "address": "172.31.29.146",
            "bridge-mtu": 9001
        },
        "config": {
            "public-address": "18.191.182.69",
            "cpus": 8,
            "cpu-allocation-rate": 4,
            "memory": "31GB",
            "memory-allocation-rate": 2,
        }
    }
