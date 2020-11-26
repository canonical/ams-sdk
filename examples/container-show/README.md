Show Container Example
======================

Demonstrates how to show a container information using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server           | required |
| `id`      | Identifier of the container     | required |


Example:

    container-show -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgvb70vb9s9jdk3dpujg

Output:

    {
        "id": "bgvb70vb9s9jdk3dpujg",
        "name": "ams-bgvb70vb9s9jdk3dpujg",
        "status": "running",
        "node": "lxd0",
        "created_at": "1970-01-01T00:00:00Z",
        "application": {
            "id": "bgvb6u7b9s9jdk3dpuj0"
        },
        "image": {},
        "network": {
            "address": "192.168.100.2",
            "public_address": "18.191.182.69",
            "services": null
        },
        "stored_logs": null,
        "error_message": ""
    }
