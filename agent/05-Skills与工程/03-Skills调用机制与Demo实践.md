Skills 调用机制与 Demo 实践

# 1. 文档目标

本文系统整理 Agent Skills 的调用机制、组织方式、渐进式披露、Runtime 与模型的分工、是否需要 Function Calling、`load_skill` 工具的取舍，以及平铺式与树级 Skills 的适用场景。

本文重点回答以下问题：

1. `SKILL.md` 到底是什么？
2. 一个 `SKILL.md` 里一般是否包含多个 Skills？
3. Skills 是如何被 Agent Runtime 发现、选择、加载和执行的？
4. Skills 是否一定要转换成 OpenAI `tools` 并通过 Function Call 调用？
5. Skill 描述应该放在 system prompt、developer message、user prompt，还是 tools 中？
6. Runtime 内部路由 + 内部加载 Skill 是如何工作的？
7. `load_skill` tool 是否是最好的实现方式？
8. 模型与 Runtime 可能交互几次？
9. 平铺式 Skills 和树级 Skills 哪个更常见，分别适合什么场景？
10. 本项目中的平铺式与树级 demo 如何理解和运行？

# 2. Agent Skills 的基本定义

## 2.1. Skill 的本质

Agent Skill 是一种将专业能力封装为标准目录结构的工程化方式。它通常不是一个单纯函数，也不是一个 API schema，而是一套“能力说明书 + 可选脚本 + 可选资源”的组合。

一个典型 Skill 目录如下：

```text
weather-advice/
├── SKILL.md
├── scripts/
│   └── fixed_weather.go
└── resources/
    └── city-list.json
```

其中：

* `SKILL.md` 是 Skill 的入口文件，描述该 Skill 什么时候使用、如何使用、执行流程、输入输出要求、注意事项等。
* `scripts/` 是可选目录，用于存放可执行脚本，例如 Go、Python、Shell、Node.js 等。
* `resources/` 是可选目录，用于存放静态资源，例如模板、规则、JSON 数据、示例文件等。
* `reference.md`、`examples.md` 等同级 Markdown 文件可以作为补充资料，但建议只做一层引用，避免过深嵌套。

## 2.2. `SKILL.md` 的标准结构

一个 `SKILL.md` 通常包含 YAML frontmatter 和 Markdown 正文：

```markdown
---
name: weather-advice
description: 提供天气摘要和出行建议。当用户询问天气、温度、降雨、风力、穿衣或是否适合出门时使用。
---

# 天气建议 Skill

## 使用场景

当用户的问题可以直接归类为“天气查询 + 简短建议”时使用本 Skill。

## 执行流程

1. 从用户请求中提取城市。
2. 如果城市缺失，询问用户具体城市。
3. 运行脚本获取天气结果。
4. 根据脚本输出生成中文回答。
```

其中最关键的是：

* `name`：Skill 的唯一标识，通常使用小写字母、数字和连字符。
* `description`：Skill 的触发描述，Agent Runtime 或模型会根据它判断该 Skill 是否与当前用户请求相关。

`description` 很重要，因为它决定 Skill 是否会被发现和加载。好的 `description` 应该同时说明：

* 这个 Skill 能做什么。
* 什么场景下应该使用它。
* 包含足够明确的触发词，例如“天气、温度、降雨、风力、穿衣、适合出门”。

## 2.3. 一个 `SKILL.md` 一般对应一个 Skill

常见最佳实践是：

```text
一个目录 = 一个 Skill
一个 SKILL.md = 一个 Skill
```

例如：

```text
skills/
├── weather-advice/
│   └── SKILL.md
├── code-review/
│   └── SKILL.md
└── git-commit-message/
    └── SKILL.md
```

不建议在一个 `SKILL.md` 中塞入多个完整 Skill。原因是：

* Skill 的发现依赖 `name + description`，一个文件包含多个 Skill 会让语义边界变模糊。
* 渐进式披露希望“命中哪个 Skill 就加载哪个 Skill”，多个 Skill 放在一个文件会导致无关内容被一起加载。
* 一个 Skill 一个目录，便于维护、测试、版本管理和权限控制。

如果需要多个相关能力，推荐使用多个 Skill 目录，或者在规模很大时使用树级结构。

# 3. Skills 与 Function Calling 的区别

## 3.1. Function Calling 是模型调用函数的协议

Function Calling 的核心是：Runtime 将 `tools` schema 发给模型，模型返回结构化的 `tool_calls`，Runtime 执行函数后将结果返回给模型。

典型流程：

```text
Runtime 提供 tools schema
  ↓
模型返回 tool_calls
  ↓
Runtime 执行函数
  ↓
Runtime 将 tool result 返回模型
  ↓
模型生成最终回答
```

模型返回的内容类似：

```json
{
  "tool_calls": [
    {
      "function": {
        "name": "get_weather",
        "arguments": "{\"city\":\"北京\"}"
      }
    }
  ]
}
```

Function Calling 适合：

* 输入输出结构明确的函数。
* 单次调用即可完成的工具。
* 数量较少、schema 清晰的工具集。

例如：

```text
get_weather(city) -> weather json
search_docs(query) -> docs
calculate_tax(amount, region) -> tax
```

## 3.2. Skills 是组织专业能力的工程结构

Skills 的核心不是让模型直接调用函数，而是将专业流程、操作规范、脚本和资源组织成可复用模块。

典型流程：

```text
Runtime 扫描 Skill 描述
  ↓
Runtime 或模型判断命中哪个 Skill
  ↓
Runtime 读取完整 SKILL.md
  ↓
Runtime 按 Skill 说明执行脚本、工具或 API
  ↓
模型基于执行结果生成回答
```

Skills 适合：

* 多步骤流程。
* 专业操作规范。
* 需要脚本、模板、参考资料配合的任务。
* 需要沉淀经验和最佳实践的任务。

例如：

* 代码审查。
* Git commit message 生成。
* PDF 处理。
* 数据分析报告。
* 项目迁移。
* 故障排查流程。

## 3.3. Skills 不一定依赖 Function Calling

Skills 与 Function Calling 可以结合，但不是同一件事。

```text
Skills:
组织专业能力、流程、脚本、资源。

Function Calling:
模型调用结构化函数或 API。
```

一个 Skill 内部可以使用 Function Calling，例如：

```text
weather-advice/SKILL.md
  ↓
要求查询天气
  ↓
Runtime 调用 get_weather function tool
  ↓
模型基于结果生成回答
```

但 `weather-advice` 这个 Skill 本身不一定要变成一个 function tool。

## 3.4. 不建议把所有 Skills 都注册成 tools

小 demo 中可以将每个 Skill 转换成 OpenAI `tools`，但在常见实践中不推荐这样做。

原因：

* Skill 数量多时，所有 tools schema 会占用大量上下文。
* 模型一次看到过多工具，选择质量可能下降。
* 很多 Skill 是操作手册，不是简单函数，难以压缩成 JSON schema。
* `SKILL.md` 可能包含脚本、流程、输出规范、示例和注意事项，不适合全部放入 tool schema。
* Function Calling 更适合“明确输入输出的函数”，不适合复杂工作流。

更合理的做法是：

```text
Skills 不直接等于 tools。
Skills 可以在内部使用 tools。
```

# 4. 渐进式披露机制

## 4.1. 渐进式披露的含义

渐进式披露是 Agent Skills 的核心机制之一。它的目标是：不在一开始把所有 Skill 的完整内容塞进模型上下文，而是按需逐步加载。

典型层次：

```text
第一层：只暴露 Skill 的 name + description
第二层：命中后读取完整 SKILL.md
第三层：需要时读取 reference.md、examples.md、resources/
第四层：需要时运行 scripts/
```

这种方式可以减少上下文污染，降低 token 消耗，并提升模型对当前任务的关注度。

## 4.2. 模型不会自己加载文件

模型本身没有文件系统能力。它不会自己去读取本地的 `SKILL.md`。

真正加载文件的是 Agent Runtime。

完整链路是：

```text
模型判断需要某个 Skill
  ↓
Runtime 收到选择结果
  ↓
Runtime 读取对应 SKILL.md
  ↓
Runtime 将完整内容追加到下一轮模型上下文
```

因此，“模型加载 Skill”只是一个便于理解的说法，实际主体是 Runtime。

## 4.3. 脚本如何执行

`SKILL.md` 可能写明：

```bash
go run scripts/fixed_weather.go --city "<city>"
```

执行流程是：

```text
模型或 Runtime 从用户请求中提取 city=北京
  ↓
Runtime 将命令模板替换成实际命令
  ↓
Runtime 调用 Shell、子进程、MCP tool 或其他执行能力
  ↓
脚本在本地或远端环境运行
  ↓
Runtime 读取 stdout/stderr
  ↓
Runtime 将结果返回给模型
```

如果用 Go 自研 Runtime，脚本执行可能类似：

```go
cmd := exec.CommandContext(ctx, "go", "run", "scripts/fixed_weather.go", "--city", "北京")
output, err := cmd.Output()
```

所以：

```text
SKILL.md 负责告诉 Runtime 应该怎么做。
Runtime 负责真正执行。
模型负责理解说明、提取参数、生成回答。
```

# 5. Runtime 内部路由 + 内部加载 Skill

## 5.1. 这是常见产品化 Agent 的做法

在 Cursor、Claude Skills、IDE Agent 或自研 Agent Runtime 中，较常见的做法是：

```text
Runtime 负责 Skill 的发现、筛选、读取、脚本执行。
模型只看到 Runtime 决定暴露给它的内容。
```

也就是说，Runtime 像一个调度层，模型像一个推理和生成核心。

## 5.2. 启动时扫描 Skill 索引

Runtime 扫描 Skills 目录：

```text
skills/
├── weather-advice/
│   └── SKILL.md
├── code-review/
│   └── SKILL.md
└── git-commit-message/
    └── SKILL.md
```

它先只读取 frontmatter：

```markdown
---
name: weather-advice
description: 提供天气摘要和出行建议。当用户询问天气、温度、降雨、风力、穿衣或是否适合出门时使用。
---
```

形成索引：

```json
[
  {
    "name": "weather-advice",
    "description": "提供天气摘要和出行建议..."
  },
  {
    "name": "git-commit-message",
    "description": "根据代码变更生成清晰的 Git commit message..."
  }
]
```

这一步通常不读取完整 `SKILL.md`。

## 5.3. 用户请求进入后进行路由

用户请求：

```text
北京今天适合出门吗？
```

Runtime 可以用多种方式判断候选 Skill：

* 关键词匹配。
* 规则分类。
* embedding 相似度检索。
* 模型分类。
* 规则 + 检索 + 模型的混合策略。

例如：

```text
用户请求包含“今天”“适合出门”
weather-advice description 包含“天气、温度、是否适合出门”
=> 命中 weather-advice
```

## 5.4. 命中后读取完整 `SKILL.md`

命中 `weather-advice` 后，Runtime 才读取：

```text
skills/weather-advice/SKILL.md
```

然后将完整说明给模型：

```text
当前命中的 Skill 是 weather-advice。

请按照以下 Skill 说明完成任务：

# 天气建议 Skill

## 执行流程
...
```

## 5.5. Runtime 执行脚本或工具

如果完整 `SKILL.md` 要求运行脚本：

```bash
go run scripts/fixed_weather.go --city "<city>"
```

Runtime 可以：

1. 自己从用户请求提取城市。
2. 或让模型先提取参数。
3. 然后执行脚本。

脚本返回：

```json
{
  "city": "北京",
  "weather": "晴",
  "temperature": "26°C",
  "advice": "适合外出，建议注意防晒并适量补水。"
}
```

Runtime 再把结果提供给模型生成最终回答。

## 5.6. 常见完整流程

```text
Runtime 扫描 SKILL.md frontmatter
  ↓
建立 Skill 索引
  ↓
用户请求进入
  ↓
Runtime 或模型根据索引选择 Skill
  ↓
Runtime 读取完整 SKILL.md
  ↓
Runtime 按说明执行脚本/工具/API
  ↓
Runtime 将结果返回模型
  ↓
模型生成最终回答
```

一句话总结：

```text
Runtime 内部路由 + 内部加载 Skill，就是 Runtime 负责“找技能、读技能、执行技能”，模型负责“理解说明、提取参数、生成回答”。
```

# 6. 模型与 Runtime 可能交互几次

## 6.1. 可能只调用模型一次

如果 Runtime 自己完成以下工作：

* Skill 路由。
* 参数提取。
* 读取完整 `SKILL.md`。
* 执行脚本。

那么模型只需要最后生成回答。

流程：

```text
Runtime:
1. 用户问：北京今天适合出门吗？
2. 关键词/embedding 命中 weather-advice
3. 读取 weather-advice/SKILL.md
4. 执行 weather.go --city 北京
5. 得到天气 JSON
6. 将用户问题 + SKILL.md + JSON 一次性给模型

Model:
生成最终回答
```

这种情况下，模型交互只有一次。

## 6.2. 可能调用模型两次

如果 Runtime 需要模型帮助选择 Skill 或提取参数，通常会有两次调用。

第一次：

```text
Skill Index + User Request
  ↓
模型选择 Skill / 提取参数
```

第二次：

```text
Full SKILL.md + Script Result + User Request
  ↓
模型生成最终回答
```

示例：

第一次模型返回：

```json
{
  "action": "load_skill",
  "skill": "weather-advice",
  "arguments": {
    "city": "北京"
  },
  "reason": "用户询问天气和是否适合出门。"
}
```

Runtime 根据结果读取完整 `SKILL.md`，执行脚本后，再进行第二次模型调用。

## 6.3. 复杂任务可能多次调用模型

复杂 Skill 可能需要多轮协作：

```text
模型规划
  ↓
Runtime 执行脚本 A
  ↓
模型判断还缺什么
  ↓
Runtime 查询数据库
  ↓
模型继续分析
  ↓
Runtime 读取文件
  ↓
模型生成最终报告
```

适合多轮的任务包括：

* 代码审查。
* 数据分析。
* 故障排查。
* 项目迁移。
* 多文件生成或修改。

## 6.4. 调用次数由 Runtime 设计决定

Skills 并不规定必须一轮、两轮或多轮。

```text
简单 Skill：
Runtime 能自己路由、提参、执行，模型可能只调用一次。

中等 Skill：
模型负责提参或规划，通常两次。

复杂 Skill：
模型和 Runtime 多轮协作，可能多次。
```

# 7. Skill 描述如何给到模型

## 7.1. 模型能看到的只有 Runtime 提供的上下文

模型本身不知道本地有什么 Skills。Runtime 必须把 Skill 索引提供给模型。

Skill 索引通常类似：

```text
Available Skills:
- weather-advice: 提供天气摘要和出行建议。当用户询问天气、温度、降雨、风力、穿衣或是否适合出门时使用。
- git-commit-message: 根据代码变更生成清晰的 Git commit message。当用户要求写提交信息、总结 staged changes 或整理提交说明时使用。
```

模型根据这段索引和用户问题判断是否需要某个 Skill。

## 7.2. 不建议把 Skill 索引放到 tools 中

`tools` 更适合描述可调用函数，例如：

* `load_skill(name)`
* `run_shell(command)`
* `query_database(sql)`

Skill 索引本身只是上下文资料，不是函数定义，所以不适合放在 `tools` 中。

如果把 Skill 索引塞进 tools，会混淆两个概念：

```text
Skill 索引：说明有哪些能力。
tools：说明模型能调用哪些函数。
```

## 7.3. 动态 Skill 索引更像 developer/context message

逻辑上，Skill 索引属于 Runtime 提供的动态上下文：

```text
它不是用户真正说的话。
也不是固定不变的系统规则。
它是 Runtime 在本轮请求中提供给模型的资料。
```

因此最合理的消息分层是：

```text
system:
固定、长期不变的 Agent 基础规则。

developer/context:
本轮 Runtime 提供的动态上下文，例如 Skill 索引、路由规则、输出 JSON 格式。

user:
用户真实请求。
```

如果模型 API 支持 `developer` 角色，可以放在 developer message 中。

示例：

```json
{
  "role": "developer",
  "content": "Available Skills:\n- weather-advice: ...\n- git-commit-message: ..."
}
```

## 7.4. 哪些模型支持 developer message

`developer` 角色不是所有模型 API 都支持。

OpenAI 新一些的接口和模型通常支持 `developer` 角色，尤其是 Responses API 或新版 Chat Completions 中的部分模型，例如 GPT-4.1、GPT-4o、o 系列等。

但很多 OpenAI 兼容服务、第三方模型、旧模型不一定支持 `developer`，常见只支持：

```text
system
user
assistant
tool
```

有些只稳定支持：

```text
system
user
assistant
```

## 7.5. 如果不支持 developer message

如果 API 不支持 `developer` 角色，最通用做法是将 Runtime 动态上下文放入 user message 的上下文区域，并用标签隔开：

```text
<runtime_context>
Available Skills:
- weather-advice: 提供天气摘要和出行建议...
- git-commit-message: 根据代码变更生成提交信息...

Routing Instruction:
If a skill is needed, return JSON with selected_skill.
</runtime_context>

<user_request>
北京今天适合出门吗？
</user_request>
```

这种做法技术上属于 user prompt，但语义上仍然是 Runtime context。

## 7.6. system prompt 与缓存问题

动态 Skill 索引不建议频繁放入 system prompt，原因是：

* system prompt 通常更适合稳定、不变的基础规则。
* 如果每轮 system prompt 都变化，会降低 prompt cache 命中概率。
* 稳定前缀越稳定，越有利于缓存。

更好的原则是：

```text
稳定内容放前面。
动态内容放后面。
固定规则放 system。
Skill 索引放 developer/context 或 user 上下文块。
```

示例：

```text
System:
你是一个支持 Skills 的 Agent Runtime。遵守安全规则，按 Runtime 提供的上下文工作。

Developer/User Context:
<available_skills>
- weather-advice: ...
- git-commit-message: ...
</available_skills>

User:
北京今天适合出门吗？
```

# 8. 模型如何告诉 Runtime 加载哪个 Skill

## 8.1. 不一定通过 Function Call

模型告诉 Runtime 加载哪个 Skill，可以有多种方式，不一定是 Function Call。

常见方式包括：

1. 结构化 JSON 输出。
2. Function Calling 调用 `load_skill`。
3. Runtime 不问模型，自己路由。

## 8.2. 方式一：结构化 JSON 输出

Runtime 第一次调用模型时，可以要求模型只返回 JSON：

```text
请只返回 JSON：
{
  "action": "load_skill" | "answer_directly",
  "skill": "<skill-name or empty>",
  "arguments": {},
  "reason": "..."
}
```

模型返回：

```json
{
  "action": "load_skill",
  "skill": "weather-advice",
  "arguments": {
    "city": "北京"
  },
  "reason": "用户询问天气和是否适合出门。"
}
```

Runtime 解析 JSON 后加载：

```text
skills/weather-advice/SKILL.md
```

这种不是 Function Call，只是结构化文本协议。它适合自研 Agent demo 或简单 Runtime。

## 8.3. 方式二：Function Call 调用 `load_skill`

Runtime 也可以只暴露一个工具：

```json
{
  "name": "load_skill",
  "description": "Load full instructions for one or more selected skills.",
  "parameters": {
    "type": "object",
    "properties": {
      "skills": {
        "type": "array",
        "items": {
          "type": "string"
        }
      }
    },
    "required": ["skills"]
  }
}
```

模型看到 Skill 索引后调用：

```json
{
  "tool_calls": [
    {
      "function": {
        "name": "load_skill",
        "arguments": "{\"skills\":[\"weather-advice\"]}"
      }
    }
  ]
}
```

Runtime 执行 `load_skill`，读取完整：

```text
skills/weather-advice/SKILL.md
```

这种是 Function Calling 风格的 Skills Runtime。

## 8.4. 方式三：Runtime 自己路由

Runtime 也可以不问模型，直接用规则或检索选择 Skill。

例如：

```text
用户请求包含“天气”“适合出门”
=> weather-advice
```

这种方式下，模型不会告诉 Runtime 加载哪个 Skill，Runtime 自己决定。

## 8.5. 常见实践是混合策略

产品化 Agent 中常见做法是混合：

```text
Runtime 先用规则/embedding 检索筛出候选 Skills
  ↓
将少量候选 Skill 描述给模型判断
  ↓
模型用结构化输出或内部协议返回选择
  ↓
Runtime 加载完整 SKILL.md
```

是否使用 Function Call，取决于 Runtime 设计。

# 9. `load_skill` tool 是否是最好的方案

## 9.1. `load_skill` tool 是合理方案

提供一个统一的 `load_skill` tool 是一种工程上很干净的方案。

优点：

* 只需要注册一个 tool。
* 不需要把所有 Skill 都注册成 tools。
* 模型可以显式请求加载一个或多个 Skill。
* Runtime 可以统一做权限、日志、缓存、路径校验。
* 适合多 Skill 场景。

示例：

```json
{
  "name": "load_skill",
  "description": "Load full instructions for one or more selected skills.",
  "parameters": {
    "type": "object",
    "properties": {
      "skills": {
        "type": "array",
        "items": {
          "type": "string"
        }
      }
    },
    "required": ["skills"]
  }
}
```

## 9.2. `load_skill` tool 不是唯一最佳

它也有缺点：

* 增加一次 tool call 往返。
* 模型可能过度调用 `load_skill`。
* 需要防止模型加载不存在或不该访问的 Skill。
* Skill 选择更偏模型主导，Runtime 控制会弱一些。

因此更准确的说法是：

```text
一定有加载 Skill 的 Runtime 能力。
不一定有暴露给模型的 load_skill tool。
```

## 9.3. 什么时候适合使用 `load_skill` tool

适合：

* 已经有 Function Calling 基础设施。
* 想让模型显式决定加载哪些 Skills。
* Skill 数量较多，但不想把每个 Skill 都注册成 tools。
* 需要统一封装权限、日志、缓存、路径校验。

不适合：

* 追求最低延迟。
* Runtime 可以非常确定地自己路由。
* 希望 Runtime 强控制 Skill 选择。
* 模型经常过度加载无关 Skill。

## 9.4. 推荐设计

自研 Agent 中比较清晰的设计是：

```text
固定 system:
Agent 基础规则。

developer/context:
Skill 索引 + 路由要求 + 只允许从索引中选择 Skill。

tools:
少量 Runtime 能力，例如：
- load_skill
- run_script 或 run_command
```

流程：

```text
第一次模型调用:
看到 Skill 索引，调用 load_skill(["weather-advice"])

Runtime:
读取完整 SKILL.md 返回 tool result

第二次模型调用:
根据完整 SKILL.md 判断是否调用 run_script

Runtime:
执行脚本

第三次模型调用:
根据脚本结果生成最终回答
```

如果追求效率，也可以不提供 `load_skill` tool，由 Runtime 内部直接加载。

# 10. 第一次与第二次模型调用的 prompt 组织

## 10.1. 第一次：只给 Skill 索引

第一次的目的通常是选择 Skill 或提取参数，不是执行 Skill。

推荐输入结构：

```text
System:
你是一个支持 Skills 的 Agent。遵守安全规则，按 Runtime 提供的上下文工作。

Developer/Context:
你只能根据下面的 Skill 索引判断是否需要使用某个 Skill。
不要编造不存在的 Skill。
如果需要 Skill，请返回 JSON。

Available Skills:
- weather-advice: 提供天气摘要和出行建议。当用户询问天气、温度、降雨、风力、穿衣或是否适合出门时使用。
- git-commit-message: 根据代码变更生成清晰的 Git commit message。当用户要求写提交信息、总结 staged changes 或整理提交说明时使用。

返回格式：
{
  "skill": "<skill-name or none>",
  "arguments": {},
  "reason": "..."
}

User:
北京今天适合出门吗？
```

模型可能返回：

```json
{
  "skill": "weather-advice",
  "arguments": {
    "city": "北京"
  },
  "reason": "用户询问天气和是否适合出门，匹配 weather-advice。"
}
```

## 10.2. Runtime 收到第一次结果后

Runtime 做以下事情：

1. 解析模型返回的 JSON。
2. 确认 `skill = weather-advice`。
3. 根据 Skill 名称找到目录 `skills/weather-advice/`。
4. 读取 `skills/weather-advice/SKILL.md`。
5. 根据情况执行脚本。

如果第一次已经提取出参数：

```json
{
  "city": "北京"
}
```

Runtime 可以直接执行：

```bash
go run scripts/fixed_weather.go --city "北京"
```

## 10.3. 第二次：给完整 `SKILL.md` 和执行结果

第二次的目的通常是按 Skill 的完整规范生成最终回答。

推荐输入结构：

```text
System:
你是一个支持 Skills 的 Agent。遵守安全规则，按 Runtime 提供的上下文工作。

Developer/Context:
你现在已经命中 Skill: weather-advice。
请严格按照这个 Skill 的说明回答用户。

完整 SKILL.md:

---
name: weather-advice
description: 提供天气摘要和出行建议。当用户询问天气、温度、降雨、风力、穿衣或是否适合出门时使用。
---

# 天气建议 Skill

## 执行流程
1. 从用户请求中提取城市。
2. 如果城市缺失，询问用户具体城市。
3. 运行脚本获取固定天气结果。
4. 根据脚本输出生成不超过 2 句话的中文回答。

## 输出要求
- 直接回答用户问题。
- 不要编造脚本没有返回的数据。
- 如果用户只需要一句话，保持简短。

Skill 执行结果:

{
  "city": "北京",
  "weather": "晴",
  "temperature": "26°C",
  "wind": "东南风 2 级",
  "advice": "适合外出，建议注意防晒并适量补水。"
}

User:
北京今天适合出门吗？
```

模型最终回答：

```text
北京今天晴，约 26°C，东南风 2 级，适合出门。建议注意防晒并适量补水。
```

## 10.4. 为什么不建议全部放 user prompt

技术上可以全部放 user prompt，但更推荐分层：

```text
System:
固定基础规则。

Developer/Context:
Runtime 指令、Skill 索引、完整 SKILL.md、执行结果、输出规则。

User:
用户原始问题。
```

如果 API 没有 developer 角色，就将 Runtime context 放进 user message 的上下文块中，并明确标记：

```text
<runtime_context>
...
</runtime_context>

<user_request>
...
</user_request>
```

# 11. 平铺式 Skills

## 11.1. 平铺式是更常见的默认做法

平铺式是最常见的 Skills 组织方式。

结构：

```text
skills/
├── weather-advice/
│   └── SKILL.md
├── code-review/
│   └── SKILL.md
├── pdf-processing/
│   └── SKILL.md
└── git-commit-message/
    └── SKILL.md
```

每个目录都是独立 Skill，Agent 通过每个 `SKILL.md` 的 `description` 判断是否命中。

## 11.2. 平铺式适用场景

适合：

* Skill 数量较少，例如 5-20 个。
* 每个 Skill 是独立任务。
* 用户意图容易直接命中。
* Skill 之间交叉较少。
* 团队希望结构简单、维护成本低。

典型例子：

```text
“帮我写 commit message” -> git-commit-message
“帮我 review 这段代码” -> code-review
“把这个接口生成文档” -> api-doc-generator
“分析这个 PDF” -> pdf-extract
```

## 11.3. 平铺式最佳实践

推荐结构：

```text
skills/
├── code-review/
│   ├── SKILL.md
│   ├── standards.md
│   └── examples.md
├── git-commit-message/
│   ├── SKILL.md
│   └── examples.md
└── api-doc-generator/
    ├── SKILL.md
    ├── templates.md
    └── scripts/
        └── extract_openapi.go
```

注意：

* `SKILL.md` 保持精简。
* 详细规范放到同级 `standards.md`、`examples.md`。
* 脚本放到 `scripts/`。
* 资源放到 `resources/`。
* 文件引用尽量只做一层。

# 12. 树级 Skills

## 12.1. 树级是高级组织方式

树级 Skills 不是最常见默认做法，而是 Skill 数量变多、领域明显复杂时的高级组织方式。

结构示例：

```text
skills/
└── weather/
    ├── SKILL.md
    ├── travel/
    │   ├── SKILL.md
    │   ├── beach/
    │   │   └── SKILL.md
    │   └── mountain/
    │       └── SKILL.md
    └── commute/
        ├── SKILL.md
        └── driving/
            └── SKILL.md
```

核心思想：

```text
父节点只做路由。
中间节点继续细分。
叶子节点才真正执行任务。
```

## 12.2. 树级适用场景

适合：

* 几十到上百个 Skills。
* 同一领域下有明显子领域。
* 父级概念太宽，直接执行会混乱。
* 子场景执行流程差异很大。
* 需要减少一次性暴露的 Skill 数量。

例如天气领域：

```text
weather
├── travel
│   ├── beach
│   ├── mountain
│   └── city-trip
├── commute
│   ├── driving
│   ├── walking
│   └── public-transit
└── outdoor-sport
    ├── running
    ├── cycling
    └── hiking
```

## 12.3. 父级 Skill 只负责路由

父级 `skills/weather/SKILL.md` 示例：

```markdown
---
name: weather
description: 路由天气相关请求。当用户询问天气、温度、降雨、风力、旅游天气、通勤天气或户外活动天气时使用。
---

# 天气路由 Skill

## 子 Skill

- `travel`：当用户的问题和旅游、旅行、景点、酒店、航班、海边或山区出行相关时使用。
- `commute`：当用户的问题和上班、上学、开车、步行、公交、地铁通勤相关时使用。
- `outdoor-sport`：当用户的问题和跑步、骑行、徒步、登山、露营等户外运动相关时使用。

## 路由规则

根据用户意图选择一个最合适的子 Skill。
如果用户意图不明确，只问一个澄清问题。
```

父节点不负责查天气，也不负责生成最终建议，只负责判断下一步进入哪个子 Skill。

## 12.4. 中间节点继续细分

`skills/weather/travel/SKILL.md` 示例：

```markdown
---
name: travel-weather
description: 路由旅游天气请求。当用户询问景点、旅行、海边、山区、露营、航班或旅行计划中的天气问题时使用。
---

# 旅游天气路由 Skill

## 子 Skill

- `beach`：当用户提到海边、海岛、沙滩、游泳、晒太阳、海风、紫外线时使用。
- `mountain`：当用户提到登山、徒步、山区、海拔、露营、山路时使用。
- `city-trip`：当用户提到城市游、景点打卡、博物馆、餐厅、购物、城市步行时使用。
```

## 12.5. 叶子节点真正执行

`skills/weather/travel/beach/SKILL.md` 示例：

```markdown
---
name: beach-travel-weather
description: 提供海边旅游天气建议。当用户询问海边、海岛、沙滩、游泳、海风、紫外线或海边出行是否合适时使用。
---

# 海边旅游天气 Skill

## 执行流程

1. 提取用户要去的城市、海边、海岛或景点名称。
2. 如果地点缺失，询问用户具体地点。
3. 运行脚本获取固定海边天气结果。
4. 重点关注风力、紫外线、降雨、温度和是否适合下水。
5. 给出简短的海边出行建议。
```

叶子节点包含完整执行逻辑、脚本、输出格式和异常处理。

## 12.6. 树级调用流程

用户请求：

```text
周末去青岛海边玩，天气适合吗？
```

调用流程：

```text
第一阶段：只暴露顶层 weather
  ↓
命中 weather
  ↓
读取 weather/SKILL.md
  ↓
发现应该进入 travel
  ↓
读取 weather/travel/SKILL.md
  ↓
发现应该进入 beach
  ↓
读取 weather/travel/beach/SKILL.md
  ↓
执行 beach 叶子 Skill
  ↓
生成最终回答
```

树级结构非常符合渐进式披露，因为每次只暴露当前节点的下一层候选 Skill。

## 12.7. 不建议树级的情况

不适合：

* 只有几个 Skill。
* 只是为了分类好看。
* 叶子 Skill 内容高度重复。
* 路由层比执行层还复杂。
* 每次都要多走几层才能完成简单任务。

判断规则：

```text
如果 description 能一句话准确命中，平铺。
如果 description 变得很长，包含很多互斥场景，考虑拆分。
如果拆分后同一领域出现 10+ 个相关 Skill，考虑树级路由。
```

# 13. 平铺式与树级的选择建议

## 13.1. 默认选择平铺式

更常见、更推荐从平铺式开始。

原因：

* 简单。
* 直观。
* 易维护。
* 易调试。
* 更符合大多数项目的 Skill 数量规模。

适合少量、边界清晰的能力：

```text
skills/
├── code-review/
├── git-commit-message/
├── api-doc-generator/
└── pdf-extract/
```

## 13.2. 规模扩大后引入树级

当同一领域内 Skill 数量变多，且场景差异明显时，引入树级：

```text
skills/
├── code-review/
├── git-commit-message/
└── data-analysis/
    ├── SKILL.md
    ├── sql/
    ├── spreadsheet/
    └── dashboard/
```

最佳实践通常是混合：

```text
大多数 Skill 平铺。
只有复杂大领域做树级。
```

# 14. 本项目 Demo 说明

## 14.1. Demo 目录

当前示例目录位于：

```text
code-case/agent/skills/
├── README.md
├── flat-example/
└── tree-example/
```

只保留两个完整 demo：

* `flat-example`：平铺式 Skills demo。
* `tree-example`：树级路由 Skills demo。

## 14.2. 平铺式 demo

目录：

```text
flat-example/
├── main.go
├── README.md
├── weather-advice/
│   ├── SKILL.md
│   └── scripts/
│       └── fixed_weather.go
└── git-commit-message/
    └── SKILL.md
```

运行：

```bash
cd code-case/agent/skills/flat-example
GO111MODULE=off go run .
```

展示流程：

```text
1. main.go 扫描当前目录下所有包含 SKILL.md 的 Skill 目录。
2. 只读取每个 Skill 的 name 和 description。
3. 用户请求写死为：北京今天适合出门吗？
4. Runtime 根据关键词命中 weather-advice。
5. 命中后读取完整 weather-advice/SKILL.md。
6. 执行 weather-advice/scripts/fixed_weather.go。
7. 根据固定 JSON 结果生成最终回答。
```

这个 demo 模拟的是：

```text
Runtime 内部路由 + 内部加载 Skill + 内部执行脚本
```

它没有调用模型，也没有使用 Function Calling，目的是展示 Skills Runtime 的核心控制流程。

## 14.3. 树级 demo

目录：

```text
tree-example/
├── main.go
├── README.md
└── weather/
    ├── SKILL.md
    ├── travel/
    │   ├── SKILL.md
    │   └── beach/
    │       ├── SKILL.md
    │       └── scripts/
    │           └── fixed_beach_weather.go
    └── commute/
        ├── SKILL.md
        └── driving/
            ├── SKILL.md
            └── scripts/
                └── fixed_driving_weather.go
```

运行：

```bash
cd code-case/agent/skills/tree-example
GO111MODULE=off go run .
```

展示流程：

```text
1. 只暴露顶层 weather Skill。
2. 命中后读取 weather/SKILL.md。
3. 根据路由规则进入 weather/travel/SKILL.md。
4. 继续命中 weather/travel/beach/SKILL.md。
5. 执行 weather/travel/beach/scripts/fixed_beach_weather.go。
6. 根据固定 JSON 结果生成最终回答。
```

这个 demo 展示：

```text
weather -> travel -> beach
```

逐层渐进式披露的树级 Skill 调用方式。

# 15. 常见误区

## 15.1. 误区一：Skills 必须通过 Function Call 调用

不对。

Skills 可以通过 Function Calling 实现，但不必然依赖 Function Calling。

更常见的是 Runtime 自己发现、加载和执行 Skill。

## 15.2. 误区二：所有 Skills 都应该注册成 tools

不推荐。

Skill 是能力说明和流程封装，tools 是模型可调用函数。二者可以结合，但不应混为一谈。

## 15.3. 误区三：一个 `SKILL.md` 可以随意放很多 Skills

不推荐。

更好的边界是：

```text
一个目录 = 一个 Skill
一个 SKILL.md = 一个 Skill
```

如果需要层级，使用目录树表达，而不是把多个完整 Skill 堆进一个文件。

## 15.4. 误区四：模型会自己读取 `SKILL.md`

不对。

模型不会自己读取本地文件。Runtime 读取文件后再放进模型上下文。

## 15.5. 误区五：两轮调用就是 Function Call

不对。

两轮调用只是交互模式，不等于 Function Calling。

两轮也可以是普通文本协议：

```text
第一轮：模型返回 JSON，表示选择哪个 Skill。
第二轮：Runtime 加载 Skill 后让模型生成最终回答。
```

只有当模型返回标准 `tool_calls`，并由 Runtime 按工具协议执行时，才是 Function Calling。

# 16. 实践建议总结

## 16.1. 默认实践

建议默认使用：

```text
平铺式 Skills
Runtime 内部路由
Runtime 内部加载 SKILL.md
Runtime 内部执行脚本
模型负责理解和生成
```

## 16.2. Skill 数量较多时

当 Skills 数量增长后：

```text
先用 embedding/规则筛选候选 Skills
再把少量候选 Skill 描述给模型判断
必要时使用树级结构降低暴露数量
```

## 16.3. 需要模型显式加载 Skill 时

可以提供一个统一的 `load_skill` tool：

```text
load_skill(skills: string[])
```

但不要把每个 Skill 都注册成独立 tool。

## 16.4. prompt 组织建议

推荐：

```text
system:
固定基础规则，保持稳定，利于缓存。

developer/context:
动态 Skill 索引、路由规则、输出格式。

user:
用户真实请求。
```

如果不支持 developer：

```text
将 runtime_context 放入 user message 的独立标签块。
```

## 16.5. 平铺与树级选择

```text
少量 Skill：平铺。
边界清晰：平铺。
几十个以上：考虑树级。
同一领域内场景差异大：树级。
只有为了分类好看：不要树级。
```

最终建议：

```text
大多数项目从平铺开始。
复杂领域再局部引入树级。
Runtime 控制加载与执行。
Function Calling 作为可选实现手段，而不是 Skills 的必需前提。
```
