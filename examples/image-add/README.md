Add Image Example
=================

Demonstrates how to add a new image using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server              | required |
| `name`    | Name of the image                  | required |
| `package-path`  | Path to the image xz tarball | required |


Example:

    image-add -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=default -package-path=./anbox-lxd-bionic-e0261ed-2018-09-17-0_arm64.tar.xz

Output:

    {
		"images": [
			"bgvarifb9s9hpbfhgcdg"
		]
	}
