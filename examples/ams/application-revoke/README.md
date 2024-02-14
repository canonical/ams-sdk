Revoke Application Example
==========================

Demonstrates how to revoke an application using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      | required |
| `id`      | Identifier of the application to revoke   | required |
| `version` | version of the application to revoke      | required |


Example:

    application-revoke -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgutrvm5nof0fqm0894g -version=0

Output:

empty

