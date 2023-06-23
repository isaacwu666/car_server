create table player
(
    id        int primary key auto_increment comment 'id',
    phone     varchar(32) comment '电话',
    pwd       varchar(64) comment '密码',
    salt      varchar(12) comment 'pwd盐',
    nick_name varchar(32) comment '昵称'

)comment '玩家信息表'
#     gf gen dao -l "mysql:root1234:root@tcp(127.0.0.1:3306)/happy_poker"