Nginx：一个官方配置好的docker镜像

反向代理：
反向代理 = 你把请求发给 Nginx，Nginx 再把请求发给你的服务。

同时它也是：静态网页服务器、负载均衡器。

Nginx 反向代理 3 大好处
隐藏真实后端端口
不用暴露 8080、9000 各种奇怪端口
统一入口
前端、Go 后端、图片静态资源，全部统一走 Nginx
负载均衡
你开 3 个 Go 容器服务，Nginx 自动平分流量，不会崩

目前最佳就是用docker
通过官方nginx镜像
快速开启nginx容器
docker run -d --name nginx-web --network my-net -p 80:80 -v 配置文件路径 nginx
name    容器名称
network  共享网络
v   挂载配置文件到容器内（自动同步）
nginx 使用 nginx 官方镜像快速启动容器
![alt text](image.png)



# 监听 80 端口（浏览器默认端口，不用输端口号）

server {
listen 80;
server\_name localhost;

```
# 所有请求，全部转发给 Go 容器
location / {
    # 核心！容器间通信：直接用 容器名:内部端口
    proxy_pass http://go-demo:8080;
    
    # 固定标配（不用改）
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

<br />

<br />

# 第一步：创建 Docker 共享网络（容器互通必须）

&#x20;

cmd 执行：

bash

运行

```
docker network create my-net

```

作用：让 Nginx 容器 和 Go 容器 住同一个局域网，互相访问。

***

# 第二步：启动你的 Go 容器

bash

运行

```
docker run -d \
--name go-app \
--network my-net \
my-go-demo

```

解释：

- `--name go-app` 给容器起名（后面 Nginx 要靠这个名字找它）
- 不映射端口也行，容器内部互通不需要宿主机端口

***

# 第三步：写 Nginx 配置（关键）

1. 在桌面新建文件夹：`nginx-conf`
2. 里面新建文件，名字：`default.conf`
3. 粘贴下面全部内容：

nginx

```
server {
    listen 80;
    server_name localhost;

    location / {
        # 核心：容器名:内部端口
        proxy_pass http://go-app:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

```

***

# 第四步：启动 Nginx 独立容器

复制下面命令，**修改你电脑的文件路径**：

bash

运行

```
docker run -d \
--name nginx-web \
--network my-net \
-p 80:80 \
-v C:\Users\Administrator\Desktop\nginx-conf\default.conf:/etc/nginx/conf.d/default.conf \
nginx

```

### 路径对照（你直接用）

你的配置文件在：

`C:\Users\Administrator\Desktop\nginx-conf\default.conf`

上面命令已经写好，直接运行就行。

***

# 第五步：直接测试

打开浏览器输入：

plaintext

```
http://localhost

```

✅ 直接看到：`hello-docker`

***

# 完整流程逻辑

1. 浏览器访问 `localhost:80`
2. 打到 **Nginx 独立容器**
3. Nginx 通过容器名 `go-app` 内网转发
4. 交给 **Go 独立容器 8080**
5. 两个容器完全独立，互不附属

