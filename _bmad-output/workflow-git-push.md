# Git 推送工作流：测试进展共享

> 目标：myorigin 收全量（代码+文档），origin 只收代码，互不干扰

## 仓库结构

| Remote | 地址 | 接收内容 |
|--------|------|----------|
| `origin` | `zstackio/terraform-provider-zstack` | 仅代码（测试代码、源码） |
| `myorigin` | `yi.qin/terraform-provider-zstack` | 全量（代码 + 测试文档 + 进度报告） |

## 分支策略

```
master          ← 代码变更（推 origin + myorigin）
test/progress   ← 进度文档（只推 myorigin）
```

- `master`：只放代码文件（`*.go`、`go.mod` 等），两个 remote 都推
- `test/progress`：基于 master，额外包含 `_bmad-output/`、`sprint*-test-results.md`、`.github/ISSUE_TEMPLATE/` 等测试文档，只推 myorigin

## 初始化（仅一次）

```bash
# 创建 test/progress 分支
git checkout -b test/progress master

# 提交现有文档
git add _bmad-output/ sprint2-test-results.md .github/ISSUE_TEMPLATE/
git commit -m "docs: initial test progress reports and bug tracking"

# 推送到 myorigin
git push -u myorigin test/progress

# 回到 master
git checkout master
```

## 日常工作流

### 完成一轮测试后的推送流程

```bash
# ── Step 1: 代码变更提交到 master ──
git add zstack/provider/*_test.go zstack/provider/*.go
git commit -m "test: <本次变更描述>"

# ── Step 2: 文档变更提交到 test/progress ──
git checkout test/progress
git merge master                          # 先同步代码
git add _bmad-output/ sprint2-test-results.md
git commit -m "docs: <本次进度更新描述>"

# ── Step 3: 推送 ──
git push myorigin test/progress           # 文档 → myorigin
git push myorigin master                  # 代码 → myorigin

# ── Step 4: 代码同步到 origin ──
git checkout master
git push origin master                    # 代码 → origin
```

### 只更新文档（无代码变更）

```bash
git checkout test/progress
git add _bmad-output/
git commit -m "docs: <更新描述>"
git push myorigin test/progress
git checkout master
```

### 只更新代码（无文档变更）

```bash
# 直接在 master 操作
git add zstack/provider/*.go zstack/provider/*_test.go
git commit -m "test: <变更描述>"
git push myorigin master
git push origin master
```

## 其他团队成员查看进度

```bash
# 拉取测试进度分支
git fetch myorigin test/progress
git checkout myorigin/test/progress

# 关键文件：
# _bmad-output/review-results.md   ← 分支审查 & Bug 清单
# _bmad-output/sprint-overview.md  ← Sprint 全局概览 & Story 跟踪
# sprint2-test-results.md          ← Sprint 2 测试结果
```

## 文件归属规则

| 文件/目录 | 归属分支 | 说明 |
|-----------|----------|------|
| `zstack/provider/*.go` | master | 源码 |
| `zstack/provider/*_test.go` | master | 测试代码 |
| `go.mod` / `go.sum` | master | 依赖 |
| `_bmad-output/*.md` | test/progress | 审查报告、Story、Sprint 概览 |
| `sprint*-test-results.md` | test/progress | 测试结果 |
| `.github/ISSUE_TEMPLATE/` | test/progress | Bug 模板 |

## 注意事项

1. **merge 方向固定**：始终是 `master → test/progress`（单向），不要反向 merge
2. **冲突概率极低**：两个分支改的文件集合不重叠
3. **feature/fix 分支**：正常合入 master，不影响此工作流
4. **test/progress 是长期分支**：不需要创建/删除，持续使用
