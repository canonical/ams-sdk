List Applications Example
=========================

Demonstrates how to list the existing applications using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      | required |

Example:

    application-list -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443

Output:

[
    {
      "id": "bgutrvm5nof0fqm0894g",
      "name": "clash of clans",
      "status_code": 2,
      "status": "ready",
      "instance_type": "a2.3",
      "boot_package": "com.supercell.clashofclans",
      "parent_image_id": "bguqesm5nof675ard8b0",
      "published": true,
      "versions": [
          {
              "number": 0,
              "parent_image_version": 0,
              "status_code": 3,
              "status": "active",
              "published": true,
              "created_at": 1547558410,
              "boot_activity": "",
              "required_permissions": [
                  "android.permission.WRITE_EXTERNAL_STORAGE",
                  "android.permission.READ_EXTERNAL_STORAGE"
              ],
              "addons": [],
              "extra_data": {
                  "com.supercell.clashofclans.obb": {
                      "target": "/data/app/com.supercell.clashofclans-1/lib",
                      "owner": "",
                      "permissions": ""
                  },
                  "game-data-folder": {
                      "target": "/sdcard/Android/data/com.supercell.clashofclans/",
                      "owner": "",
                      "permissions": ""
                  }
              }
          }
      ],
      "addons": null,
      "created_at": 0,
      "immutable": false
    }
]

