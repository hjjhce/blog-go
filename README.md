"# blog-go" 
whyspacex@aliyun.com

单页应用
REST API风格


docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 mysql:5.7
docker run --name phpmyadmin -d --link mysql:db -p 9002:80 phpmyadmin/phpmyadmin

后台功能
    1.用户管理
    2.文章管理
        新增
        修改
        删除
        查找
    3.评论管理
        新增
        删除
    4.分类管理
        文章
        视频
    5.设置
        首页设置
            title
            banner图
        爬虫设置
        导航栏设置
        底部设置
    6.爬虫模块
    7.图片算法
    8.性能剖析（pprof, new relic）
    