LLM调用使用说明

# 1. 01-LLM调用.md

## 1.1. 01-LLM调用.md
本文件是LLM调用目录的介绍文件，详细说明了各类大语言模型的接口调用方式、参数配置、多模态处理（如图像调用）以及不同模型的缓存优化机制。通过规范化的调用指南，帮助开发者高效、稳定地集成多种模型能力。

# 2. 调用核心要素

## 2.1. 基础调用流程
LLM调用通常涉及身份验证、端点配置、请求参数构建和响应处理四个核心步骤。支持同步和流式（Streaming）两种主要调用模式。

## 2.2. 多模态支持
现代模型支持图像、音频等多种输入形式。图像调用通常支持URL直接引用或Base64编码上传，需根据模型限制选择合适的解析分辨率和格式。

# 3. 支持的模型系列

## 3.1. OpenAI 系列
涵盖 GPT-4o, GPT-4 Turbo 等模型。使用标准的 `/v1/chat/completions` 接口，支持流式输出和多模态图像输入。详见 [03-OpenAI系列.md](file:///Users/aliang/Documents/document/knowladge/agent/03-LLM%E8%B0%83%E7%94%A8/03-OpenAI%E7%B3%BB%E5%88%97.md)。

## 3.2. Claude 系列 (Anthropic)
涵盖 Claude 3.5 Sonnet, Claude 3 Opus 等模型。使用 Messages API，系统提示词独立配置，支持显式 Prompt Caching。详见 [04-Claude系列.md](file:///Users/aliang/Documents/document/knowladge/agent/03-LLM%E8%B0%83%E7%94%A8/04-Claude%E7%B3%BB%E5%88%97.md)。

## 3.3. Gemini 系列 (Google)
涵盖 Gemini 1.5 Pro/Flash 等模型。支持超长上下文（最高 2M+ Tokens）和原生多模态，提供 Context Caching 功能。详见 [05-Gemini系列.md](file:///Users/aliang/Documents/document/knowladge/agent/03-LLM%E8%B0%83%E7%94%A8/05-Gemini%E7%B3%BB%E5%88%97.md)。
