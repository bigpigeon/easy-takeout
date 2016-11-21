### 使用方式
  easy-takeout [flags]


  easy-takeout [commands] [flags]

### 可用命令简介

**generate**

使用模板解析生成html文件


执行该命令会将template/中的html文件渲染并写入public/


在template/static中的文件不会被解析，并且会原封不动的拷贝到到public中


可以通过 "."的方式调用Config中的属性 e.g. {{.BaseUrl}} => http://127.0.0.1


关于渲染规则可以看[render库](//github.com/easy-takeout/easy-takeout/tree/master/backend/render)

**migrate**

在数据库中创建/修改表格

**server**

一切准备就绪后，就可以开启服务

**print**

以toml格式打印当前配置


