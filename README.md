pentaconta
==========

A service runner

config
------

Searches for a configfile in working dir and etc named pentaconta.yaml or pentaconta.json
The config can be changed by giving a parameter `pentaconta -config custom/main`. This would look
in `workingdir/custom/main.json` and `/etc/custom/main.json` as well as the yaml variants. Also
an absolute path can be given.

    services:
      pc_stable_service:
        executable: /usr/bin/pc_stable
        arguments: [--foo, bar]
    
    fs-triggers:
      pc_stable_bin_trigger:
        path: /usr/bin/pc_stable
        services: [pc_stable_service]

      pc_stable_conf_trigger:
        path: /etc/pc_stable.yaml
        services: [pc_stable_service]

This would define one service. The given executable is started and if it crashes ot terminates regularly
it will be restarted after a pause of 1 second.

The two fs-triggers are responsible that the service also gets restarted when the executable itself gets
updated or when the config file gets updated.

logs
----

All the output from stdout and er are sent to stdout as well as some logging info from pentaconta itself:

    2017-04-21 23:01 49.012619 OUT pc_stable: Stable main started
    2017-04-21 23:01 50.012928 OUT pc_stable: I'm doing fine
    2017-04-21 23:01 51.013093 OUT pc_stable: I'm doing fine
    2017-04-21 23:02 01.015132 OUT pc_stable: I'm doing fine
    2017-04-21 23:02 01.310071 PEN pc_stable: Terminated service with signal: interrupt <--- service gets sigint from config file change
    2017-04-21 23:02 01.610011 PEN pc_stable: Sigint worked
    2017-04-21 23:02 02.313112 OUT pc_stable: Stable main started
    2017-04-21 23:02 03.313456 OUT pc_stable: I'm doing fine
