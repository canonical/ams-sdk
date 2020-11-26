Show Addon Example
==================

Demonstrates how to show an addon information using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      | required |
| `name`    | Name of the addon          | required |


Example:

    addon-show -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -name=debugger

Output:

	{
	    "name": "debugger",
	    "versions": {
		    "0": {
		        "size": "238B",
		        "created-at": "2018-09-17 11:15:44 +0000 UTC"
		    }
	    }
	}

