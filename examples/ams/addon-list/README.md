List Addons Example
===================

Demonstrates how to list the installed addons using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server      | required |


Example:

    addon-list -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 

Output:

	[
	    {
		    "name": "debugger",
		    "versions": [
				{
					"version": 0,
					"fingerprint": "23b0b26f3c48a675b0fe3ab999fc0cfa3950ee051649c2bf45b7882976641eac",
					"size": 238,
					"created_at": 1547564337
				}
		    ],
		    "used_by": null
	    }
	]

