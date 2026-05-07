LiteLLM接口

# 1. 文档

swagger 页面链接：<http://18.143.27.161:4000/>
测试环境base url：<http://18.143.27.161:4000/>

# 2. 接口列表

## 2.1. TOP1:

### 2.1.1. chat/completions

接口路径：V1/chat/completions
聊天接口示例:

```
curl -X POST 'http://18.143.27.161:4000/v1/chat/completions' -H 'accept: application/json' -H 'x-litellm-api-key: sk-UgfOWp_fazEt9Np2aqHveA' -H 'Content-Type: application/json' -d '{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "Hello, how are you?"}]}'
```

## 2.1.2. key management

接口路径：POST /key/generate
创建 key，需要使用管理员权限

```
curl -X 'POST' \  'http://18.143.27.161:4000/key/generate' \  -H 'accept: application/json' \  -H 'x-litellm-api-key: sk-1234' \  -H 'Content-Type: application/json' \  -d '{  "key_alias": "alias-asdf",  "duration": "30d",  "models": [],  "spend": 0,  "max_budget": null,  "user_id": null,  "team_id": null,  "agent_id": null,  "max_parallel_requests": null,  "metadata": {},  "tpm_limit": null,  "rpm_limit": null,  "budget_duration": "30d",  "allowed_cache_controls": [],  "config": {},  "permissions": {},  "blocked": false,  "aliases": {},  "budget_id": null,  "allowed_routes": [],  "key_type": "default",  "auto_rotate": false}'
```

参数

```
{
  "key_alias": "alias-asdf",
  "duration": "30d",
  "models": [],
  "spend": 0,
  "max_budget": null,
  "user_id": null,
  "team_id": null,
  "agent_id": null,
  "max_parallel_requests": null,
  "metadata": {},
  "tpm_limit": null,
  "rpm_limit": null,
  "budget_duration": "30d",
  "allowed_cache_controls": [],
  "config": {},
  "permissions": {},
  "blocked": false,
  "aliases": {},
  "budget_id": null,
  "allowed_routes": [],
  "key_type": "default",
  "auto_rotate": false
}
```

```
{
  "key_alias": "alias-asdf",           // 🔑 密钥别名 - 用户自定义的密钥名称，便于识别管理
  "duration": "30d",                    // ⏰ 密钥有效期 - 支持 30s/30m/30h/30d/30w/30mo 格式，设置密钥多久后过期，（不填写则永不过期，带测试）
  "models": [],                        // 🤖 允许使用的模型列表 - 空数组表示允许所有模型，也可指定如 ["gpt-4", "claude-3"]
  "spend": 0,                          // 💰 已消费金额 - 密钥当前已使用的金额（单位：美元）
  "max_budget": null,                   // 💸 最大预算 - 密钥最多可消费的总金额，null 表示无限制
  "user_id": null,                      // 👤 用户ID - 关联到特定用户的ID，用于跟踪该用户的所有密钥消费
  "team_id": null,                      // 👥 团队ID - 关联到特定团队的ID，密钥使用团队的限额和设置
  "agent_id": null,                     // 🤖 代理ID - 用于追踪通过代理发出的请求
  "max_parallel_requests": null,       // 🔄 最大并发请求数 - 密钥允许的最大并发请求数，null 表示无限制
  "metadata": {},                      // 📋 元数据 - 自定义的键值对数据，用于存储额外信息如 {"source": "api"}
  "tpm_limit": null,                   // 📊 TPM限制 - 每分钟令牌数限制（Tokens Per Minute），null 表示无限制
  "rpm_limit": null,                   // 📊 RPM限制 - 每分钟请求数限制（Requests Per Minute），null 表示无限制
  "budget_duration": "30d",            // 📅 预算重置周期 - 预算多久重置一次，如 "30d" 表示每月重置
  "allowed_cache_controls": [],       // 💾 允许的缓存控制 - 允许使用的缓存策略，如 ["no-cache", "no-store"]
  "config": {},                        // ⚙️ 配置 - 密钥的特定配置，如自定义模型参数
  "permissions": {},                   // 🔐 权限设置 - 控制密钥的特殊权限，如 {"allow_pii_controls": true}
  "blocked": false,                    // 🚫 阻止状态 - true 表示禁用此密钥，false 表示正常启用
  "aliases": {},                       // 🔗 模型别名 - 模型的昵称映射，如 {"gpt4": "gpt-4-turbo"}
  "budget_id": null,                   // 💳 预算ID - 关联已创建的预算ID，使用预定义的预算设置
  "allowed_routes": [],                // 🛣️ 允许的路由 - 限制密钥可访问的API路由，如 ["chat/completions"]
  "key_type": "default",                // 📌 密钥类型 - default/llm_api/management/read_only，决定密钥的访问权限范围
  "auto_rotate": false                 // 🔄 自动轮换 - true 表示启用自动轮换，false 表示手动管理
}
```

响应：

```
{
  "key": "sk-2rJ...kQF",                   // 🔑 生成的密钥 - 新创建的明文 API Key
  "key_name": "sk-2rJ...kQF",              // 🔖 密钥缩略名 - 用于 UI 展示的缩写形式
  "expires": "2026-05-01T12:34:56.789012+00:00",  // ⏰ 过期时间 - 由 duration 计算得出，null 表示永不过期
  "token_id": "f2f1f9e6...",               // 🆔 Token ID - 数据库中 token 的哈希标识
  "created_at": "2026-04-01T12:00:00.789012+00:00",  // 📅 创建时间
  "updated_at": "2026-04-01T12:00:00.789012+00:00",  // 🔄 更新时间
  "created_by": "admin_user_id",          // 👤 创建者 - 调用该接口的用户 ID
  "updated_by": "admin_user_id",           // 👤 更新者 - 最后修改该密钥的用户 ID

  "key_alias": "alias-asdf",               // 🔑 密钥别名 - 用户自定义的密钥名称
  "duration": "30d",                       // ⏰ 有效期 - 请求时传入，30d 表示 30 天后过期
  "models": [],                            // 🤖 允许模型列表 - 空数组表示允许所有模型
  "spend": 0,                              // 💰 已消费金额 - 密钥当前已使用的金额（美元）
  "max_budget": null,                      // 💸 最大预算 - 密钥最多可消费的总金额，null 表示无限制
  "user_id": null,                         // 👤 用户ID - 关联到特定用户的 ID
  "team_id": null,                         // 👥 团队ID - 关联到特定团队的 ID
  "agent_id": null,                        // 🤖 代理ID - 用于追踪通过代理发出的请求
  "max_parallel_requests": null,           // 🔄 最大并发请求数 - null 表示无限制
  "metadata": {},                          // 📋 元数据 - 自定义的键值对数据
  "tpm_limit": null,                       // 📊 TPM限制 - 每分钟令牌数限制，null 表示无限制
  "rpm_limit": null,                       // 📊 RPM限制 - 每分钟请求数限制，null 表示无限制
  "budget_duration": "30d",                // 📅 预算重置周期 - "30d" 表示每月重置
  "allowed_cache_controls": [],            // 💾 允许的缓存控制 - 允许使用的缓存策略
  "config": {},                            // ⚙️ 配置 - 密钥的特定配置
  "permissions": {},                       // 🔐 权限设置 - 控制密钥的特殊权限
  "blocked": false,                        // 🚫 阻止状态 - true 表示禁用此密钥
  "aliases": {},                           // 🔗 模型别名 - 模型的昵称映射
  "budget_id": null,                       // 💳 预算ID - 关联已创建的预算 ID
  "allowed_routes": [],                   // 🛣️ 允许的路由 - 限制密钥可访问的 API 路由
  "key_type": "default",                   // 📌 密钥类型 - default/llm_api/management/read_only
  "auto_rotate": false,                    // 🔄 自动轮换 - true 表示启用自动轮换
  "rotation_interval": null,                // 🔄 轮换间隔 - auto_rotate=true 时需要，如 "30d"
  "organization_id": null,                 // 🏢 组织ID - 关联到特定组织
  "project_id": null,                      // 📁 项目ID - 关联到特定项目

  "soft_budget": null,                     // �软预算 - 触发警告的阈值，null 表示未设置
  "send_invite_email": null,               // 📧 发送邀请邮件 - 是否发送邀请邮件给用户
  "key_type": "default",                   // 📌 密钥类型 - 决定密钥的访问权限范围
  "router_settings": {},                   // ⚙️ 路由器设置 - 密钥特定的路由器配置
  "access_group_ids": [],                 // 👥 访问组ID列表 - 定义密钥可访问的模型

  "litellm_budget_table": null             // 💳 预算表对象 - 若绑定预算，返回相应对象
}
```

接口路径：key/list
请求体详情可以看swagger，主要参数有:
page,size,user\_id,team\_id,key\_alias,sort\_by,sort\_order,status

成功的响应：

```
{
  "keys": [
    "f5c216b9dd1db35e1b2bcc594a0ee55f4f76742b052c3bfeedcbb7294cc50103",
    "d788174a044904fc5a68f59a740fcf51a96a19d1c7e12d9396e08fb43429f020",
    "415230dbce2afc2654f3238c5cced9725bd8f781e5f64837c32fa73292e4e910"
  ],
  "total_count": 3,
  "current_page": 1,
  "total_pages": 1
}
```

<br />

<br />

接口路径：/key/update 更改key信息

<br />

请求体:

```
{
  "key": "string",                               // ⭐ 必须：要更新的 API Key。可以是明文 sk-xxx 格式，也可以是 token_id (SHA-256 哈希值)。
                                                 //    例如: "sk-abcde12345" 或 "d788174a044904fc5a68f59a740fcf51a96a19d1c7e12d9396e08fb43429f020"

  "key_alias": "string",                         // 🔑 密钥别名 - 用户自定义的密钥名称，便于识别管理。
  "duration": "string",                          // ⏰ 密钥有效期 - 支持 30s/30m/30h/30d/30w/30mo 格式，设置密钥多久后过期。null 表示永不失效。
  "models": [],                                  // 🤖 允许使用的模型列表 - 空数组表示允许所有模型，也可指定如 ["gpt-4", "claude-3"]。
  "spend": 0,                                    // 💰 已消费金额 - 密钥当前已使用的金额（单位：美元）。通常由系统自动更新，手动设置需谨慎。
  "max_budget": 0,                               // 💸 最大预算 - 密钥最多可消费的总金额，null 表示无限制。
  "user_id": "string",                           // 👤 用户ID - 关联到特定用户的ID，用于跟踪该用户的所有密钥消费。
  "team_id": "string",                           // 👥 团队ID - 关联到特定团队的ID，密钥使用团队的限额和设置。
  "agent_id": "string",                          // 🤖 代理ID - 用于追踪通过代理发出的请求。
  "max_parallel_requests": 0,                    // 🔄 最大并发请求数 - 密钥允许的最大并发请求数，null 表示无限制。
  "metadata": {                                  // 📋 元数据 - 自定义的键值对数据，用于存储额外信息如 {"source": "api"}。
    "additionalProp1": {}                        // 示例：{"source": "dashboard-update", "department": "engineering"}
  },
  "tpm_limit": 0,                                // 📊 TPM限制 - 每分钟令牌数限制（Tokens Per Minute），null 表示无限制。
  "rpm_limit": 0,                                // 📊 RPM限制 - 每分钟请求数限制（Requests Per Minute），null 表示无限制。
  "budget_duration": "string",                   // 📅 预算重置周期 - 预算多久重置一次，如 "30d" 表示每月重置。
  "allowed_cache_controls": [],                  // 💾 允许的缓存控制 - 允许使用的缓存策略，如 ["no-cache", "no-store"]。
  "config": {},                                  // ⚙️ 配置 - 密钥的特定配置，如自定义模型参数。
  "permissions": {},                             // 🔐 权限设置 - 控制密钥的特殊权限，如 {"allow_pii_controls": true}。
  "model_max_budget": {},                        // 💸 模型最大预算 - 模型级别的预算限制，例如 {"gpt-4": 20.0, "claude-3-opus": 30.0}。
  "model_rpm_limit": {                           // 📊 模型 RPM 限制 - 模型级别的每分钟请求数限制，例如 {"gpt-4": 50, "claude-3-opus": 30}。
    "additionalProp1": {}
  },
  "model_tpm_limit": {                           // 📊 模型 TPM 限制 - 模型级别的每分钟令牌数限制，例如 {"gpt-4": 50000, "claude-3-opus": 30000}。
    "additionalProp1": {}
  },
  "guardrails": [                                // 🔒 护栏 - 启用的护栏列表，例如 ["content-moderation"]。
    "string"                                     // ⭐ **企业版功能**
  ],
  "policies": [                                  // 📜 策略 - 应用的策略列表，例如 ["default-policy"]。
    "string"                                     // ⭐ **企业版功能**
  ],
  "prompts": [                                   // 💬 提示词 - 允许使用的提示词列表，例如 ["welcome-prompt"]。
    "string"
  ],
  "blocked": true,                               // 🚫 阻止状态 - true 表示禁用此密钥，false 表示正常启用。
  "aliases": {},                                 // 🔗 模型别名 - 模型的昵称映射，如 {"gpt4": "gpt-4-turbo"}。
  "object_permission": {                         // 🔐 对象权限 - 密钥特定的对象访问权限。
    "mcp_servers": [                             // 允许访问的 MCP 服务器列表。
      "string"                                   // ⭐ **企业版功能**
    ],
    "mcp_access_groups": [                       // 允许访问的 MCP 访问组列表。
      "string"                                   // ⭐ **企业版功能**
    ],
    "mcp_tool_permissions": {                    // 允许访问的 MCP 工具权限。
      "additionalProp1": [                       // ⭐ **企业版功能**
        "string"
      ],
      "additionalProp2": [
        "string"
      ],
      "additionalProp3": [
        "string"
      ]
    },
    "vector_stores": [                           // 允许访问的向量存储列表。
      "string"                                   // ⭐ **企业版功能**
    ],
    "agents": [                                  // 允许访问的代理列表。
      "string"                                   // ⭐ **企业版功能**
    ],
    "agent_access_groups": [                     // 允许访问的代理访问组列表。
      "string"                                   // ⭐ **企业版功能**
    ],
    "models": [                                  // 允许访问的模型列表（在 object_permission 内部）。
      "string"                                   // ⭐ **企业版功能**
    ]
  },
  "budget_id": "string",                         // 💳 预算ID - 关联已创建的预算ID，使用预定义的预算设置。
  "tags": [                                      // 🏷️ 标签 - 用于跟踪消费和/或基于标签的路由。
    "string"                                     // ⭐ **企业版功能**
  ],
  "enforced_params": [                           // 强制参数 - 密钥的强制参数列表。
    "string"                                     // ⭐ **企业版功能**
  ],
  "allowed_routes": [],                          // 🛣️ 允许的路由 - 限制密钥可访问的 API 路由，如 ["/chat/completions"]。
  "allowed_passthrough_routes": [                // 🛣️ 允许的透传路由 - 限制密钥可访问的透传 API 路由。
    "string"
  ],
  "allowed_vector_store_indexes": [              // 允许的向量存储索引列表。
    {                                            // ⭐ **企业版功能**
      "index_name": "string",
      "index_permissions": [
        "read"
      ]
    }
  ],
  "rpm_limit_type": "guaranteed_throughput",     // 📊 RPM 限制类型 - "best_effort_throughput", "guaranteed_throughput", "dynamic"。
  "tpm_limit_type": "guaranteed_throughput",     // 📊 TPM 限制类型 - "best_effort_throughput", "guaranteed_throughput", "dynamic"。
  "router_settings": {                           // ⚙️ 路由器设置 - 密钥特定的路由器配置。
    "routing_strategy_args": {                   // ⭐ **企业版功能**
      "additionalProp1": {}
    },
    "routing_strategy": "string",                // ⭐ **企业版功能**
    "model_group_retry_policy": {                // ⭐ **企业版功能**
      "additionalProp1": {}
    },
    "model_group_affinity_config": {             // ⭐ **企业版功能**
      "additionalProp1": [
        "string"
      ],
      "additionalProp2": [
        "string"
      ],
      "additionalProp3": [
        "string"
      ]
    },
    "allowed_fails": 0,                          // ⭐ **企业版功能**
    "cooldown_time": 0,                          // ⭐ **企业版功能**
    "num_retries": 0,                            // ⭐ **企业版功能**
    "timeout": 0,                                // ⭐ **企业版功能**
    "max_retries": 0,                            // ⭐ **企业版功能**
    "retry_after": 0,                            // ⭐ **企业版功能**
    "fallbacks": [                               // ⭐ **企业版功能**
      {
        "additionalProp1": {}
      }
    ],
    "context_window_fallbacks": [                // ⭐ **企业版功能**
      {
        "additionalProp1": {}
      }
    ],
    "model_group_alias": {}                      // ⭐ **企业版功能**
  },
  "access_group_ids": [                          // 👥 访问组ID列表 - 定义密钥可以访问的模型。
    "string"                                     // ⭐ **企业版功能**
  ],
  "temp_budget_increase": 0,                     // 📈 临时预算增加 - 临时增加 Key 的预算。
  "temp_budget_expiry": "2026-04-01T07:41:45.233Z", // 📅 临时预算过期时间 - 临时预算的过期时间 (ISO 8601 格式)。
                                                 // ⭐ **企业版功能**
  "auto_rotate": true,                           // 🔄 自动轮换 - true 表示启用自动轮换。
  "rotation_interval": "string",                 // 🔄 轮换间隔 - 自动轮换的频率（例如 '30d', '90d'）。当 auto_rotate=true 时必须设置。
  "organization_id": "string"                    // 🏢 组织ID - 关联到特定组织。
}
```

<br />

成功过的请求参数

```



{
  "key": "d788174a044904fc5a68f59a740fcf51a96a19d1c7e12d9396e08fb43429f020",
  "token": "d788174a044904fc5a68f59a740fcf51a96a19d1c7e12d9396e08fb43429f020",
  "key_name": "sk-...wNVQ",
  "key_alias": "my-new-alias-for-user",
  "soft_budget_cooldown": false,
  "spend": 0,
  "expires": null,
  "models": [],
  "aliases": {},
  "config": {},
  "router_settings": {},
  "user_id": "2039160439543857153",
  "team_id": "2039160439543857152",
  "agent_id": null,
  "project_id": null,
  "permissions": {},
  "max_parallel_requests": null,
  "metadata": {},
  "blocked": true,
  "tpm_limit": null,
  "rpm_limit": null,
  "max_budget": 100,
  "budget_duration": null,
  "budget_reset_at": null,
  "allowed_cache_controls": [],
  "allowed_routes": [],
  "policies": [],
  "access_group_ids": [],
  "model_spend": {},
  "model_max_budget": {},
  "budget_id": null,
  "organization_id": null,
  "object_permission_id": null,
  "created_at": "2026-04-01T06:46:14.032000+00:00",
  "created_by": "default_user_id",
  "updated_at": "2026-04-01T08:09:31.472000+00:00",
  "updated_by": "default_user_id",
  "last_active": null,
  "rotation_count": 0,
  "auto_rotate": false,
  "rotation_interval": null,
  "last_rotation_at": null,
  "key_rotation_at": null,
  "litellm_budget_table": null,
  "litellm_organization_table": null,
  "litellm_project_table": null,
  "object_permission": null,
  "jwt_key_mappings": null
}
```

<br />

成功的响应：

```



{
  "key": "d788174a044904fc5a68f59a740fcf51a96a19d1c7e12d9396e08fb43429f020",
  "token": "d788174a044904fc5a68f59a740fcf51a96a19d1c7e12d9396e08fb43429f020",
  "key_name": "sk-...wNVQ",
  "key_alias": "my-new-alias-for-user",
  "soft_budget_cooldown": false,
  "spend": 0,
  "expires": null,
  "models": [],
  "aliases": {},
  "config": {},
  "router_settings": {},
  "user_id": "2039160439543857153",
  "team_id": "2039160439543857152",
  "agent_id": null,
  "project_id": null,
  "permissions": {},
  "max_parallel_requests": null,
  "metadata": {},
  "blocked": true,
  "tpm_limit": null,
  "rpm_limit": null,
  "max_budget": 100,
  "budget_duration": null,
  "budget_reset_at": null,
  "allowed_cache_controls": [],
  "allowed_routes": [],
  "policies": [],
  "access_group_ids": [],
  "model_spend": {},
  "model_max_budget": {},
  "budget_id": null,
  "organization_id": null,
  "object_permission_id": null,
  "created_at": "2026-04-01T06:46:14.032000+00:00",
  "created_by": "default_user_id",
  "updated_at": "2026-04-01T08:09:31.472000+00:00",
  "updated_by": "default_user_id",
  "last_active": null,
  "rotation_count": 0,
  "auto_rotate": false,
  "rotation_interval": null,
  "last_rotation_at": null,
  "key_rotation_at": null,
  "litellm_budget_table": null,
  "litellm_organization_table": null,
  "litellm_project_table": null,
  "object_permission": null,
  "jwt_key_mappings": null
}
```

<br />

## 2.1.3. Internal User management

接口路径：/user/new
创建内部用户

参数说明

```
{
  "user_id": "string",                           // 👤 用户ID - 可选。指定唯一的用户ID，若不传则系统自动生成 UUID。
  "user_alias": "string",                        // 👤 用户别名 - 可选。为该用户设置一个易于识别的名称。
  "user_email": "string",                        // 📧 用户邮箱 - 可选。关联用户的电子邮件地址。
  "user_role": "string",                         // 🔐 用户角色 - 可选。可选值: "proxy_admin", "proxy_admin_viewer", "internal_user", "internal_user_viewer"。
  "teams": [],                                   // 👥 团队列表 - 可选。用户所属的团队 ID 列表，或包含角色和预算的团队对象列表。
  "organizations": [],                           // 🏢 组织列表 - 可选。用户所属的组织 ID 列表。
  
  "key_alias": "string",                         // 🔑 密钥别名 - 为 /user/new 自动生成的 API Key 设置别名。
  "duration": "string",                          // ⏰ 密钥有效期 - 支持 30s/30m/30h/30d/30w/30mo 格式，设置自动生成的密钥多久后过期。
  "models": [],                                  // 🤖 允许使用的模型列表 - 空数组表示允许所有模型，也可指定如 ["gpt-4", "claude-3"]。
  
  "max_budget": 0.0,                             // 💸 最大总预算 - 该用户（跨所有密钥）允许消费的总金额（单位：美元）。
  "budget_duration": "string",                   // 📅 预算重置周期 - 预算多久重置一次，如 "30d" 表示每月重置。
  "soft_budget": 0.0,                            // 🔔 软预算 - 达到此金额时触发警告，但不会拦截请求。
  
  "tpm_limit": 0,                                // 📊 TPM限制 - 每分钟令牌数限制（Tokens Per Minute）。
  "rpm_limit": 0,                                // 📊 RPM限制 - 每分钟请求数限制（Requests Per Minute）。
  "max_parallel_requests": 0,                    // 🔄 最大并发请求数 - 允许的最大并发请求数。
  
  "model_max_budget": {},                        // 💸 模型最大预算 - 模型级别的预算限制，例如 {"gpt-4": 20.0}。
  "model_rpm_limit": {},                         // 📊 模型 RPM 限制 - 模型级别的每分钟请求数限制，例如 {"gpt-4": 50}。
  "model_tpm_limit": {},                         // 📊 模型 TPM 限制 - 模型级别的每分钟令牌数限制，例如 {"gpt-4": 50000}。
  
  "metadata": {},                                // 📋 元数据 - 自定义的键值对数据，用于存储额外信息。
  "config": {},                                  // ⚙️ 配置 - 用户特定的特定配置（已弃用，建议使用 metadata）。
  
  "guardrails": [],                              // 🔒 护栏 - 启用的护栏列表，例如 ["content-moderation"]。⭐ **企业版功能**
  "policies": [],                                // 📜 策略 - 应用的策略列表。⭐ **企业版功能**
  "prompts": [],                                 // 💬 提示词 - 允许使用的提示词 ID 列表。
  
  "object_permission": {                         // 🔐 对象权限 - 用户特定的对象访问权限。
    "mcp_servers": ["string"],                   // 允许访问的 MCP 服务器列表。⭐ **企业版功能**
    "vector_stores": ["string"],                 // 允许访问的向量存储列表。⭐ **企业版功能**
    "agents": ["string"]                         // 允许访问的代理列表。⭐ **企业版功能**
  },
  
  "auto_create_key": true,                       // 🔄 自动创建密钥 - 默认为 true，即在创建用户时自动返回一个 API Key。
  "send_invite_email": false,                    // 📧 发送邀请邮件 - 是否向 user_email 发送邀请邮件。
  "blocked": false,                              // 🚫 阻止状态 - true 表示禁用此用户。
  "aliases": {},                                 // 🔗 模型别名 - 模型的昵称映射，如 {"gpt4": "gpt-4-turbo"}。
  "allowed_cache_controls": [],                  // 💾 允许的缓存控制 - 允许使用的缓存策略。
  "sso_user_id": "string"                        // 🆔 SSO用户ID - 关联单点登录系统的用户 ID。
}
```

<br />

<br />

## 2.1.4. team management

接口路径：/team/new
创建团队

我们可以看到 litellm/proxy/management\_endpoints/ 目录下有非常多文件，比如：
key\_management\_endpoints.py (管理 API Key 的接口，例如 /key/generate)
team\_endpoints.py (管理团队的接口，例如 /team/new)
internal\_user\_endpoints.py (管理内部用户的接口，例如 /user/new)
budget\_management\_endpoints.py (管理预算的接口)
model\_management\_endpoints.py (管理可用模型的接口)

<br />

<br />

<br />

<br />

## 2.1.5.日志：

<br />

接口路径：/spend/logs 详细聊天日志

<br />

<br />

<br />

/team/daily/activity (GET 请求) 聚合聊天日志

- 端点 : /team/daily/activity
- 文件 : litellm/proxy/management\_endpoints/team\_endpoints.py
- team\_ids : 支持 。它接受一个逗号分隔的 team\_ids 字符串。如果未提供，则返回所有团队的数据。
- start\_date , end\_date : 支持 ，用于指定活动期间。
- model : 支持 ，按模型名称过滤。
- api\_key : 支持 。它接受一个可选的 api\_key 参数。如果用户不是管理员且未提供 api\_key ，它将按用户自己的 API 密钥进行过滤。






# 模型列表

有，而且基本已经满足你说的“两种模型名”，只是分布在两个接口里。
GET /v1/models（或 /models）
返回的是用户可用/可调用的模型名（即路由层模型组名，通常就是你配置给用户用的名字，比如 claude-4-6-sonnet）。
GET /v1/model/info（或 /model/info）
返回更完整的 deployment 信息，包含：
model_name：用户使用时的模型名
litellm_params.model：底层实际模型名（比如 bedrock/...、anthropic/... 这种）