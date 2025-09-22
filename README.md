# 临时

在goframe提供的达梦数据库驱动上修改

原驱动地址：[https://github.com/gogf/gf/tree/master/contrib/drivers/dm/v2](https://github.com/gogf/gf/tree/master/contrib/drivers/dm/v2)

修改说明：

在使用达梦V8数据库时 遇到一下问题

1、执行gen dao时会读取所有表（在读取字段信息时，需要对用户手DBA权限，导致生成dao时读取了所有表） 修改dm_tables.go文件查询表的SQL语句 增加owner条件

2、执行gen dao时无法获取表字段信息（是否主键、是否自增、备注、类型、长度等数据） 修改dm_table_fields.go文件查询SQL语句
