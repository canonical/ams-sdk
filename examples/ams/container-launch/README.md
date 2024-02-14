Launch Container Example
========================

Demonstrates how to launch a new container using AMS SDK

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server              | required |
| `id`      | Identifier of the application      | required |
| `version` | Version of the application         | optional |
| `node`    | In which node to launch the container | optional |
| `raw`     | Launched a container from a specific image instead of an application if it's set to true | optional | 
| `instance-type` | Instance type to use when launching a container from an image instead of an application | optional |

*NOTE*: If the optional parameter `raw` is set to true, this parameter `id` will accept an image id 
or image name instead. As a result, a container will be launched from a specific image instead of
an application


Example:

    container-launch -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 id=bgutrvm5nof0fqm0894g

Output:

    {
		"containers": [
			"bgv0afe5nof0fqm089b0"
		]
	}


