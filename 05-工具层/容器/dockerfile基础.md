Dockerfile 描述构建步骤；在包含该文件的目录执行 docker build 生成镜像。

功能：主要就是等价于 我们在控制台里手动执行的命令。
规定好：用什么环境、复制哪些代码、装依赖、开什么端口、最后运行你的程序。

注意：复制，前是宿主路径，后是容器路径。

指令	作用
FROM	基础镜像，必须是第一条有效指令（除解析器指令外）。
WORKDIR	设置工作目录，后续 RUN/CMD 的默认路径。
COPY	从构建上下文复制文件到镜像（推荐优先于 ADD）。
RUN	  构建时执行命令（安装依赖、编译等）。
ENV	  设置环境变量。
EXPOSE	声明容器监听端口（文档性质，不自动发布到宿主机）。
CMD	容器默认启动命令，可被 docker run 末尾参数覆盖。
ENTRYPOINT	入口；与 CMD 组合时常用作固定可执行文件，CMD 只传参数。


### 构建与运行
```bash
# 构建镜像
docker build -t myapp:v1 .

# 运行容器
docker run -d -p 8080:8080 --name myapp myapp:v1
```

### 记忆口诀
**"FROM定基础 → WORKDIR定目录 → COPY搬文件 → RUN装依赖 → EXPOSE开端口 → CMD定启动"**

### 关键注意点
- 每条指令都会创建一个新的镜像层（层数越少越好）
- `COPY` 和 `ADD` 的区别：`ADD` 支持URL和解压，一般用 `COPY` 更安全
- `CMD` vs `ENTRYPOINT`：`CMD` 可被覆盖，`ENTRYPOINT` 是固定入口
- 利用缓存：把不常变的放前面（如依赖），常变的放后面（如源码）

### 镜像分层机制

#### 详细原理

**1. 镜像的本质**

Docker镜像不是"一个大文件"，而是**多层只读文件系统的叠加**。每条Dockerfile指令都会对文件系统进行一次修改，从而生成一个新的镜像层。

```
┌─────────────────────┐
│  层4: CMD ["./app"]  │  ← 元数据层
├─────────────────────┤
│  层3: COPY . .       │  ← 文件变化层
├─────────────────────┤
│  层2: RUN npm install│  ← 文件变化层
├─────────────────────┤
│  层1: FROM node:20   │  ← 基础镜像层
└─────────────────────┘
```

**2. 为什么要分层？（三大核心原因）**

**① 缓存复用（最重要）**

```dockerfile
FROM node:20          # 层1：基础镜像
WORKDIR /app          # 层2：创建目录
COPY package.json .   # 层3：复制依赖文件
RUN npm install       # 层4：安装依赖（耗时！）
COPY . .              # 层5：复制源码
```

- 你只改了业务代码，没改 `package.json`
- Docker发现层1-4没变化，直接复用缓存
- 只重新构建层5，构建从2分钟 → 5秒

**② 存储共享（省空间）**

```
容器A: nginx + 前端代码
容器B: nginx + 后端代码
容器C: nginx + 静态资源
```

- 三个容器都用了 `FROM nginx`
- nginx基础镜像只存一份，三个容器共享同一层
- 如果合并成一层，就要存三份完整的nginx

**③ 版本控制 & 回滚**

- 每层都有唯一的 **SHA256 ID**
- 可以追踪：哪层改了什么
- 出问题可以回滚到某一层

**3. 实际验证（查看镜像分层）**

```bash
# 构建镜像
docker build -t myapp .

# 查看镜像分层
docker history myapp
```

输出示例：
```
IMAGE          CREATED        SIZE
abc123         2 min ago      0B     CMD ["./app"]
def456         2 min ago      50MB   COPY . .
ghi789         3 min ago      120MB  RUN npm install
jkl012         5 min ago      0B     COPY package.json .
mno345         5 min ago      0B     WORKDIR /app
pqr678         1 hour ago     150MB  FROM node:20
```

**4. 层数太多的问题**

虽然分层好，但也不是越多越好：

```dockerfile
# ❌ 不好的写法（层数多）
RUN apt update
RUN apt install -y nginx
RUN apt install -y curl
RUN apt install -y vim

# ✅ 好的写法（合并层）
RUN apt update && \
    apt install -y nginx curl vim && \
    rm -rf /var/lib/apt/lists/*
```

**原因**：
- 层数多 → 镜像体积大（每层都保留中间文件）
- 层数多 → 构建慢
- 层数多 → 网络传输慢

---

#### 总结

> **Docker分层 = Git的commit机制**
> - 每次提交（指令）生成一个新快照（层）
> - 可以复用、共享、回滚
> - 但commit太多也不好，要合理合并

**面试常考点**：

**Q: 为什么Dockerfile要把 `COPY package.json` 和 `RUN npm install` 分开写，而不是直接 `COPY . .` 然后 `RUN npm install`？**

**A**: 利用分层缓存机制。`package.json` 不常变，单独复制后，Docker可以缓存 `npm install` 的结果。后续只改业务代码时，跳过耗时的依赖安装，大幅提升构建速度。



