### 目的:

让集体订餐变得简单

### 流程:

在订餐平台（饿了么/美团/百度）嵌入订餐按钮->进入本地菜单管理->锁定菜单 & 自动分单->由管理员开始订餐->通过扫二维码付款

### 目录暂时定:

    config.yaml  # 主配置文件，前后端共用
    main.go      # 后端程序入口
    backend/     # 后端目录入口
    template/    # 静态模板路经
	public/      # 生成的静态html文件路经