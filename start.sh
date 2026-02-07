#!/bin/bash

# Elysia LLM Tool 启动脚本

set -e

echo "======================================"
echo "  Elysia LLM Tool 启动脚本"
echo "======================================"

# 检查环境变量文件
if [ ! -f .env ]; then
    echo "警告: .env 文件不存在"
    echo "正在从 .env.example 创建 .env 文件..."
    cp .env.example .env
    echo "请编辑 .env 文件，填入你的 API Key"
    exit 1
fi

# 检查是否已生成proto代码
if [ ! -f proto/llm/llm.pb.go ]; then
    echo "检测到proto代码未生成，正在生成..."
    cd proto && make && cd ..
    echo "proto代码生成完成"
fi

# 安装依赖
echo "检查依赖..."
go mod download
go mod tidy

# 编译项目
echo "编译项目..."
go build -o elysia-llm-tool main.go

# 启动服务
echo "启动服务..."
./elysia-llm-tool
