#!/bin/bash

# Enable service
systemctl enable argonone.service

# Add user
id -u gopi &>/dev/null || useradd --system -G i2c,video gopi

# Add directories and permissions
install -o gopi -g gopi -m 750 -d /opt/gopi/etc
install -o gopi -g gopi -m 755 -d /opt/gopi/var

# Start service
systemctl start argonone.service
