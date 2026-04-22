.PHONY: release release-major release-minor release-patch release-auto help

# 版本发布快捷命令
release:  ## 发布补丁版本 (vX.X.X -> vX.X.(X+1))
	./release.sh patch

release-major:  ## 发布主版本 (vX.X.X -> v(X+1).0.0)
	./release.sh major

release-minor:  ## 发布次版本 (vX.X.X -> vX.(X+1).0)
	./release.sh minor

release-auto:  ## 自动判断版本类型 (基于 commit 消息)
	./release.sh auto

help:  ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "示例:"
	@echo "  make release       # 发布补丁版本"
	@echo "  make release-minor # 发布次版本"
	@echo "  make release-auto  # 自动判断版本类型"
