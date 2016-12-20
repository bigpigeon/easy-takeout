## 引言

easy-takeout的目的是让一个团队一些订餐变得更便利


### 如何安装

可以使用源码安装:

    go get github.com/easy-takeout/easy-takeout
	
或者直接下载二进制包:

    TODO

### 如何使用

确保有可用的redis服务器


根据[config.toml.temp](config.toml.temp)配置好config.toml


执行

    easy-takeout migrate -config config.toml
    easy-takeout generate -config config.toml
	easy-takeout server -config config.toml



更多关于easy-takeout的使用可以看[这里](//github.com/easy-takeout/easy-takeout/tree/master/backend/command)




