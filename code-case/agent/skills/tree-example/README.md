# 树级 Skills 示例

这个目录是一个可运行的树级 Skills demo：父节点只做路由，中间节点继续细分，叶子节点才真正执行任务。

```text
tree-example/
├── main.go
└── weather/
    ├── SKILL.md
    ├── travel/
    │   ├── SKILL.md
    │   └── beach/
    │       ├── SKILL.md
    │       └── scripts/
    │           └── fixed_beach_weather.go
    └── commute/
        └── driving/
            ├── SKILL.md
            └── scripts/
                └── fixed_driving_weather.go
```

## 适用场景

树级适合一个大领域下有很多细分 Skill，且不同场景的执行流程明显不同。

例如天气领域可以继续分成：

- 旅游天气
- 通勤天气
- 户外运动天气
- 农业天气
- 施工天气

## 运行方式

```bash
GO111MODULE=off go run .
```

## 调用流程

1. `main.go` 先只暴露顶层 `weather` 的 `name` 和 `description`。
2. 命中后读取 `weather/SKILL.md`，只暴露下一层子 Skill。
3. 如果用户意图是旅游，继续读取 `weather/travel/SKILL.md`。
4. 如果进一步命中海边旅游，读取 `weather/travel/beach/SKILL.md`。
5. 叶子 Skill 执行脚本并生成最终回答。

## Demo 输入

```text
周末去青岛海边玩，天气适合吗？
```

这个输入会走完整路由：

```text
weather -> travel -> beach
```

最后执行：

```bash
go run scripts/fixed_beach_weather.go --place "青岛海边"
```
