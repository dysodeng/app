#!/bin/bash

HOOK_DIR=$(git rev-parse --git-path hooks)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 创建符号链接
ln -sf "$SCRIPT_DIR/hooks/pre-commit" "$HOOK_DIR/pre-commit"

# 确保 hook 脚本有执行权限
chmod +x "$SCRIPT_DIR/hooks/pre-commit"

echo "Git hooks installed successfully!"