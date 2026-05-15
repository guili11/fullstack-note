1️⃣ 基础概念

Docker Compose：Docker 官方的多容器编排工具，用一个 YAML 文件管理多个容器，统一启动、停止、网络和卷挂载。

核心理念：把原本多个 docker run 命令整合成一个可维护、可重现的配置文件。

核心关键词：

关键词    说明

services    服务定义，每个容器一个 service
image    使用的镜像
build    使用 Dockerfile 构建镜像
ports    宿主机端口映射到容器端口
volumes    容器与宿主机或其他卷共享文件/数据
depends_on    启动顺序依赖
networks    容器网络，服务互联

---

2️⃣ 额外功能

环境变量（env）：

可以在 environment 或 .env 文件中定义，方便配置数据库密码、API URL 等。

environment:

- DB_HOST=db
- DB_USER=root
- DB_PASS=123456

Docker 网络：

默认 Compose 会创建一个隔离网络，容器之间可以用服务名互相访问。

networks:
  mynet:
    driver: bridge
具体访问：例如nginx代理后端容器，其实就是域名替换为容器名了：

```
       location /api/user {
          proxy_pass http://backend:8080;
          proxy_set_header Host $host;
          proxy_set_header X-Real-IP $remote_addr;
        }
```

卷（volumes）：

持久化数据或挂载代码目录，用于开发热更新或数据库持久化。

volumes:

- ./frontend:/app  # 挂载本地前端代码
- db_data:/var/lib/mysql  # 数据库持久化

依赖（depends_on）：

控制容器启动顺序。

depends_on:

- db

重启策略（可选）：

容器异常退出后自动重启。

restart: always

---

3️⃣ 文件示例

```yaml
version: "3.9"

services:
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
    depends_on:
      - backend

  backend:
    build: ./backend
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
    environment:
      - DB_HOST=db
      - DB_USER=root
      - DB_PASS=123456
    depends_on:
      - db

  db:
    image: mysql:8
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: 123456
    volumes:
      - db_data:/var/lib/mysql

volumes:
  db_data:

networks:
  default:
    driver: bridge
```

说明：

frontend → 前端开发服务（可挂载代码热更新）

backend → 后端服务，依赖数据库

db → MySQL 数据库，持久化数据

volumes → 数据卷

networks → 默认桥接网络，服务间可用服务名访问

---

💡 记忆点（面试快速答法）：

> “Docker Compose 用 YAML 管理多容器：
> services 定义容器，
> volumes 挂载代码或数据，
> networks 容器互联，
> depends_on 控制启动顺序，
> environment 配置变量。
> 开发环境可挂载代码热更新，生产环境用镜像和卷持久化。”
