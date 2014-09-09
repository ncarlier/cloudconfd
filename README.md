cloudconfd
==========

A very simple cloud-config http server.

This can help when your cloud provider don't provide an easy way to setup user-data.

Installation
------------

Linux binaries for release [0.0.1](https://github.com/ncarlier/cloudconfd/releases)

* [amd64](https://github.com/ncarlier/cloudconfd/releases/download/v0.0.1/cloudconfd-linux-amd64-v0.0.1.tar.gz)

Download the version you need, untar, and install to your PATH.

    $ wget https://github.com/ncarlier/cloudconfd/releases/download/v0.0.1/cloudconfd-linux-amd64-v0.0.1.tar.gz
    $ tar xvzf cloudconfd-linux-amd64-v0.0.1.tar.gz
    $ ./cloudconfd

Create your own user-data template in the "templates" directory. The filename define the configuration name.

**Example:** "./templates/cloud-config.yaml" will define the "cloud-config" configuration. 

Check this project directory for a sample.

Create a configuration file in a directory named according the configuration name in the 'conf' directory. This file is a YAML file and must be named according the mac address of the node you want to configure. You have to replace ":" char by "_" char (for a better file system compatibility).

**Example:** ./conf/cloud-config/da_0f_f7_b2_a1_b3.yaml

Again, check this project directory for a sample.

And launch the server:

    cloudconfd -l :9090

Now you can fetch the wanted configuration like this:

    curl localhost:9090/cloud-config/da:0f:f7:b2:a1:b3

This will fetch the configuration named "cloud-config" for the MAC address da:0f:f7:b2:a1:b3.

Use case
--------

A simple use case of cloudconfd is to boot a CoreOS node with a minimalist configdrive configured to link cloud-init with cloudconfd server.

First, create a minimalist user-data setup:

    #!/bin/bash
    MAC=`ifconfig ens18 | grep -o -E '([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}'`
    URL="http://192.168.0.2:9090/cloud-init/${MAC}"
    coreos-cloudinit --from-url="${URL}"

Adapt the URL to your cloudconfd server location.

Then create an ISO file:

    mkdir -p /tmp/new-drive/openstack/latest
    cp user_data /tmp/new-drive/openstack/latest/user_data
    mkisofs -R -V config-2 -o configdrive.iso /tmp/new-drive
    rm -r /tmp/new-drive

Mount the iso alongside your CoreOs node and boot the node.

And voila. The CoreOs is set up thanks to cloudconfd.
