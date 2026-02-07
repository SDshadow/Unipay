

# 架构设计说明（Architecture Design）

## 一、项目背景与目标

在真实业务中，支付系统需要面对以下典型问题：

- 不同支付渠道（支付宝、微信等）**接口形态、参数、签名方式高度不一致**
- 第三方支付系统**不可靠**（超时、重复回调、验签失败）
- 业务系统不希望直接耦合支付协议细节

本项目的目标不是实现一个完整商业支付中台，而是：

> **设计一个“可扩展、可插拔、可测试”的支付执行内核**，
>  用统一的执行模型抽象不同支付渠道与支付方式。

------

## 二、整体架构概览

项目整体采用 **Pipeline + Plugin** 的执行模型，核心由三部分组成：

```
API 层
  ↓
Service（业务用例）
  ↓
Provider（支付渠道）
  ↓
Rocket + Plugin Pipeline（执行内核）
```

### 核心思想

- **Rocket**：支付执行过程的“上下文容器”
- **Plugin**：支付流程中的一个独立步骤
- **Pipeline**：按顺序执行 Plugin，完成一次支付动作

------

## 三、核心概念设计

### 1️⃣ Rocket（执行上下文）

Rocket 是整个支付流程中**唯一的状态载体**，用于在插件之间传递数据。

```
type Rocket struct {
    Params   map[string]any   // 原始业务参数
    Payload  map[string]any   // 支付协议参数（逐步构建）
    Response any              // 第三方返回结果
    Result   Result           // 统一支付结果
}
```

设计原则：

- **Plugin 之间不直接依赖**
- 只通过 Rocket 读写数据
- Rocket 在不同阶段呈现不同“形态”，但结构统一

------

### 2️⃣ Plugin（支付插件）

每个 Plugin 表示支付流程中的**一个明确阶段**：

| Plugin       | 职责               |
| ------------ | ------------------ |
| BuildPlugin  | 构建支付请求参数   |
| SignPlugin   | 对请求进行签名     |
| HttpPlugin   | 发起 HTTP 请求     |
| VerifyPlugin | 校验返回签名       |
| ParsePlugin  | 解析响应为统一结果 |

统一接口：

```
type Plugin interface {
    Handle(r *Rocket) error
}
```

设计特点：

- **单一职责**
- 可自由组合
- 易于测试（可 Mock 任意 Plugin）

------

### 3️⃣ Pipeline（插件执行管道）

Pipeline 负责**按顺序执行 Plugin 链**：

```
Rocket
  ↓
Build → Sign → HTTP → Verify → Parse
  ↓
Result
```

Pipeline 本身不关心业务含义，只保证：

- 顺序执行
- 出错即中断
- Rocket 在插件间安全传递

------

## 四、Provider 设计（支付渠道抽象）

### Provider 的角色

Provider 不是直接“发请求”，而是：

> **负责为某一支付渠道，组装一条合适的 Plugin Pipeline**

```
type Provider interface {
    Pay(ctx context.Context, order Order) (Result, error)
}
```

示例（支付宝）：

```
func (a *AlipayProvider) Pay(order Order) (Result, error) {
    rocket := NewRocket(order)

    pipeline := NewPipeline(
        BuildPlugin{},
        SignPlugin{},
        HttpPlugin{},
        VerifyPlugin{},
        ParsePlugin{},
    )

    return pipeline.Execute(rocket)
}
```

这样做的好处：

- 不同支付方式 = 不同 Plugin 组合
- 新增支付渠道 **不影响核心执行逻辑**
- Rocket 与 Pipeline 完全复用

------

## 五、Service 层设计（业务视角）

Service 层是**业务与支付系统的边界**，负责：

- 参数校验
- 选择 Provider
- 屏蔽支付协议细节

```
type PaymentService struct {
    providerFactory ProviderFactory
}

func (s *PaymentService) Pay(req PayRequest) (PayResult, error) {
    provider := s.providerFactory.Get(req.Channel)
    return provider.Pay(req.Order)
}
```

业务系统只关心：

- 我用哪个渠道
- 我是否支付成功

------

## 六、为什么不用“传统 SDK 直调模式”

### SDK 直调的问题

- 业务系统直接感知支付协议
- 多业务线难以统一治理
- 协议升级成本高

### 本项目的取舍

- 不追求“功能全”
- 优先体现：
  - 抽象能力
  - 执行模型
  - 可扩展性

这是一个**偏“支付内核设计”的工程实践**。

------

## 七、测试策略

- **Plugin 单测**：验证每个插件行为
- **Pipeline 集成测试**：Mock HTTP，验证完整支付流程
- 不依赖真实支付渠道

------

## 八、设计取舍说明

本项目在设计中 **刻意避免**：

- 过度中台化（如多级 orchestrator）
- 与业务无关的通用 util 抽象
- 为“未来可能用到”而设计的复杂结构

目标是：

> **在有限复杂度下，最大化体现系统设计能力**

------

## 九、总结

该项目通过 **Rocket + Plugin + Pipeline** 的模型：

- 统一了不同支付渠道的执行结构
- 降低了支付协议对业务的侵入
- 提供了一个可扩展、可测试的支付执行内核
