Delete Image Example
====================

Demonstrates how to delete an image or an image version using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server               | required |
| `id`      | Id or name of the image to delete   | required |
| `version` | Version of the image to delete      | optional |


Example:

    image-delete -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgvarifb9s9hpbfhgcdg

Output:

    empty
