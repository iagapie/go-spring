# Spring CMS

### Run CMS
```shell
make run
```

http://localhost:8000

### Build plugin
```shell
make build_plugin_demo
```

### Build app
```shell
make build
```

### Generate asymmetric RSA for JWT
```shell
make cert
```

### Create backend user
```shell
./go-spring user:create -n Name -e name@gmail.com -p "Admin123"
```
