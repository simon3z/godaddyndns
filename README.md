# Installing and Upgrading

Installing:

    $ GO111MODULE=on go get -ldflags="-s -w" github.com/simon3z/godaddyndns/cmd/godaddyndnsd

Cross-compiling for a Raspberry Pi:

    $ GO111MODULE=on GOARCH=arm GOARM=5 go get -ldflags="-s -w" github.com/simon3z/godaddyndns/cmd/godaddyndnsd

# Usage

    $ godaddyndnsd -h
    Usage of godaddyndnsd:
      -c string
            configuration file path
      -i int
            check interval in seconds (default 300)

# Installing as Systemd Service

Checkout the repository and build locally:

    # go build ./cmd/godaddyndnsd

Install the godaddyndnsd binary and configuration:

    # cp -vi godaddyndnsd /usr/local/sbin/godaddyndnsd
    # (umask 077 && cp -vi init/godaddyndnsd.conf /etc/godaddyndnsd.conf)

Edit the configuration file:

    # vim /etc/godaddyndnsd.conf

Install the godaddyndnsd systemd service:

    # cp -vi init/godaddyndnsd.service /etc/systemd/system/godaddyndnsd.service

Reload systemd services:

    # systemctl daemon-reload

Enable the godaddyndnsd systemd service:

    # systemctl start godaddyndnsd.service
    # systemctl enable godaddyndnsd.service

Check the godaddyndnsd service status:

    # systemctl status godaddyndnsd.service
