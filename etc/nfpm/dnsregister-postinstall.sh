#!/bin/bash

# Enable service
deb-systemd-helper enable dnsregister.service

# Add user
id -u gopi &>/dev/null || useradd --system -G i2c,video gopi

# Add /var dir
install -o gopi -m 755 -d /opt/gopi/var

# Start service
deb-systemd-invoke start dnsregister.service
