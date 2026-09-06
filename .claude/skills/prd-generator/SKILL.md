---
name: prd-generator
description: 通过交互式问答生成PRD需求文档。触发条件：(1) 用户说 /prd (2) 用户说"写需求"、"需求文档"、"生成PRD"等自然语言。工作流程：检查现有文档 -> 交互式收集需求信息 -> 生成并写入docs/需求.md
---

# PRD生成器

## 工作流程

### 1. 检查现有文档
```bash
node .claude/skills/prd-generator/scripts/file-check.js
```

- NOT_EXIST: 直接创建新文档
- EMPTY: 覆盖空文档
- EXISTS: 显示前5行，询问是否覆盖

### 2. 交互式收集信息
使用AskUserQuestion工具一次性收集所有信息：

```javascript
AskUserQuestion({
  questions: [
    {
      question: "需求的一句话描述是什么？",
      header: "需求描述",
      options: [],
      multiSelect: false
    },
    {
      question: "业务背景和目标是什么？",
      header: "业务目标",
      options: [],
      multiSelect: false
    },
    {
      question: "目标用户是谁？",
      header: "目标用户",
      options: [],
      multiSelect: false
    },
    {
      question: "核心功能点有哪些？（请逐个列出）",
      header: "功能点",
      options: [],
      multiSelect: true
    },
    {
      question: "非功能性需求有哪些？",
      header: "非功能需求",
      options: [
        {label: "性能要求", description: "响应时间、并发量等"},
        {label: "安全要求", description: "权限、加密等"},
        {label: "兼容性", description: "浏览器、设备等"}
      ],
      multiSelect: true
    }
  ]
})
```

### 3. 生成文档
使用 [prd-template.md](references/prd-template.md) 格式生成完整文档

### 4. 写入文件
```bash
node .claude/skills/prd-generator/scripts/file-create.js "<文档内容>"
```

输出：docs/需求.md
