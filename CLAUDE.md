
## Git Conventions

**Commit Message Format:**
- **Use Chinese for commit messages**
- Follow Conventional Commits specification with emoji prefixes
- Format: `<emoji> <type>[scope]: <description>`
- Add body for complex changes explaining what and why

**Commit Types & Emojis:**
| Type | Emoji | Description | Example |
|------|-------|-------------|---------|
| `feat` | ✨ | New feature | ✨ feat(采购): 添加采购订单审批功能 |
| `fix` | 🐛 | Bug fix | 🐛 fix(库存): 修复库存计算错误 |
| `refactor` | ♻️ | Code refactoring | ♻️ refactor(常量): 提取枚举常量到独立包 |
| `docs` | 📝 | Documentation | 📝 docs: 更新API文档 |
| `style` | 🎨 | Code style | 🎨 style: 统一代码格式 |
| `perf` | ⚡️ | Performance | ⚡️ perf: 优化数据库查询性能 |
| `test` | ✅ | Testing | ✅ test: 添加单元测试 |
| `chore` | 🔧 | Maintenance | 🔧 chore: 更新依赖版本 |
| `ci` | 👷 | CI/CD | 👷 ci: 优化CI配置 |

**Example Commit:**
```
♻️ refactor(常量): 将模型中的枚举常量提取到独立的常量包

- 按模型创建独立的常量文件(bom, inbound_order, inventory_lot等)
- 从模型文件中移除常量定义
- 更新handlers、repositories、api包中所有引用

这样可以将领域常量与模型定义分离，提高代码可维护性，减少模型与常量值之间的耦合。
```