Show Application Example
========================

Demonstrates how to show the information of an existing application using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server                   | required |
| `id`      | Identifier of the application to show   | required |


Example:

    application-show -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgutrvm5nof0fqm0894g

Output:

{
    "id": "bgutrvm5nof0fqm0894g",
    "name": "clash of clans",
    "status": "ready",
    "published": true,
    "immutable": false,
    "Config": {
        "instance-type": "a2.3",
        "boot-package": "com.supercell.clashofclans"
    },
    "versions": {
        "0": {
            "image": "bguqesm5nof675ard8b0 (version 0)",
            "published": true,
            "status": "active",
            "required-permissions": [
                "android.permission.WRITE_EXTERNAL_STORAGE",
                "android.permission.READ_EXTERNAL_STORAGE"
            ],
            "extra-data": {
                "com.supercell.clashofclans.obb": {
                    "target": "/data/app/com.supercell.clashofclans-1/lib"
                },
                "game-data-folder": {
                    "target": "/sdcard/Android/data/com.supercell.clashofclans/"
                }
            }
        }
    }
}


