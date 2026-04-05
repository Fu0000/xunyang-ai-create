#!/bin/bash
# 寻氧AI 服务端部署脚本
set -e

PROJECT_DIR="/opt/xiaoye-ai"
REPO_URL="git@github.com:Fu0000/xunyang-ai-create.git"
SSH_KEY="/root/.ssh/id_ed25519"
BACKEND_PORT=8092

echo "===> [1/6] 拉取/更新代码..."
if [ -d "$PROJECT_DIR/.git" ]; then
  cd $PROJECT_DIR && git pull origin main
else
  git clone $REPO_URL $PROJECT_DIR
  cd $PROJECT_DIR
fi

echo "===> [2/6] 检查 .env.prod 配置..."
if [ ! -f "$PROJECT_DIR/backend/.env.prod" ]; then
  echo "ERROR: 请先创建 $PROJECT_DIR/backend/.env.prod 文件！"
  echo "模板："
  cat "$PROJECT_DIR/backend/.env"
  exit 1
fi

echo "===> [3/6] 构建并启动后端容器..."
cd $PROJECT_DIR
docker compose down --remove-orphans 2>/dev/null || true
docker compose up -d --build

echo "===> [4/6] 构建前端..."
cd $PROJECT_DIR/frontend
npm install --legacy-peer-deps
npm run build

echo "===> [5/6] 构建前端管理后台..."
cd $PROJECT_DIR/frontend-admin
npm install --legacy-peer-deps
npm run build

echo "===> [6/6] 拷贝前端静态资源..."
mkdir -p /var/www/xiaoye-ai/frontend
mkdir -p /var/www/xiaoye-ai/admin
cp -r $PROJECT_DIR/frontend/dist/* /var/www/xiaoye-ai/frontend/
cp -r $PROJECT_DIR/frontend-admin/dist/* /var/www/xiaoye-ai/admin/

echo ""
echo "✅ 部署完成！"
echo "- 后端 API:   http://8.140.214.182:$BACKEND_PORT"
echo "- 前端静态:   /var/www/xiaoye-ai/frontend/"
echo "- 管理后台:   /var/www/xiaoye-ai/admin/"
echo ""
echo "请确保 Nginx 配置已指向正确路径并重载。"
