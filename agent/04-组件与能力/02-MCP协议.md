MCP协议详解

# 1. MCP协议历史与演进

## 1.1. 背景与起源
随着大语言模型（LLM）能力的飞速发展，它们在理解和生成文本方面展现出前所未有的潜力。然而，LLM本身是“封闭”的，它们的能力局限于其训练数据，无法直接与外部世界进行实时交互、执行代码或调用外部服务。这极大地限制了LLM在实际应用中的价值。为了弥补这一鸿沟，业界开始探索如何为LLM构建一个标准化的“手脚”，使其能够与外部工具（Tools）进行高效、可靠的通信。Model Context Protocol (MCP) 正是在这样的背景下应运而生，旨在解决LLM与外部工具之间集成复杂、互操作性差的问题。

## 1.2. 发展历程与重要里程碑
MCP协议的演进是LLM与外部世界融合的必然趋势。其发展历程中，以下几个关键时间节点和事件值得关注：
*   **2023年**：OpenAI发布了“函数调用”API和ChatGPT插件框架，初步解决了LLM与外部工具集成的问题，但这些方案通常是供应商特定的，缺乏通用性。
*   **2024年11月**：Anthropic正式推出Model Context Protocol (MCP) 作为一个开放标准和开源框架。其目标是标准化AI系统（如LLM）与外部工具、系统和数据源的集成与数据共享方式，旨在解决信息孤岛和传统系统带来的“N×M”数据集成难题。MCP借鉴了语言服务器协议（LSP）的消息流思想，并基于JSON-RPC 2.0进行传输。
*   **2024年11月**：MCP发布时，同步推出了Python, TypeScript, C#, Java等多种编程语言的软件开发工具包（SDKs），并由Anthropic维护参考MCP服务器实现。
*   **2025年3月**：OpenAI正式宣布采纳MCP，并将其集成到包括ChatGPT桌面应用在内的产品中，这标志着MCP在行业内获得了广泛认可。
*   **2025年4月**：Google DeepMind也宣布拥抱Anthropic的MCP标准，进一步巩固了MCP作为行业通用协议的地位。然而，同期安全研究人员也发布分析报告，指出MCP存在多项安全问题，包括提示注入、工具权限可能导致数据泄露以及外观相似的工具可能被恶意替换等。
*   **2025年6月18日**：MCP发布了上一个规范版本，重点关注结构化工具输出、基于OAuth的授权、服务器发起的交互式用户引导以及改进的安全最佳实践。
*   **2025年9月**：MCP Registry预览版上线，这是一个开放的目录和API，用于索引和发现MCP服务器，进一步促进了MCP生态系统的发展。
*   **2025年11月25日**：MCP迎来了一周年纪念，并计划发布下一个规范版本，其发布候选版本（RC）定于2025年11月11日。新版本将致力于异步操作、无状态性和可扩展性、服务器身份识别以及官方扩展等方面的改进。
*   **2025年12月**：Anthropic将MCP捐赠给Agentic AI Foundation (AAIF)，这是一个由Linux基金会资助的基金，由Anthropic、Block和OpenAI共同创立，并得到其他公司的支持，旨在推动MCP的开放治理和社区发展。

MCP协议的出现，标志着LLM从单纯的“文本处理引擎”向“智能体（Agent）”迈出了关键一步，使其能够真正地成为一个能够感知、决策、行动的智能实体。

# 1. 02-MCP协议.md

## 1.1. 02-MCP协议.md
本文件详细介绍Model Context Protocol的架构设计、通信机制和实现原理，涵盖标准化模型工具交互协议。

# 2. 协议概述

## 2.1. MCP定义
Model Context Protocol (MCP) 是一个开放且标准化的协议，旨在定义大语言模型（LLM）与外部工具和服务之间进行高效、可靠通信的规范。它提供了一套统一的接口和约定，使得不同平台和技术栈的LLM能够无缝地调用各种外部功能。

### 没有MCP会遇到什么问题？
在没有MCP这类标准化协议的情况下，LLM与外部工具的集成会面临诸多挑战：
1.  **紧耦合与高成本**：每个LLM与工具的集成都需要定制开发，导致系统紧耦合，难以维护和扩展。
2.  **互操作性差**：不同LLM或工具之间缺乏统一的通信标准，难以实现跨平台、跨服务的互操作。
3.  **工具发现与管理困难**：LLM难以动态发现和理解可用的外部工具及其功能，限制了其能力边界。
4.  **安全性与可靠性风险**：缺乏统一的安全机制和错误处理规范，可能导致数据泄露、调用失败等问题。

### MCP如何解决这些问题？
MCP通过以下机制有效解决了上述问题：
1.  **标准化通信**：定义了统一的消息格式（如JSON）和请求-响应模式，确保LLM与工具之间的数据交换清晰、可预测。
2.  **松耦合架构**：通过抽象工具接口，将LLM的决策逻辑与工具的具体实现解耦，使得工具可以独立开发、部署和更新，提高了系统的灵活性和可维护性。
3.  **统一API接口**：提供了一套标准化的API接口，LLM可以通过这套接口调用任何符合MCP规范的工具，极大地简化了集成过程。
4.  **工具发现与描述**：通过工具定义（Tool Definition）机制，工具可以清晰地描述其功能、参数和预期行为，LLM可以动态地发现和理解这些工具，从而扩展其能力。
5.  **安全与可靠性保障**：协议中包含了错误处理、权限控制等机制，确保工具调用的安全性和稳定性，降低了潜在风险。
6.  **多传输协议支持**：支持多种底层传输协议（如stdio、socket、HTTP），使其能够适应不同的部署环境和性能需求。

# 3. 核心概念

## 3.1. 工具(Tools)
工具是MCP中的基本执行单元，每个工具提供特定的功能，如代码执行、文件操作、网络请求等。

## 3.2. 资源(Resources)
资源是MCP中的数据对象，提供只读访问权限，如文件内容、数据库查询结果、API响应等。

## 3.3. 提示(Prompts)
提示机制允许服务器向客户端请求额外信息，支持动态的交互式工具调用过程。

# 4. 通信机制

## 4.1. 请求响应模式
MCP采用标准的请求-响应通信模式，客户端发起工具调用请求，服务器返回执行结果。

## 4.2. 消息格式
消息使用JSON格式进行序列化，包含方法名、参数、返回值等字段，支持错误处理和状态跟踪。

## 4.3. 传输协议
MCP支持多种传输协议，包括stdio、socket、HTTP等，适应不同的部署环境和性能要求。

# 5. MCP详细运行原理

## 5.1. 协议架构
MCP协议的设计精妙之处在于其对LLM与外部Tool交互流程的标准化和异步处理能力。在深入理解其运行机制之前，我们首先明确MCP中的两个核心主体：**客户端（Client）** 和 **服务端（Server）**。

*   **MCP客户端 (Client)**：通常是集成在大语言模型（LLM）应用中的模块，负责发起Tool调用请求、接收Tool执行结果以及管理与MCP服务端的连接。它代表了LLM与外部Tool交互的“发起方”和“接收方”。
*   **MCP服务端 (Server)**：负责托管和管理具体的外部Tool。它接收来自MCP客户端的Tool调用请求，执行相应的Tool逻辑，并将结果返回给客户端。MCP服务端是外部Tool的“提供方”和“执行方”。

其核心运行原理可以概括为以下几个关键环节：

## 5.2. 客户端与服务端的连接与发现
MCP客户端在启动时连接MCP服务端，服务端通过SSE连接发送endpoint事件告知Tool调用URL，实现客户端与服务端的连接与Tool发现。

### 5.2.1. 客户端启动与连接
    *   在智能体应用启动时，MCP客户端会根据配置（例如，项目启动时输入的值）尝试连接MCP服务端。这可以是连接一个远程的HTTP服务器，也可以是启动一个本地的MCP服务器实例。
    *   **MCP服务端的加载方式**：
        *   **直接调用HTTP远程服务器**：这是最简单、最常见的部署方式，客户端直接通过HTTP请求与远程部署的MCP服务端通信。
        *   **通过指令运行本地服务器代码**：
            *   **代码来自GitHub等远程仓库**：开发者可以将MCP服务端代码拉取到本地，并通过命令行指令运行。这种方式的缺点是，如果需要修改服务端逻辑，需要先修改代码再重新运行。
            *   **代码是本地可修改的**：开发者可以直接在本地开发和修改MCP服务端代码，然后运行。
        *   **关于服务器运行位置的疑问**：通过指令运行的本地MCP服务器，通常是作为一个独立的进程运行，而不是客户端服务的一个子进程。这意味着它们拥有独立的生命周期和资源，可以独立启动、停止和管理。

### 5.2.2. SSE连接与Endpoint发现
    *   MCP客户端（特别是HTTP客户端）会与MCP服务端建立一个**Server-Sent Events (SSE)** 连接。这是一个持久化的单向连接，允许服务端主动向客户端推送事件。
    *   当服务端检测到有新的客户端通过SSE连接到自己时，它会发送一个特殊的 `endpoint` 事件。这个事件会告诉客户端MCP服务端用于接收Tool调用请求的URL（即Endpoint），类似于服务端向客户端声明“我有哪些Tool可以调用，你可以向我发送Tool调用请求了”。
    *   **流程示意图**：
        ```
        [客户端] <───────────────────────────────────────────> [服务器]
            │                                                      │
            │ 1. 建立SSE连接                                        │
            │ ──────────────────────────────────────────────────>  │
            │                                                      │
            │ 2. 返回endpoint事件（仅一次，告知Tool调用URL）         │
            │ <──────────────────────────────────────────────────  │
        ```

## 5.3. Tool调用与结果返回

1.  **LLM的Tool选择与参数准备**：
    *   智能体中的LLM在处理用户请求时，如果判断需要调用外部Tool来完成任务，它会根据Tool的描述（通过MCP客户端的 `list_tools` 功能获取）选择合适的Tool，并准备好调用该Tool所需的参数。
    *   LLM客户端需要具备聊天功能，并且在聊天过程中能够识别并携带Tool调用信息（即Function Call）。

2.  **客户端发起Tool调用**：
    *   MCP客户端接收到LLM的Tool调用指令后，会使用之前通过 `endpoint` 事件获取到的URL，向该URL发送一个POST请求。这个请求包含了要调用的Tool名称和LLM提供的参数。
    *   MCP客户端需要提供 `list_tools`（列出所有可用Tool）和 `call_tool`（执行Tool调用）这两个核心功能。

3.  **服务端处理与结果推送**：
    *   MCP服务端接收到Tool调用请求后，会执行相应的Tool逻辑。
    *   Tool执行完成后，服务端会将执行结果封装成一个或多个 `message` 事件，并通过之前建立的SSE连接推送给客户端。
    *   **流程示意图（续）**：
        ```
        [客户端] <───────────────────────────────────────────> [服务器]
            │                                                      │
            │ 3. 向endpoint URL发送POST请求（携带Tool调用信息）    │
            │ ──────────────────────────────────────────────────>  │
            │                                                      │
            │ 4. 通过SSE发送message事件（可多次，返回Tool执行结果）  │
            │ <──────────────────────────────────────────────────  │
        ```

## 5.4. 读写分离机制

MCP协议的一个重要特性是其**读写分离**的设计。这意味着：
*   **写操作（Tool调用）**：通常通过HTTP POST请求发送到特定的Endpoint URL，是同步的请求-响应模式。
*   **读操作（Tool执行结果、事件通知）**：通过SSE连接进行，是异步的、服务端推送模式。

这种分离确保了Tool调用的可靠性，同时又能够高效地将Tool执行的实时进展和最终结果异步地通知给客户端，提高了系统的响应性和并发处理能力。

## 5.5. 客户端与LLM客户端的协同

为了实现完整的智能体Tool调用能力，通常需要两个客户端协同工作：
1.  **MCP客户端**：负责与MCP服务端建立连接、发现Tool、发送Tool调用请求以及接收Tool执行结果。它专注于MCP协议层面的通信。
2.  **LLM客户端**：负责与大语言模型进行交互，处理聊天功能。它需要能够解析LLM返回的Function Call指令，并将其转发给MCP客户端执行；同时，也需要将MCP客户端返回的Tool执行结果反馈给LLM，以便LLM继续生成响应。

通过上述机制，MCP协议为智能体提供了一个强大、灵活且可扩展的Tool调用框架，使得LLM能够真正地“走出”文本世界，与外部服务和数据进行深度交互。

# 6. 示例文档与资源

为了帮助开发者更好地理解和应用MCP协议，以下是一些官方及社区提供的示例文档和资源：

## 6.1. MCP Server实现

*   **GitHub MCP Server (Go语言)**：官方发布的MCP Server实现，采用Go语言编写，提供了高性能和可靠性。
    `https://github.com/github/github-mcp-server`
*   **Python MCP Server**：一个Python实现的MCP服务端示例，展示了如何用Python构建MCP兼容的服务。
    `https://github.com/modelcontextprotocol/servers/blob/main/src/fetch/src/mcp_server_fetch/server.py`
*   **官方StreamableHttp实现 (TypeScript)**：TypeScript版本的官方StreamableHttp实现，为Web环境下的MCP通信提供了参考。
    `https://github.com/modelcontextprotocol/typescript-sdk/blob/main/src/server/streamableHttp.ts`

## 6.2. MCP 文档与社区

*   **Claude MCP 官方文档**：详细介绍了MCP协议的各个方面，是学习MCP的权威资源。
    `https://www.claudemcp.com/zh/docs/introduction`
*   **GitHub上的MCP资源集合**：由社区维护的MCP相关资源集合，包含了各种示例和参考资料。
    `https://github.com/cyanheads/model-context-protocol-resources`
*   **MCP社区平台**：`https://mcp.so/` 是一个专注于Model Context Protocol (MCP) 的社区驱动平台，主要功能是为AI开发者提供第三方MCP服务器的集中管理与集成服务。

## 6.3. 客户端SDK与工具

*   **官方Python SDK（推荐首选）**：用于与MCP服务器进行交互的官方Python开发工具包，推荐开发者优先使用。
    `https://github.com/modelcontextprotocol/python-sdk`
    **安装方式**：`pip install mcp-client`
*   **社区高级客户端 (Ultimate MCP Client)**：
    GitHub地址：`https://github.com/Dicklesworthstone/ultimate_mcp_client`
    **特色功能**：
    *   自动服务发现（mDNS/Zeroconf）
    *   双协议支持（SSE/Stdio）
    *   Web可视化调试界面
*   **CherryStudio 官网**：`https://docs.cherry-ai.com/`
*   **CherryStudio GitHub**：`https://github.com/CherryHQ/cherry-studio`
*   **CherryStudio 文档**：`https://github.com/CherryHQ/cherry-studio-docs/blob/main/README.md`

## 6.4. 应用案例

*   **掘金关于接入高德地图的文档**：一个实际应用案例，展示了如何通过MCP协议接入高德地图服务。
    `https://juejin.cn/post/7487810035385368639`

# 5. 示例文档与资源

为了帮助开发者更好地理解和应用MCP协议，以下是一些官方及社区提供的示例文档和资源：

## 5.1. MCP Server实现

*   **GitHub MCP Server (Go语言)**：官方发布的MCP Server实现，采用Go语言编写，提供了高性能和可靠性。
    `https://github.com/github/github-mcp-server`
*   **Python MCP Server**：一个Python实现的MCP服务端示例，展示了如何用Python构建MCP兼容的服务。
    `https://github.com/modelcontextprotocol/servers/blob/main/src/fetch/src/mcp_server_fetch/server.py`
*   **官方StreamableHttp实现 (TypeScript)**：TypeScript版本的官方StreamableHttp实现，为Web环境下的MCP通信提供了参考。
    `https://github.com/modelcontextprotocol/typescript-sdk/blob/main/src/server/streamableHttp.ts`

## 5.2. MCP 文档与社区

*   **Claude MCP 官方文档**：详细介绍了MCP协议的各个方面，是学习MCP的权威资源。
    `https://www.claudemcp.com/zh/docs/introduction`
*   **GitHub上的MCP资源集合**：由社区维护的MCP相关资源集合，包含了各种示例和参考资料。
    `https://github.com/cyanheads/model-context-protocol-resources`
*   **MCP社区平台**：`https://mcp.so/` 是一个专注于Model Context Protocol (MCP) 的社区驱动平台，主要功能是为AI开发者提供第三方MCP服务器的集中管理与集成服务。

## 5.3. 客户端SDK与工具

*   **官方Python SDK（推荐首选）**：用于与MCP服务器进行交互的官方Python开发工具包，推荐开发者优先使用。
    `https://github.com/modelcontextprotocol/python-sdk`
    **安装方式**：`pip install mcp-client`
*   **社区高级客户端 (Ultimate MCP Client)**：
    GitHub地址：`https://github.com/Dicklesworthstone/ultimate_mcp_client`
    **特色功能**：
    *   自动服务发现（mDNS/Zeroconf）
    *   双协议支持（SSE/Stdio）
    *   Web可视化调试界面
*   **CherryStudio 官网**：`https://docs.cherry-ai.com/`
*   **CherryStudio GitHub**：`https://github.com/CherryHQ/cherry-studio`
*   **CherryStudio 文档**：`https://github.com/CherryHQ/cherry-studio-docs/blob/main/README.md`

## 5.4. 应用案例

*   **掘金关于接入高德地图的文档**：一个实际应用案例，展示了如何通过MCP协议接入高德地图服务。
    `https://juejin.cn/post/7487810035385368639`