# MCP 与 Agent 完整权限指南

本版本允许经过认证的 MCP 或 Agent 客户端读取和修改当前用户拥有的财务数据。所谓“完整权限”是指覆盖账户、余额、交易、分类、标签和标签组的创建、查询、更新与删除；它不绕过登录、Token 类型、用户归属、字段校验、转账关联或被引用数据的删除保护。

## 服务端开关

```dotenv
EBK_MCP_ENABLE_MCP=true
EBK_MCP_MCP_ALLOWED_REMOTE_IPS=
EBK_SECURITY_ENABLE_API_TOKEN=true
EBK_SECURITY_API_TOKEN_ALLOWED_REMOTE_IPS=
```

MCP 地址为 `/mcp`，要求专用 MCP Bearer Token。Agent Skill 的 REST 命令使用用户 API Token。两种 Token 都应设置有效期、独立保存并定期轮换；不要使用普通登录 Token 代替。

## MCP 工具范围

| 领域 | 查询 | 写操作 |
| --- | --- | --- |
| 账户 | `query_all_accounts`、`query_all_accounts_balance` | `create_account`、`update_account`、`adjust_account_balance`、`delete_account` |
| 交易 | `query_transactions` | `add_transaction`、`update_transaction`、`delete_transaction` |
| 分类 | `query_all_transaction_categories` | `create_transaction_category`、`update_transaction_category`、`delete_transaction_category` |
| 标签 | `query_all_transaction_tags` | 标签与标签组的 create/update/delete |
| 汇率 | `query_latest_exchange_rates` | 只读 |

更新和删除前先查询稳定 ID。支持 `dry_run` 的工具应先预演，再由用户确认正式提交。直接修正账户当前余额使用 `update_account.balance`；只有明确需要生成可审计的余额修改交易时才使用 `adjust_account_balance`。`update_account` 支持账户全部可变字段，`update_transaction` 支持时间、类型、分类、账户、金额、标签、图片、备注、隐藏金额和地理位置等全部可变字段；未传字段保持不变。

## Agent Skill

技能目录：

```text
skills/ezbookkeeping/
├── SKILL.md
└── scripts/
    ├── ebktools.ps1
    └── ebktools.sh
```

配置：

```dotenv
EBKTOOL_SERVER_BASEURL=https://book.example.com
EBKTOOL_TOKEN=你的用户API令牌
```

列出命令：

```bash
sh skills/ezbookkeeping/scripts/ebktools.sh list
```

```powershell
skills\ezbookkeeping\scripts\ebktools.ps1 list
```

复杂更新可用 JSON 文件原样传递官方 REST 请求体：

```bash
sh skills/ezbookkeeping/scripts/ebktools.sh \
  --body-file ./request.json transactions-modify
```

```powershell
skills\ezbookkeeping\scripts\ebktools.ps1 `
  -bodyFile .\request.json transactions-modify
```

## 安全建议

- 生产环境启用 HTTPS，并设置唯一的 `EBK_SECURITY_SECRET_KEY`。
- 用 IP/CIDR 白名单限制 MCP 与 API Token 来源。
- 为自动化单独创建 Token，不与浏览器会话共享。
- 批量删除、跨账户转账和大额变更先使用 `dry_run`。
- 定期备份数据库与附件；下一阶段建议增加逐条审计日志和一键回滚。
