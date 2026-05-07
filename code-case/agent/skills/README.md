# Skills Demo

这个目录只保留两个完整 demo：

```text
skills/
├── flat-example/  # 平铺式 Skills demo
└── tree-example/  # 树级路由 Skills demo
```

## 平铺式 Demo

运行：

```bash
cd flat-example
GO111MODULE=off go run .
```

展示流程：

1. 扫描同级目录下的 `SKILL.md`。
2. 只读取每个 Skill 的 `name` 和 `description` 作为索引。
3. 根据用户请求命中 `weather-advice`。
4. 命中后读取完整 `weather-advice/SKILL.md`。
5. 执行 `weather-advice/scripts/fixed_weather.go`。
6. 基于固定 JSON 结果生成最终回答。

## 树级 Demo

运行：

```bash
cd tree-example
GO111MODULE=off go run .
```

展示流程：

1. 先只暴露顶层 `weather` Skill。
2. 命中后读取 `weather/SKILL.md`。
3. 根据路由规则进入 `weather/travel/SKILL.md`。
4. 继续命中叶子 Skill `weather/travel/beach/SKILL.md`。
5. 执行 `weather/travel/beach/scripts/fixed_beach_weather.go`。
6. 基于固定 JSON 结果生成最终回答。
