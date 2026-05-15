# 一,ngingx基础概念理解

前言：无论是前端还是后端，都是部署在服务器上，等待用户去访问url，然后从服务器拿东西，然后在用户端显示的（浏览器）

***

### 1️⃣ Nginx 到底是什么？

Nginx 是一个高性能的 HTTP & 反向代理服务器

它可以做 **Web 服务器**（直接提供静态资源 HTML/JS/CSS/图片等）

它也可以做 **反向代理服务器**（前端请求**只发到nginx上**，后端服务器不对外暴露，只由nginx可以访问，nginx决定前端请求最终访问到哪个后端--这也是nginx屌的根本）

还支持 负载均衡、缓存、限流、SSL/TLS 终端等功能

核心特性：

高并发：异步、事件驱动模型（处理大量并发连接而不占用太多资源）

轻量、快速、稳定

模块化：HTTP 模块、邮件代理模块、流式模块等

> 简单理解：Nginx 是前端和后端之间的“交通指挥官”，它可以决定请求去哪里，谁来响应，以及如何快速返回给客户端。

***

### 2️⃣ Nginx 解决了什么问题？

问题 A：前端请求直接到后端，压力大

如果所有客户端都直接访问后端：

每个请求都占用一个线程 / 进程

**高并发下容易挂掉**

Nginx 解决方案：

使用异步事件驱动，处理大量并发连接

作为反向代理，把请求分发到多台后端服务器（**负载均衡**）

***

问题 B：静态资源请求效率低

每次访问都让后端生成 HTML/JS/CSS → 低效率

Nginx 解决方案：

可以直接**托管静态文件**

响应速度快，**减少后端压力**

***

问题 C：**安全和隐藏后端**

直接暴露后端服务地址：

可能泄露 IP

安全防护能力弱

Nginx 解决方案：

客户端只看到 Nginx 的地址

可以做 HTTPS 终端、限流、防火墙等

后端地址隐藏，安全性提高

***

问题 D：负载均衡和高可用

单台后端挂掉 → 整个服务不可用

高流量下单台后端处理不过来

Nginx 解决方案：

配置多个后端节点

轮询 / 权重 / 最少连接等策略分发请求

自动健康检查，挂掉节点不再分流请求

***

### 3 核心理解

> Nginx = 高性能的前端入口 + 请求调度器 + 静态资源服务器 + 安全守护者

总结：

- HTML → 静态文件（Nginx 容器/服务器）
- 动态数据 → 后端 API（Nginx 代理）

它解决的核心问题：

1. 性能问题（高并发处理，快速响应）
2. 负载问题（分发请求，保护后端）
3. 静态资源问题（减轻后端压力，提升访问速度）
4. 安全问题（隐藏真实服务，SSL/TLS、防护）

***

问题 5：Docker 场景

- 在 Docker Compose 环境下，前端、Nginx、后端分别部署在不同容器，Nginx 的作用是什么？
- 如果前端是 SPA（单页应用），Nginx 主要负责什么？后端主要负责什么？

<br />

# 二，nginx的使用 **（安装+配置 +启动）**

> 这里暂时先介绍 docker环境下的使用    但是核心其实就是 **安装+配置 +启动**

### 1. 安装

这里使用的是docker，其实就是直接使用 nginx镜像了，有两种使用方法：

- 直接拉取官方的nginx镜像，容器启动的时候挂载前端文件到volume，适合静态资源托管（方便更新资源），开发环境还是用vite proxy （二者基本一致）
- 以nginx为基础镜像，使用前端文件 搭建 镜像，直接开启容器，不需要挂载， 适合 生产环境

挂载举例：  **实际上，这里已经是一个http服务器了**

docker run -d -p 80:80 \
-v /path/to/html:/usr/share/nginx/html \
\--name mynginx nginx

注意：

挂载 Volume 给 Nginx，只能挂载 磁盘上的静态文件，开发模式下内存里的文件 Nginx 是看不到的。

这里只可以挂载硬盘中的静态文件，html啊，图片啊，等。

vue这种是不可以npm run dev的， 除非你打包好了再进来调试，不过这就太麻烦了。

镜像举例：dist+nginx.conf

```
# 构建阶段，仅编译，用完就丢
FROM node:18-alpine AS builder
WORKDIR /app

# 1. 复制依赖文件，利用Docker缓存
COPY package*.json ./
RUN npm install

# 2. 复制全部源码
COPY . .

# 3. 构建生产环境产物，生成dist目录
RUN npm run build

# 运行阶段，极简环境，只有产物和必要的运行时
FROM nginx:alpine

# 1. 把构建好的静态文件复制到Nginx默认发布目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 2. 复制自定义的Nginx配置文件，替换默认配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 3. 暴露Nginx的80端口（容器内部端口）
EXPOSE 80

# 4. 前台运行Nginx，防止容器启动后直接退出
CMD ["nginx", "-g", "daemon off;"]
```

1.也是多阶段构造嘛，先打包环境 然后切换到运行环境

2.复制文件注意事项（多阶段必须注意）：

打包好的文件dist到 发布页面

nginx配置文件 xxx.conf替换默认配置（很重要，后续介绍配置文件的作用）

3.保证前台运行nginx，避免自动退出

<br />

### 2.配置

前言：前面提到了，nginx那么多的作用，那么作用肯定是要配置项去管理的呀，配置文件就是为了管理作用，例如请求具体发到哪个服务器上面，

<br />

其实这个东西是一个共性的，我们使用那么多的软件啊，组件啊，有很多功能，功能都是要靠配置文件去配置的，所以配置文件的时候我们关心需要什么功能就可以了

<br />

明确功能需求
问自己：我要软件实现哪些功能？

示例：

- Nginx：静态资源托管、SPA 路由、API 代理、缓存优化
- Docker：端口映射、卷挂载、环境变量

功能映射到配置指令

- **静态资源** → root, try\_files
- **API 代理** → proxy\_pass
- 缓存控制 → expires, add\_header
- 错误处理 → error\_page
- 日志/监控 → access\_log, log\_format

```
# nginx.conf
worker_processes  1;

events {
    worker_connections 1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    keepalive_timeout  65;

    server {
        listen       80;
        server_name  localhost;

        # 静态资源根目录
        root /usr/share/nginx/html;
        index index.html index.htm;

        # SPA 前端路由处理
        location / {
            try_files $uri $uri/ /index.html;
        }

        # API 代理示例（可选，开发用）
        # location /api/user {
        #     proxy_pass http://backend:8080;
        #     proxy_set_header Host $host;
        #     proxy_set_header X-Real-IP $remote_addr;
        # }
        # location /api/words {
        #     proxy_pass http://backend:8081;
        #     proxy_set_header Host $host;
        #     proxy_set_header X-Real-IP $remote_addr;
        # }

        # 静态资源缓存（可选）
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 30d;
            add_header Cache-Control "public";
        }

        # 错误页（可选）
        error_page 404 /index.html;
    }
}
```

三点核心：

静态资源根目录设置

spa路由处理

原理问法

“为什么 Vue/React 项目刷新页面会 404？”

答：浏览器直接请求服务器路径，但服务器默认找不到对应文件。

解决方案问法
“如何解决？”

答：Nginx try\_files $uri $uri/ /index.html;，把未匹配路径都返回首页，由前端路由处理

api代理

### nginx启动

这一点其实在安装里详细写了

<br />

### 使用策略：

- 开发环境：使用 Vite 的开发服务器 (npm run dev)，通过 Vite 的 proxy 解决跨域问题，支持热更新，方便快速开发和调试。

- 生产环境：前端代码打包生成 dist，通过自定义 Nginx 镜像部署，Nginx 挂载静态资源（HTML、JS、CSS、图片等），处理生产环境下的跨域和缓存优化，保证环境稳定可靠。”
  
  <br />
