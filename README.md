# 临时

达梦数据库驱动
在[https://github.com/gogf/gf/tree/master/contrib/drivers/dm/v2](https://github.com/gogf/gf/tree/master/contrib/drivers/dm/v2)基础上修改

主要修改内容

1、执行gen dao时会读取所有表 修改dm_tables.go文件查询表的SQL语句

2、执行gen dao时无法获取表字段信息（主键、自增、备注） 修改dm_table_fields.go文件查询SQL语句
生成md格式
