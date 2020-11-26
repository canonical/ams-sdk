List Container Logs Example
========================

Demonstrates how to list log files of a container

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server          | required |
| `id`      | Identifier of the container    | required |

*NOTE*:  These log files only exist when the container runs into error state.


Example:

    container-list-log -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgva3qu5nof0fqm089cg

Output:

    [
	    "system.log",
	    "container.log",
	    "android.log",
	    "console.log"
    ]
