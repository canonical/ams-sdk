List Images Example
===================

Demonstrates how to list available images using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server               | required |

Example:

    image-list -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443

Output:

    [
        {
            "id": "bgvarifb9s9hpbfhgcdg",
            "name": "default",
            "versions": [
                {
                    "version": 0,
                    "fingerprint": "8174a75351c660008c1db3ff17c4bfe78774bc020ba581397576fff349b095b1",
                    "size": 293870712,
                    "created_at": 1547611593,
                    "status_code": 3,
                    "status": "active"
                }
            ],
            "status_code": 3,
            "status": "active",
            "used_by": [
                "bgvb6u7b9s9jdk3dpuj0"
            ]
        }
    ]


