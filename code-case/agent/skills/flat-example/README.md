# 平铺式 Skills 示例

这个目录是一个可运行的平铺式 Skills demo：每个目录就是一个独立 Skill，Agent 通过每个 `SKILL.md` 的 `description` 直接判断是否命中。

```text
flat-example/
├── main.go
├── weather-advice/
│   ├── SKILL.md
│   └── scripts/
│       └── fixed_weather.go
└── git-commit-message/
    └── SKILL.md
```

## 适用场景

平铺式适合 Skill 数量较少、边界清晰、用户意图容易直接判断的场景，例如：

- 查询天气并给出建议
- 生成 Git commit message
- 代码审查
- 生成接口文档

## 运行方式

```bash
GO111MODULE=off go run .
```

## 调用流程

1. `main.go` 扫描当前目录下所有包含 `SKILL.md` 的 Skill 目录。
2. Agent 先读取所有 Skill 的 `name` 和 `description`。
2. 用户请求进入后，Agent 根据 `description` 判断命中哪个 Skill。
3. 命中后读取该 Skill 的完整 `SKILL.md`。
4. 如果 `SKILL.md` 要求执行脚本，再运行对应脚本。
5. Agent 根据脚本结果和输出要求生成最终回答。

## Demo 输入

```text
北京今天适合出门吗？
```

这个输入会命中 `weather-advice`，并执行：

```bash
go run scripts/fixed_weather.go --city "北京"
```
