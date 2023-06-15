# Grafiport
This is a Tool to export Grafana objects and store them in JSON Files the filesystem.
In addition to this it can restore the config from the Files and import them back into Grafana.

## Usage/Installation

You can install Grafiport either as a container (preferred method) or standalone.

### Docker

Just start a container with following command.
You have to set username, password, url as environment variables.
The default exporting/restoring path inside the container is /output.
Use Volume Mounts to export/restore to/from your file system.
```
docker container -e username=<username> -e password=<password> -e url=<url> --name grafiport grafiport/grafiport
```

Following vars can be set via environment variables:
Note: Don't the directory variable in the container environment, unless you know what you are doing. 

| Variable  | Flag   | Description                                             |
|-----------|--------|---------------------------------------------------------|
| username  | String | Username of Grafana instance                            |
| password  | String | Password of Grafana instance                            |
| url       | String | Url of the Grafana Instance                             |
| directory | String | Target or origin folder of Grafana Config Objects       |
| restore   | Bool   | Activate restore mode. Default: false                   |
| alerting  | Bool   | If restoring/exporting config objects consider Alerting |


### non Container environment
Download the executable for your systems OS/architecture in the Release and just execute it. No further installation needed
In non container environments the same variables can be set to configure Grafiport.
You have to set them preferribly as Flags instead of environment variables.
But you can set settings as environment variables as well. 
```

./grafiport -username <username> -password <password> -url <url>

```
If you need more help just use the -h flag to display the Usage

## Contribute

Feel free to contribute! Use it, test it, give feedback or even Document and code.
Help is always welcome.

## Useful Information

### Restore of Alerting Config Objects not recommended

Few Objects in Grafana are not supposed to be configured via api. They result would be an object
which you can't configure anymore [explanation][provisioned_ressources].
More to come...

[provisioned_ressources]: https://grafana.com/docs/grafana/latest/alerting/set-up/provision-alerting-resources/view-provisioned-resources
