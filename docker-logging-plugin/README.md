# docker-logging-plugin

The goal of this project is create a custom logger driver for docker.

[Docker Engine managed plugin system](https://docs.docker.com/engine/extend/)

[Writing a Docker Log Driver plugin — Basics](https://software-factotum.medium.com/writing-a-docker-log-driver-plugin-7275d99d07be)

[Writing a Docker Log Driver plugin — Part II](https://software-factotum.medium.com/writing-a-docker-log-driver-plugin-f94cee827a0a)

## Docker log plugins

[docker-log-logstash](https://github.com/rchicoli/docker-log-logstash)

[docker-plugin-multilogger](https://github.com/ekristen/docker-plugin-multilogger)

[docker-multilogger-plugin](https://github.com/allgdante/docker-multilogger-plugin)

[splunk docker-logging-plugin ](https://github.com/splunk/docker-logging-plugin)

## Issue

With docker 20.10 during enabling a plugin we will get error:
> Error response from daemon: dial unix /run/docker/plugins/c07c04bd10b78bbd4e1c6b975f75a794f06c5428838791eb0bf67a827672d295/customlogdriver.sock: connect: no such file or directory

There is a open issue in docker/go-plugins-helpers - [Plugin enable command is looking for .sock file in /run/docker/plugins/<docker-plugin-id> and not /run/docker/plugins #113](https://github.com/docker/go-plugins-helpers/issues/113)

