<!--
title: "Suricata module"
description: "The Suricata modules retrieves IDS metrics from your suricata node instance. The collected metrics are forwarded to your netdata agent."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/suricata/README.md
sidebar_label: "Suricata module"
-->

# Suricata module

The Suricata modules retrieves IDS metrics from your suricata node instance. The collected metrics are forwarded to your netdata agent.

## Charts

This module produces Intrusion Detection System charts with security metrics collected from your Suricata node instance.

The collected metrics and charts are,
1. Total Events per sec
2. Intrusion detection events per second
3. ...

## Configuration

Edit the `go.d/suricata.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/suricata.conf
```

Disabled by default. Should be explicitly enabled
in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

```yaml
# go.d.conf
modules:
  suricata: yes
```

Here is an suricata configuration with several jobs:

```yaml
jobs:
  - name: suricata
    charts:
      num: 3
      dimensions: 5

  - name: hidden_suricata
    hidden_charts:
      num: 3
      dimensions: 5
```

---

For all available options, see the suricata
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/suricata.conf).

## Troubleshooting

To troubleshoot issues with the `suricata` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m suricata
```


# Dependencies for suricata module
#### Suricata installation
1. [Suricata installation in your node ](https://suricata.readthedocs.io/en/latest/install.html#install-advanced).
2. [Suricata necessary libraries](https://suricata.readthedocs.io/en/latest/install.html#dependencies)
3. Enable unix socket in /etc/suricata/suricata.yaml configuration file.
```
# Unix command socket that can be used to pass commands to Suricata.
# An external tool can then connect to get information from Suricata
# or trigger some modifications of the engine. Set enabled to yes
# to activate the feature. In auto mode, the feature will only be
# activated in live capture mode. You can use the filename variable to set
# the file name of the socket.
unix-command:
  enabled: yes
  filename: /var/run/suricata/suricata-command.socket
```

#### netdata agent
1. 
#### Enable go.d.plugin module