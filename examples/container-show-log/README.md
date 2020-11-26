Show Container Log Example
========================

Demonstrates how to show the specified log file of a container

Parameters
-----

You have to provide the following parameters in any order:

| Name      | Description           | Attribute  |
| --------- |:--------------------  | :--------: |
| `cert`    | Path to the file with the client certificate to use to connect to AMS | required |
| `key`     | Path to the file with the client key to use to connect to AMS  | required |
| `url`     | URL of the AMS server              | required |
| `id`      | Identifier of the container | required |
| `log-name`| Name of the container log   | required |


Example:

    container-log -cert=./client.crt -key=./client.key -url=https://<ams_ip_address>:8443 -id=bgvb70vb9s9jdk3dpujg -log-name=container.log

Output:
    ....
    ....
	lxc 20190114075135.363 DEBUG    lxc_start - start.c:__lxc_start:1402 - Tearing down virtual network devices used by container "default".
	lxc 20190114075135.363 INFO     lxc_conf - conf.c:lxc_delete_network:3072 - Interface "(null)" with index 114 already deleted or existing in different network namespace.
	lxc 20190114075135.423 INFO     lxc_conf - conf.c:lxc_delete_network:3105 - Removed interface "vethR8S81A" from host.
	lxc 20190114075135.463 INFO     lxc_conf - conf.c:run_script_argv:435 - Executing script "/usr/share/lxcfs/lxc.reboot.hook" for container "default", config section "lxc".
