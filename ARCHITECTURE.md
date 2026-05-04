# Agent Ontogenesis Framework (AOF)

## 1. 系统定位

AOF 是一个让通用Agent发生**形态发生**（Morphogenesis）的框架。它不依赖Prompt工程、外挂RAG或硬编码工作流，而是通过一套声明式协议驱动Agent重组自身的认知结构、能力边界和交互模式，实现从通用到专用的真正分化。

---

## 2. 核心概念

| 概念 | 定义 | 类比 |
|------|------|------|
| **Morphogen** | 声明式专用化协议，定义领域本体、能力清单和演化规则 | 干细胞分化指令 |
| **Crystal** | 可内化的能力单元，Agent持有它而非检索它 | 细胞器 |
| **Germline** | 可跨实例遗传的能力基因库 | 生殖系DNA |
| **Soma** | 单次运行时的临时特化实例 | 体细胞 |

---

## 3. 系统架构

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              AOF 系统架构                                         │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  Layer 5: 用户定义层                                                               │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Morphogen Specification                                                │   │
│  │  - 领域本体 (Ontology)                                                   │   │
│  │  - 能力晶体声明 (Capability Declaration)                                  │   │
│  │  - 演化规则 (Evolution Rules)                                             │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                                    ▼                                            │
│  Layer 4: 编译与合成层                                                             │
│  ┌─────────────────────────────┐  ┌─────────────────────────────────────────┐   │
│  │  Morphogen Compiler         │  │  Crystal Forge                          │   │
│  │  ├─ Parser                  │  │  ├─ Loader                              │   │
│  │  ├─ Validator               │──│  ├─ Synthesizer                         │   │
│  │  └─ Emitter                 │  │  └─ Binder                              │   │
│  └─────────────────────────────┘  └─────────────────────────────────────────┘   │
│                                    │                                            │
│                                    ▼ CrystalConfig[]                            │
│  Layer 3: 分化层                                                                   │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Agent Differentiator                                                   │   │
│  │  ├─ Injector        (注入晶体配置)                                        │   │
│  │  ├─ Activator       (激活能力边界)                                        │   │
│  │  └─ Belief Calibrator (校准置信度边界)                                    │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                                    ▼ SpecializedAgent                           │
│  Layer 2: 运行时层                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │  Ontogeny Runtime                                                       │   │
│  │  ├─ Cortex          (认知皮层：上下文构建与晶体选择)                        │   │
│  │  ├─ Actuator        (执行驱动：工具调用与外部交互)                         │   │
│  │  └─ Reflector       (反思引擎：质量评估与经验提炼)                         │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│                    ┌───────────────┴───────────────┐                           │
│                    ▼                               ▼                           │
│  Layer 1: 持久化层                                                                 │
│  ┌──────────────────────────────┐  ┌──────────────────────────────────────────┐ │
│  │  Crystal Repository          │  │  Germline Archive                        │ │
│  │  (运行时晶体缓存与动态管理)    │  │  (可遗传进化记录与版本管理)                 │ │
│  └──────────────────────────────┘  └──────────────────────────────────────────┘ │
│                                                                                 │
│  Layer 0: 基座层                                                                   │
│  ┌────────────────────────┐  ┌──────────────────┐  ┌──────────────────────────┐ │
│  │  LLM Kernel            │  │  Meta-Cognition  │  │  Sandbox                 │ │
│  │  (大模型统一接口)        │  │  (元认知能力)      │  │  (安全执行环境)            │ │
│  └────────────────────────┘  └──────────────────┘  └──────────────────────────┘ │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. 组件详解

### 4.1 Morphogen Specification（形态发生素协议）

**职责**：用户声明专用化需求的唯一入口。

**核心结构**：

```yaml
morphogen:
  identity:
    domain: string           # 领域标识
    ontology: string         # 领域本体文件路径
    
  capabilities:
    - id: string
      type: enum             # knowledge | tool | protocol | memory | belief
      source: object         # 原材料定义
      activation: enum       # always | contextual | manual
      
  evolution:
    rules:
      - trigger: string      # 条件表达式
        action: string       # 执行动作
    inherit: boolean         # 是否允许遗传
```

**设计原则**：
- 描述"能力结构"而非"行为表现"
- 声明式而非命令式
- 可版本化管理

---

### 4.2 Morphogen Compiler（形态发生素编译器）

**职责**：将声明式协议转化为系统可执行配置。

| 子组件 | 输入 | 输出 | 职责 |
|-------|------|------|------|
| **Parser** | Morphogen YAML | AST | 解析DSL语法，生成抽象语法树 |
| **Validator** | AST | 验证报告 | 检查引用一致性、权限合法性、晶体间冲突 |
| **Emitter** | 验证后的AST | `CrystalConfig[]` | 输出标准化晶体配置清单 |

**关键输出**：

```typescript
interface CrystalConfig {
  id: string;
  type: 'knowledge' | 'tool' | 'protocol' | 'memory' | 'belief';
  source: SourceConfig;      // 原材料位置与加载方式
  synthesis: SynthesisConfig; // 合成策略
  activation: ActivationPolicy; // 激活策略
  dependencies: string[];     // 依赖的其他晶体
}
```

---

### 4.3 Crystal Forge（能力晶体熔炉）

**职责**：将原材料加工为Agent可持有的运行时晶体。

| 子组件 | 输入 | 输出 | 职责 |
|-------|------|------|------|
| **Loader** | `CrystalConfig` + 外部资源 | 原材料包 | 加载知识图谱、工具定义、语法文件、Schema等 |
| **Synthesizer** | 原材料包 | 结构化晶体 | 构建推理索引、生成语义表征、编译交互策略 |
| **Binder** | 结构化晶体 | `RuntimeCrystal` | 绑定执行语义：指定如何被调用、如何与LLM上下文结合 |

**五类晶体加工方式**：

| 类型 | 原材料 | 加工产物 | Binder绑定方式 |
|------|-------|---------|--------------|
| Knowledge | 文本/图谱/数据库 | 程序性知识图谱 + 推理链索引 | 上下文注入 + 推理引导 |
| Tool | OpenAPI/函数定义 | 工具语义表征 + 参数模板 | Function Call绑定 |
| Protocol | BNF/对话样本 | 自适应交互策略网络 | 对话管理器绑定 |
| Memory | Schema定义 | Episodic存储结构 | 记忆检索器绑定 |
| Belief | 能力边界定义 | 概率化信念图 | 元认知过滤器绑定 |

---

### 4.4 Agent Differentiator（Agent分化器）

**职责**：将通用Agent基座转化为专用Agent实例。

| 子组件 | 输入 | 输出 | 职责 |
|-------|------|------|------|
| **Injector** | Base Agent + `CrystalConfig[]` | 待激活Agent | 将晶体配置注入Agent元认知层 |
| **Activator** | 待激活Agent + `RuntimeCrystal[]` | 激活态Agent | 按策略标记初始可用晶体 |
| **Belief Calibrator** | 激活态Agent + Belief Crystal | 校准后Agent | 设定能力边界：自主/请示/拒绝 |

**分化公式**：

```
SpecializedAgent = BaseKernel + Σ(Activated Crystals) + BeliefBoundary + EpisodicMemory
```

**关键原则**：分化是"叠加"而非"覆盖"。通用能力被专用晶体"遮蔽"或"增强"，不会被删除。

---

### 4.5 Ontogeny Runtime（个体发生运行时）

**职责**：专用Agent的实际运行与演化环境。

| 子组件 | 输入 | 输出 | 职责 |
|-------|------|------|------|
| **Cortex** | 用户输入 + 激活晶体 + 信念 + 记忆 | 推理上下文 | 构建认知皮层：选择进入LLM上下文的信息 |
| **Actuator** | LLM输出 + Tool Crystal | 执行结果 | 驱动工具调用、外部交互、状态变更 |
| **Reflector** | 输入/输出/执行轨迹 | 反思报告 | 评估推理质量、工具使用、交互策略 |

**单次交互数据流**：

```
用户输入
    │
    ▼
┌─────────┐    ┌──────────┐    ┌─────────┐    ┌──────────┐
│  Cortex │───►│ LLM Kernel│───►│ Actuator│───►│ 外部系统  │
│(构建上下文)│    │ (推理)   │    │ (执行)   │    │ (工具/API)│
└────┬────┘    └──────────┘    └────┬────┘    └────┬─────┘
     │                               │              │
     │                               └──────────────┘
     │                                  执行结果
     │
     ▼
┌───────────┐    ┌────────────────┐    ┌──────────────────┐
│ Reflector │───►│ Crystal Repository│   │ Germline Archive │
│ (反思提炼)  │    │ (更新晶体状态)    │    │ (选择性遗传)     │
└───────────┘    └────────────────┘    └──────────────────┘
```

---

### 4.6 Crystal Repository（晶体仓库）

**职责**：运行时晶体的缓存与动态管理。

**核心功能**：
- 按Agent实例隔离存储
- 支持动态加载/卸载（运行中根据新任务激活新晶体）
- 支持晶体融合（多颗晶体组合为临时复合能力）
- 跟踪晶体使用频率和效果，为Cortex的选择提供数据

**接口**：

```typescript
interface CrystalRepository {
  load(instanceId: string, crystalId: string): RuntimeCrystal;
  unload(instanceId: string, crystalId: string): void;
  fuse(instanceId: string, crystalIds: string[]): RuntimeCrystal;
  getActive(instanceId: string): RuntimeCrystal[];
  updateStats(instanceId: string, crystalId: string, metrics: Metrics): void;
}
```

---

### 4.7 Germline Archive（生殖系档案）

**职责**：跨会话的能力进化持久化与版本管理。

**核心功能**：
- 存储Reflector评估为"高价值"的新晶体
- 管理晶体版本（同一领域多版本共存）
- 维护进化树（追踪能力从初始配置到当前状态的演进路径）
- 提供遗传接口（新实例化时自动加载Germline精华）

**遗传策略**：

```
Reflector评估 ──► 质量阈值筛选 ──► (可选)人工确认 ──► 写入Germline
```

**接口**：

```typescript
interface GermlineArchive {
  inherit(domain: string, generation: number): CrystalConfig[];
  commit(domain: string, crystal: RuntimeCrystal, lineage: Lineage): void;
  getLineage(domain: string): EvolutionTree;
  rollback(domain: string, version: string): void;
}
```

---

### 4.8 Foundation Layer（基座层）

**职责**：整个系统的底层依赖。

| 组件 | 职责 | 接口 |
|------|------|------|
| **LLM Kernel** | 封装多厂商大模型调用，提供统一complete/chat接口 | `complete(prompt, context) -> response` |
| **Meta-Cognition** | 提供自我监控能力：当前状态、推理路径、置信度 | `monitor(thought) -> meta_report` |
| **Sandbox** | 工具执行的安全沙箱：网络隔离、权限控制、超时熔断 | `execute(tool, params) -> result` |

---

## 5. 数据模型

### 5.1 Morphogen

```yaml
morphogen:
  version: "1.0"
  identity:
    domain: string
    ontology_ref: string
    description: string
  capabilities:
    - id: string
      type: enum           # knowledge | tool | protocol | memory | belief
      config: object
      activation: enum     # always | contextual | manual
      dependencies: [string]
  evolution:
    rules:
      - trigger: string
        action: string
        priority: number
    inherit: boolean
    max_generations: number
```

### 5.2 Crystal

```yaml
crystal:
  id: string
  type: enum
  version: string
  state: enum              # latent | active | fused | deprecated
  payload: object          # 类型决定结构
  activation_policy:
    mode: enum             # always | conditional | triggered
    condition: expression
  meta:
    created_at: timestamp
    source_morphogen: string
    confidence: float
    usage_count: number
    last_used: timestamp
```

### 5.3 Agent Instance

```yaml
agent_instance:
  id: string
  base_kernel_ref: string
  morphogen_ref: string
  crystals:
    - crystal_id: string
      state: enum          # active | latent
      fusion_parent: string
  runtime:
    current_beliefs: object
    episodic_memory: [memory_event]
    turn_count: number
    context_window_usage: number
  lineage:
    generation: number
    parent_instance: string
    mutations: [crystal_id]
```

---

## 6. 组件交互时序

### 6.1 从定义到首次运行

```
用户
 │
 ├── 1.提交 Morphogen Specification
 │
 ▼
Morphogen Compiler
 │
 ├── 2.解析(Parse) ──► 3.验证(Validate) ──► 4.输出(Emit)
 │
 ▼
Crystal Forge
 │
 ├── 5.加载(Load) ──► 6.合成(Synthesize) ──► 7.绑定(Bind)
 │
 ▼
Agent Differentiator
 │
 ├── 8.注入(Inject) ──► 9.激活(Activate) ──► 10.校准(Calibrate)
 │
 ▼
SpecializedAgent (就绪)
 │
 ├── 11.接收用户输入
 │
 ▼
Ontogeny Runtime
 │
 ├── 12.Cortex构建上下文 ──► 13.LLM推理 ──► 14.Actuator执行
 │
 ▼
用户 (收到响应)
 │
 ├── 15.Reflector反思 ──► 16.更新Repository ──► 17.(可选)遗传到Germline
 │
 ▼
系统 (完成一次交互闭环)
```

### 6.2 运行时动态演化

```
用户输入
    │
    ▼
Cortex ──► [评估当前激活晶体是否足够]
    │
    ├── 足够 ──► 继续正常推理
    │
    └── 不足 ──► 向Crystal Repository请求加载新晶体
                    │
                    ▼
              新晶体激活 ──► 重新构建上下文 ──► 继续推理
                    │
                    ▼
              Reflector评估新晶体效果 ──► 更新使用统计
```

---

## 7. 关键设计决策

### 决策1：晶体 vs 插件

| 维度 | 插件 | 晶体 |
|------|------|------|
| 关系 | 外挂，通过API调用 | 内化，Agent持有并用于推理 |
| 影响 | 不改变Agent本身 | 改变Agent的认知结构 |
| 类比 | 计算器 | 心算能力 |

### 决策2：动态激活 vs 全量加载

专用Agent可能持有数十颗晶体，但单次交互仅需激活3-5颗。Cortex根据输入语义动态选择，避免上下文爆炸。

### 决策3：Belief Calibrator 的必要性

没有置信度边界的Agent会幻觉。Belief Crystal让Agent明确知道：
- **自主域**：可以独立决策的任务
- **请示域**：需要用户确认的任务
- **拒绝域**：超出能力边界必须拒绝的任务

### 决策4：Germline 的版本化

Agent演化不总是正向的。版本化支持：
- 回滚到稳定版本
- A/B测试不同进化路径
- 可视化进化树理解能力演进

---

## 8. 与传统方案对比

| 维度 | 传统方案 | AOF |
|------|---------|-----|
| 专用化方式 | System Prompt + RAG + 硬编码工作流 | Morphogen协议 + 晶体内化 + 自适应编排 |
| 知识处理 | 向量检索，运行时注入文本 | 程序性知识图谱，Agent持有并用于推理 |
| 工具使用 | 预定义Function Call | 可学习语义的工具表征 |
| 交互模式 | 预定义DAG/状态机 | 自适应交互策略 |
| 进化能力 | 静态，需人工重新配置 | 运行时自我改进 + 跨会话遗传 |
| 上下文管理 | 全量加载，容易溢出 | 动态激活，按需加载 |

---

## 9. 部署视图

```
┌─────────────────────────────────────────────────────────────┐
│                      用户应用层                               │
│         (业务系统 / Chat界面 / API服务)                        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      AOF 核心层                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Compiler   │  │ Differentiator│ │ Ontogeny Runtime   │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│ Crystal Store │  │ Germline DB   │  │ LLM Provider  │
│ (晶体缓存)     │  │ (进化档案)     │  │ (大模型服务)   │
└───────────────┘  └───────────────┘  └───────────────┘
```

**部署单元**：
- AOF Core：作为库或服务部署
- Crystal Store：Redis/内存（热数据）+ 对象存储（冷数据）
- Germline DB：关系型/图数据库（支持版本和谱系查询）
- LLM Provider：多厂商适配，可切换

---

## 10. 接口总览

### 10.1 用户接口

```python
# 定义专用化协议
morphogen = Morphogen.from_file("medical.yaml")

# 编译
config = Compiler.compile(morphogen)

# 分化（创建专用Agent）
agent = Differentiator.differentiate(base_kernel, config)

# 运行
response = await agent.run("患者发热3天，咳嗽...")

# 手动触发遗传（也可由evolution规则自动触发）
await agent.inherit_to_germline()
```

### 10.2 组件间接口

| 调用方 | 被调用方 | 接口 |
|-------|---------|------|
| Compiler | Forge | `forge.build(config: CrystalConfig) -> RuntimeCrystal` |
| Differentiator | Repository | `repo.register(instanceId, crystals: RuntimeCrystal[])` |
| Runtime | Repository | `repo.activate(instanceId, crystalId) / repo.deactivate(...)` |
| Runtime | LLM Kernel | `llm.complete(context: Context) -> Response` |
| Reflector | Germline | `germline.commit(domain, crystal, lineage)` |
| Reflector | Repository | `repo.updateMetrics(instanceId, crystalId, metrics)` |
