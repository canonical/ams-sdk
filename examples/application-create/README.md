Create Application Example
==========================

Demonstrates how to create a new application using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      | required |
| `name`    | Name of the addon          | required |
| `package-path` | Path to the application tarball  | required |

The package must be a valid tar.bz2 file with this struture:

    * Application_folder
        - manifest.yaml
        - app.apk

Where:

* `Application_folder` Folder with the same name of the application to install
* `manifest.yaml` contains metadata for the application. Including:
    - `name` The desired name for the new application
    - `image` The id or name of an existing image. It is used as the base for the application image.
    - `instance-type` Size of the machines in terms of memory and CPU to use for the application containers
    - `addons` List of addons this application uses
    - `extra-data` Any additional data needed by the application, provided in a free format
* `app.apk` The Android installable file for this application. Must have exactly that _app.apk_ name


Example of application package content

    clashofclans
	├── app.apk
	└── manifest.yaml

Example of manifest.yaml:

    name: Clash of Clans
    image: default
    instance-type: a2.3
    required_permissions: []
    addons: 
    - debugger
    extra-data:


Example:

    application-create -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -package-path=./clashofclans.tar.bz2
    
Output:

	{
		"applications": [
			"bgutrvm5nof0fqm0894g"
		]
	}

