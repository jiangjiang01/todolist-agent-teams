#!/bin/bash

# Todo List 项目一键停止脚本
# 停止后端和前端服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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

# 停止后端
stop_backend() {
    if [ -f "$BACKEND_PID_FILE" ]; then
        BACKEND_PID=$(cat "$BACKEND_PID_FILE")
        if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
            log_info "停止后端服务 (PID: $BACKEND_PID)..."
            kill "$BACKEND_PID" 2>/dev/null || true
            # 等待进程结束
            for i in {1..10}; do
                if ! ps -p "$BACKEND_PID" > /dev/null 2>&1; then
                    break
                fi
                sleep 0.5
            done
            # 如果还没停止，强制杀死
            if ps -p "$BACKEND_PID" > /dev/null 2>&1; then
                log_warn "强制停止后端服务..."
                kill -9 "$BACKEND_PID" 2>/dev/null || true
            fi
            log_success "后端服务已停止"
        else
            log_warn "后端服务未运行 (PID: $BACKEND_PID)"
        fi
        rm -f "$BACKEND_PID_FILE"
    else
        log_warn "未找到后端 PID 文件"
    fi

    # 额外清理：杀死可能遗留的 Go 进程
    pkill -f "go run cmd/server/main.go" 2>/dev/null && log_info "清理遗留的后端进程" || true
}

# 停止前端
stop_frontend() {
    if [ -f "$FRONTEND_PID_FILE" ]; then
        FRONTEND_PID=$(cat "$FRONTEND_PID_FILE")
        if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
            log_info "停止前端服务 (PID: $FRONTEND_PID)..."
            kill "$FRONTEND_PID" 2>/dev/null || true
            # 等待进程结束
            for i in {1..10}; do
                if ! ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
                    break
                fi
                sleep 0.5
            done
            # 如果还没停止，强制杀死
            if ps -p "$FRONTEND_PID" > /dev/null 2>&1; then
                log_warn "强制停止前端服务..."
                kill -9 "$FRONTEND_PID" 2>/dev/null || true
            fi
            log_success "前端服务已停止"
        else
            log_warn "前端服务未运行 (PID: $FRONTEND_PID)"
        fi
        rm -f "$FRONTEND_PID_FILE"
    else
        log_warn "未找到前端 PID 文件"
    fi

    # 额外清理：杀死可能遗留的 Vite 进程
    pkill -f "vite" 2>/dev/null && log_info "清理遗留的前端进程" || true
}

# 显示状态
show_status() {
    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║           Todo List 项目 - 停止服务                  ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# 主函数
main() {
    show_status

    # 停止服务
    stop_backend
    stop_frontend

    # 清理日志文件（可选）
    # rm -f "$PROJECT_DIR/backend.log" "$PROJECT_DIR/frontend.log"

    echo ""
    log_success "✅ 所有服务已停止"
    echo ""
}

# 运行主函数
main
