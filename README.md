# ad-infrastructure-api

Register Teams and Users API for autogenerate [ad-infrastructure](https://github.com/Ivanhahanov/ad-infrastructure)
configs.

## Deploy

### Docker

coming soon...

### Golang

Change config name `mv config.yml.example config.yml ` and edit it.

Config parameters:
* **users**
* **teams**

You can use one of parameters for create only players or only teams. If users and teams specified it means that infrastructure has teams with users. 

Resource parameters:
* **memory** - count of memory (default: 2048)
* **vcpu** - number of virtual cpu (default: 2)

Run Api use golang

```
go run main.go
```
