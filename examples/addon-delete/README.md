Delete Addon Example
====================

Demonstrates how to delete an addon using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      |  required |
| `name`    | Name of the addon          |  required |
| `version` | Version of addon           |  optional |


Example:

    addon-delete -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=debugger

Output:

empty
