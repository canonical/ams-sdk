Show Image Example
==================

Demonstrates how to show an image information using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server               | required |
| `id`      | ID or name of the image to show     | required |


Example:

    image-show -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgvarifb9s9hpbfhgcdg

Output:

    {
        "id": "bgvarifb9s9hpbfhgcdg",
        "name": "default",
        "status": "active",
        "versions": {
            "0": {
            "size": "280.26MB",
            "created-at": "2019-01-16 04:06:33 +0000 UTC",
            "status": "active"
            }
        }
    }