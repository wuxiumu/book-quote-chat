#!/usr/bin/env bash
set -e
echo "📚 Bootstrapping QuoteChat Frontend…"

# 1) Node 版本
command -v nvm >/dev/null 2>&1 || {
  echo "❌ nvm not found, please install first."; exit 1;
}
nvm install
nvm use

# 2) 安装依赖（干净模式）
npm ci

# 3) 启动
npm run dev