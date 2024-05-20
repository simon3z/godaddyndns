# Installing and Upgrading

Installing:

    $ GO111MODULE=on go get -ldflags="-s -w" github.com/simon3z/nsdyndns/cmd/nsdyndnsd

Cross-compiling for a Raspberry Pi:

    $ GO111MODULE=on GOARCH=arm GOARM=5 go get -ldflags="-s -w" github.com/simon3z/nsdyndns/cmd/nsdyndnsd

# Usage

    $ nsdyndnsd -h
    Usage of nsdyndnsd:
      -c string
            configuration file path
      -i int
            check interval in seconds (default 300)

# Installing as Systemd Service

Checkout the repository and build locally:

    # go build ./cmd/nsdyndnsd

Install the nsdyndnsd binary and configuration:

    # cp -vi nsdyndnsd /usr/local/sbin/nsdyndnsd
    # (umask 077 && cp -vi init/nsdyndnsd.conf /etc/nsdyndnsd.conf)

Edit the configuration file:

    # vim /etc/nsdyndnsd.conf

Install the nsdyndnsd systemd service:

    # cp -vi init/nsdyndnsd.service /etc/systemd/system/nsdyndnsd.service

Reload systemd services:

    # systemctl daemon-reload

Enable the nsdyndnsd systemd service:

    # systemctl start nsdyndnsd.service
    # systemctl enable nsdyndnsd.service

Check the nsdyndnsd service status:

    # systemctl status nsdyndnsd.service
