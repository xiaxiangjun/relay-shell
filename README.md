# relay-shell



利用中继服务器转发 SSH 连接，实现本地与远程主机之间基于 SSH 协议的通信。

## 工程说明

适用场景:

* 实现跨平台、跨网络的本机与远程主机之间的互通
* 支持嵌入式设备系统、Linux 系统、Windows 系统的客户端
* 在无法方便安装 ssh-server 的环境中使用
* 在远程主机无法访问本地网络环境的情况下使用
* 在本地与远程主机无法互通的环境中使用

实现以下功能：

- 利用中继服务器连接本地与远程主机
- 实现远程主机通过代理访问本地网络资源

适合以下用户:

* 开发人员：可通过 IDE 工具快速同步代码，并在远程主机上进行编译、运行、调试操作

## 使用方法

* 需要一台带有公网IP的主机（或者是本机与远程主机都能访问的主机）来部署 `relay-center`。
* 在远程主机上运行 `relay-shell` 客户端工具，获取唯一的机器码。
* 在本地使用 `relay-shell` 工具监听一个端口，并使用机器码和密码与远程主机建立连接。
* 通过其他 SSH 客户端工具访问本地监听的端口（无需输入密码）。



## 使用第三方库说明

| 第三方库                   | 使用方式 | 备注             |
| -------------------------- | -------- | ---------------- |
| github.com/google/gops     | 引用     | 监视运行状态     |
| github.com/quic-go/quic-go | 引用     | quic协议         |
| golang.org/x/net           | 引用     | 网络扩展         |
| gopkg.in/yaml.v3           | 引用     | 配置yaml文件解析 |

