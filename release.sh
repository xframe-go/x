#!/bin/bash

# Go Module 版本发布脚本
# 用法: ./release.sh [major|minor|patch|auto]
#   major  - 主版本号升级 (如: 1.0.0 -> 2.0.0)
#   minor  - 次版本号升级 (如: 1.0.0 -> 1.1.0)
#   patch  - 补丁版本号升级 (如: 1.0.0 -> 1.0.1，默认)
#   auto   - 自动判断版本类型 (根据 commit 消息)

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 默认补丁版本升级
VERSION_TYPE="${1:-patch}"

# 获取当前分支
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# 检查是否有未提交的更改
if ! git diff-index --quiet HEAD --; then
    echo -e "${RED}❌ 有未提交的更改，请先提交${NC}"
    git status --short
    exit 1
fi

# 检查是否有未推送的提交
UNPUSHED=$(git log @{u}..HEAD --oneline 2>/dev/null || true)
if [ -n "$UNPUSHED" ]; then
    echo -e "${YELLOW}⚠️  检测到未推送的提交:${NC}"
    echo "$UNPUSHED"
    echo ""
fi

# 获取最新版本号
get_latest_version() {
    local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

    if [ -z "$latest_tag" ]; then
        echo "0.0.0"
    else
        # 去掉 v 前缀
        echo "${latest_tag#v}"
    fi
}

# 解析版本号
parse_version() {
    local version=$1
    local major=$(echo "$version" | cut -d. -f1)
    local minor=$(echo "$version" | cut -d. -f2)
    local patch=$(echo "$version" | cut -d. -f3)
    echo "$major $minor $patch"
}

# 自动判断版本类型
auto_detect_version_type() {
    local last_commit=$(git log -1 --pretty=%B)
    local last_commit_lower=$(echo "$last_commit" | tr '[:upper:]' '[:lower:]')

    if echo "$last_commit_lower" | grep -qE "^feat!|break|major|breaking change"; then
        echo "major"
    elif echo "$last_commit_lower" | grep -qE "^feat|minor|new feature"; then
        echo "minor"
    else
        echo "patch"
    fi
}

# 应用版本升级
bump_version() {
    local current=$1
    local type=$2

    read -r major minor patch <<< $(parse_version "$current")

    case "$type" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo -e "${RED}❌ 未知的版本类型: $type${NC}"
            exit 1
            ;;
    esac

    echo "$major.$minor.$patch"
}

# 确认操作
confirm() {
    local message=$1
    echo -en "${YELLOW}🤔 $message [y/N]${NC} "
    read -r response
    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

# 主流程
main() {
    echo -e "${BLUE}📦 Go Module 版本发布脚本${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    # 获取当前版本
    CURRENT_VERSION=$(get_latest_version)

    # 自动检测版本类型
    if [ "$VERSION_TYPE" = "auto" ]; then
        VERSION_TYPE=$(auto_detect_version_type)
        echo -e "${GREEN}🔍 自动检测版本类型: ${BOLD}$VERSION_TYPE${NC}"
    fi

    # 计算新版本
    NEW_VERSION=$(bump_version "$CURRENT_VERSION" "$VERSION_TYPE")
    NEW_TAG="v$NEW_VERSION"

    # 显示信息
    echo -e "当前版本:  ${YELLOW}$CURRENT_VERSION${NC}"
    echo -e "版本类型:  ${YELLOW}$VERSION_TYPE${NC}"
    echo -e "新版本:    ${GREEN}$NEW_VERSION${NC}"
    echo -e "新标签:    ${GREEN}$NEW_TAG${NC}"
    echo -e "当前分支:  ${YELLOW}$CURRENT_BRANCH${NC}"
    echo ""

    # 确认
    if ! confirm "确认发布版本 $NEW_TAG ?"; then
        echo -e "${RED}❌ 操作已取消${NC}"
        exit 0
    fi

    echo ""

    # 推送到远端
    echo -e "${BLUE}📤 推送代码到远端...${NC}"
    if git push origin "$CURRENT_BRANCH"; then
        echo -e "${GREEN}✅ 代码推送成功${NC}"
    else
        echo -e "${RED}❌ 代码推送失败${NC}"
        exit 1
    fi

    # 创建标签
    echo -e "${BLUE}🏷️  创建标签 $NEW_TAG...${NC}"
    if git tag -a "$NEW_TAG" -m "Release $NEW_TAG"; then
        echo -e "${GREEN}✅ 标签创建成功${NC}"
    else
        echo -e "${RED}❌ 标签创建失败${NC}"
        exit 1
    fi

    # 推送标签
    echo -e "${BLUE}📤 推送标签到远端...${NC}"
    if git push origin "$NEW_TAG"; then
        echo -e "${GREEN}✅ 标签推送成功${NC}"
    else
        echo -e "${RED}❌ 标签推送失败${NC}"
        exit 1
    fi

    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}🎉 版本 $NEW_TAG 发布成功！${NC}"
    echo ""
    echo -e "其他用户可以通过以下命令获取新版本:"
    echo -e "${YELLOW}go get github.com/xframe-go/x@$NEW_TAG${NC}"
    echo -e "${YELLOW}go get -u github.com/xframe-go/x${NC}"
}

main "$@"
