#!/bin/bash

# Todo List 项目一键启动脚本
# 同时启动后端 (Go) 和前端 (Vite) 服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"

# PID 文件
BACKEND_PID_FILE="$PROJECT_DIR/.backend.pid"
FRONTEND_PID_FILE="$PROJECT_DIR/.frontend.pid"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 清理函数
cleanup() {
    log_info "正在停止服务..."

    # 停止后端
    if [ -f "$BACKEND_PID_FILE" ]; then
        BACKEND_PID=$(cat "$BACKEND_PID_FILE")
        if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
            log_info "停止后端服务 (PID: $BACKEND_PID)"
            kill "$BACKEND_PID" 2>/dev/null || true
        fi
        rm -f "$BACKEND_PID_FILE"
    fi

    # 停止前端
    if [ -f "$FRONTEND_PID_FILE" ]; then
        FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
        if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
            log_info "停止前端服务 (PID: $FRONTEND_PID)"
            kill "$FRONTEND_PID" 2>/dev/null || true
        fi
        rm -f "$FRONTEND_PID_FILE"
    fi

    # 额外清理：杀死可能遗留的进程
    pkill -f "go run cmd/server/main.go" 2>/dev/null || true
    pkill -f "vite" 2>/dev/null || true

    log_success "所有服务已停止"
    exit 0
}

# 捕获退出信号
trap cleanup SIGINT SIGTERM

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    # 检查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装，请先安装 Go: https://golang.org/dl/"
        exit 1
    fi

    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装，请先安装 Node.js: https://nodejs.org/"
        exit 1
    fi

    # 检查 npm
    if ! command -v npm &> /dev/null; then
        log_error "npm 未安装"
        exit 1
    fi

    log_success "依赖检查完成"
}

# 启动后端
start_backend() {
    log_info "启动后端服务..."

    cd "$BACKEND_DIR"

    # 检查 go.mod
    if [ ! -f "go.mod" ]; then
        log_error "后端 go.mod 不存在"
        exit 1
    fi

    # 安装依赖（如果需要）
    if [ ! -d "vendor" ] && [ -f "go.mod" ]; then
        log_info "安装 Go 依赖..."
        go mod download
    fi

    # 启动后端服务
    nohup go run cmd/server/main.go > "$PROJECT_DIR/backend.log" 2>&1 &
    BACKEND_PID=$!
    echo $BACKEND_PID > "$BACKEND_PID_FILE"

    # 等待后端启动
    sleep 2

    if ps -p $BACKEND_PID > /dev/null; then
        log_success "后端服务已启动 (PID: $BACKEND_PID)"
        log_info "后端地址: ${GREEN}http://localhost:8080${NC}"
    else
        log_error "后端服务启动失败，请查看日志: $PROJECT_DIR/backend.log"
        cat "$PROJECT_DIR/backend.log"
        exit 1
    fi
}

# 启动前端
start_frontend() {
    log_info "启动前端服务..."

    cd "$FRONTEND_DIR"

    # 检查 package.json
    if [ ! -f "package.json" ]; then
        log_error "前端 package.json 不存在"
        exit 1
    fi

    # 检查 node_modules
    if [ ! -d "node_modules" ]; then
        log_info "安装前端依赖..."
        npm install
    fi

    # 启动前端服务
    log_info "正在启动 Vite 开发服务器..."
    nohup npm run dev > "$PROJECT_DIR/frontend.log" 2>&1 &
    FRONTEND_PID=$!
    echo $FRONTEND_PID > "$FRONTEND_PID_FILE"

    # 等待前端启动
    sleep 3

    if ps -p $FRONTEND_PID > /dev/null; then
        log_success "前端服务已启动 (PID: $FRONTEND_PID)"
        log_info "前端地址: ${GREEN}http://localhost:5173${NC}"
    else
        log_error "前端服务启动失败，请查看日志: $PROJECT_DIR/frontend.log"
        cat "$PROJECT_DIR/frontend.log"
        cleanup
        exit 1
    fi
}

# 显示启动信息
show_banner() {
    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║           Todo List 项目 - 一键启动                  ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# 主函数
main() {
    show_banner

    # 检查是否有遗留的 PID 文件
    if [ -f "$BACKEND_PID_FILE" ] || [ -f "$FRONTEND_PID_FILE" ]; then
        log_warn "检测到遗留的 PID 文件，正在清理..."
        cleanup
        sleep 1
    fi

    # 检查依赖
    check_dependencies

    # 启动服务
    start_backend
    start_frontend

    echo ""
    log_success "🚀 所有服务启动完成！"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "  ${GREEN}后端 API:${NC}    http://localhost:8080/api"
    echo -e "  ${GREEN}前端应用:${NC}    http://localhost:5173"
    echo -e "  ${GREEN}后端日志:${NC}    $PROJECT_DIR/backend.log"
    echo -e "  ${GREEN}前端日志:${NC}    $PROJECT_DIR/frontend.log"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    log_info "按 ${YELLOW}Ctrl+C${NC} 停止所有服务"
    echo ""

    # 保持脚本运行
    while true; do
        sleep 1

        # 检查服务是否还在运行
        BACKEND_RUNNING=false
        FRONTEND_RUNNING=false

        if [ -f "$BACKEND_PID_FILE" ]; then
            BACKEND_PID=$(cat "$BACKEND_PID_FILE")
            if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
                BACKEND_RUNNING=true
            fi
        fi

        if [ -f "$FRONTEND_PID_FILE" ]; then
            FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
            if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
                FRONTEND_RUNNING=true
            fi
        fi

        # 如果服务停止了，自动退出
        if ! $BACKEND_RUNNING || ! $FRONTEND_RUNNING; then
            log_warn "检测到服务已停止"
            cleanup
        fi
    done
}

# 运行主函数
main
