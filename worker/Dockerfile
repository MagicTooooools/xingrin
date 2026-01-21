# Worker Dockerfile - Go 版本
# 多阶段构建：编译扫描工具 + 构建 Go Worker

# ============================================
# 阶段 1: 编译 Go 扫描工具
# ============================================
FROM golang:1.24 AS tools-builder

ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1

# 安装编译依赖（构建阶段，不固定版本以获取最新安全补丁）
# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpcap-dev \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# 编译 massdns（puredns 依赖）
RUN git clone https://github.com/blechschmidt/massdns.git /tmp/massdns && \
    cd /tmp/massdns && \
    make && \
    cp bin/massdns /usr/local/bin/massdns

# 安装 ProjectDiscovery 等 Go 扫描工具
RUN go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest && \
    go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest && \
    go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest && \
    go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest && \
    go install -v github.com/projectdiscovery/katana/cmd/katana@latest && \
    go install -v github.com/tomnomnom/assetfinder@latest && \
    go install -v github.com/ffuf/ffuf/v2@latest && \
    go install -v github.com/d3mondev/puredns/v2@latest && \
    go install -v github.com/yyhuni/xingfinger@latest && \
    go install -v github.com/hahwul/dalfox/v2@latest

# ============================================
# 阶段 2: 构建 Go Worker
# ============================================
FROM golang:1.24 AS worker-builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 先复制依赖文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并构建
COPY . .
RUN CGO_ENABLED=0 go build -o /worker ./cmd/worker

# ============================================
# 阶段 3: 运行时镜像
# ============================================
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# 安装运行时依赖 + Python 环境
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libpcap-dev \
    nmap \
    masscan \
    curl \
    git \
    python3 \
    python3-pip \
    python3-venv \
    pipx \
    && rm -rf /var/lib/apt/lists/*

# 建立 python 软链接
RUN ln -s /usr/bin/python3 /usr/bin/python

# 目录约定：
# - /opt/orbit          仅存放数据（会被宿主机挂载覆盖）
# - /opt/orbit-tools    仅存放工具（来自镜像，不应被挂载覆盖）
RUN mkdir -p \
    /opt/orbit/wordlists \
    /opt/orbit/workspace \
    /opt/orbit-tools/bin \
    /opt/orbit-tools/share \
    /opt/orbit-tools/pipx

# 复制 resolvers.txt（DNS 服务器列表，内容固定）
COPY resources/resolvers.txt /opt/orbit/wordlists/resolvers.txt

# pipx：将 Python 工具安装到项目私有工具目录，避免落到 /root/.local/bin
ENV PIPX_HOME=/opt/orbit-tools/pipx
ENV PIPX_BIN_DIR=/opt/orbit-tools/bin
ENV PATH=/opt/orbit-tools/bin:$PATH

# 安装 Python 扫描工具
RUN pipx install uro && \
    pipx install waymore && \
    pipx install dnsgen

# 安装 Sublist3r（目录型工具，放到 /opt/orbit-tools/share）
RUN git clone https://github.com/aboul3la/Sublist3r.git /opt/orbit-tools/share/Sublist3r && \
    pip3 install --no-cache-dir -r /opt/orbit-tools/share/Sublist3r/requirements.txt --break-system-packages

# 从 tools-builder 复制 Go 扫描工具（统一放到 /opt/orbit-tools/bin）
COPY --from=tools-builder /go/bin/* /opt/orbit-tools/bin/
COPY --from=tools-builder /usr/local/bin/massdns /opt/orbit-tools/bin/massdns

# 从 worker-builder 复制 Go Worker
COPY --from=worker-builder /worker /usr/local/bin/worker

# 默认以 root 运行（最省心）：
# - 避免宿主机挂载目录权限不一致导致的写入失败
# - 扫描工具（如 masscan）在某些场景下需要更高权限
# 如需最小权限运行，可改为非 root 用户并按需添加 capability。
CMD ["worker"]
