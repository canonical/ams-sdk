Publish Application Example
===========================

Demonstrates how to publish an application using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      | required |
| `id`      | Identifier of the application to publish  | required |
| `version` | version of the application to publish     | required |


Example:

    application-publish -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgutrvm5nof0fqm0894g -version=0

Output:

empty

