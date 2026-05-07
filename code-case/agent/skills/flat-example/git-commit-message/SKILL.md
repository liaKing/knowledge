---
name: git-commit-message
description: 根据代码变更生成清晰的 Git commit message。当用户要求写提交信息、总结 staged changes 或整理提交说明时使用。
---

# Git Commit Message Skill

## 使用场景

当用户需要根据代码变更生成提交信息时使用本 Skill。

示例：

- `帮我写一个 commit message`
- `根据当前 diff 生成提交说明`
- `总结 staged changes`

## 执行流程

1. 查看当前 Git 变更。
2. 判断变更类型：新增功能、修复、重构、文档、测试或工程配置。
3. 生成一句简洁的提交标题。
4. 如有必要，补充一段正文说明为什么这样改。

## 固定示例输出

本示例不真实读取 Git diff，执行结果可以写死为：

```text
feat(agent): add skill examples

Add flat and tree-based skill examples to explain common skill organization patterns.
```

## 输出要求

- 标题使用英文，保持简洁。
- 正文说明变更目的，而不是逐文件罗列。
- 不要把密钥、`.env`、凭证文件写进提交说明。
