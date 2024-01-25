https://songxwn.com/Sublime_Text4_4169/

【MySQL】MySQL忘记密码或修改密码的方法
1.MySQL修改新密码方法
2. MySQL忘记密码，重置密码方法
1.MySQL修改新密码方法
记得原密码情况下，修改新密码：
登录到数据库后，输入 set password for 用户名@localhost = ‘新密码’; 来设置新的密码，别忘记分号哦。
如图所示：为本机localhost MySQL数据库系统中 root用户修改新密码为 admin

修改root@localhost用户的密码为admin: set password for root@localhost = 'admin';

在这里插入图片描述

2. MySQL忘记密码，重置密码方法
忘记登录密码情况下，通过以步骤行重置MySQL数据库系统的用户登录密码。

1.使用管理员身份打开cmd，确保关闭mysql服务，cmd输入命令: net stop mysql

ps:笔者安装的mysql版本是
在这里插入图片描述

Server version: 8.0.12，而我MySQL服务名称 为 MySQL80，所以我使用 net stop mysql80命令关闭mysql服务

在这里插入图片描述

2.将目录从默认c盘位置切换到mysqld.exe的安装目录（如我的目录：D:\Program Files\MySQL\MySQL Server 5.7\bin）

在这里插入图片描述

则在cmd黑窗口输入如下命令，切换到mysqld.exe的安装目录

（一般是 xxx\MySQL\MySQL Server 5.7\bin 目录下）

在这里插入图片描述

3.跳过密码验证

由于 mysqld --skip-grant-tables 命令实测在mysql8.0.12版本中已失效。

MySQL 8.0.x 版本推荐使用命令 mysqld --console --skip-grant-tables --shared-memory

低版本MySQL数据库，使用mysqld --skip-grant-tables

停止mysql服务后，输入mysqld --skip-grant-tables

或者如下图命令：mysqld -nt --skip-grant-tables

以上两条命令都可以：

在这里插入图片描述
执行到这里就只会有光标在一闪一闪无法继续写命令或输入任何命令了，故重新再打开一个cmd窗口

4.在新的cmd窗口中进行如下操作（这一步是否以管理员身份打开新cmd窗口都可以）

切换到mysqld.exe的安装目录，以无账号密码方式登录MySQL，然后重置数据库系统 root用户的密码为admin

flush privileges;

set password for root@localhost='root2023'

忘记密码情况下，重置密码完成！