List Nodes Example
==================

Demonstrates how to list available nodes using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server                | required |

Example:

    node-list -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443

Output:

    [
        {
            "name": "lxd0",
            "address": "172.31.29.146",
            "public_address": "18.191.182.69",
            "network_bridge_mtu": 9001,
            "cpus": 8,
            "cpu_allocation_rate": 4,
            "memory": "31GB",
            "memory_allocation_rate": 2,
            "status_code": 4,
            "status": "online",
            "is_master": true,
            "disk_size": "13GB",
            "gpu_slots": 0
        }
    ]
