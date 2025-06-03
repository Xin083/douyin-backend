#!/bin/bash
# 删除 public/storage 文件或目录（无论是软链接、文件还是目录）
# if [ -e "./public/storage" ]; then
#   unlink -rf ./public/storage
#   echo "已删除 public/storage"
# else
#   echo "public/storage 不存在，无需删除"
# fi

# 删除 public/storage 软链接
if [ -L "public/storage" ]; then
  unlink public/storage
  echo "已删除 public/storage 软链接"
else
  echo "public/storage 不是软链接，无需删除"
fi
