# ezBookkeeping Pixel

[![开源协议](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![自动构建](https://github.com/Panda-995/ezbookkeeping-pixel/actions/workflows/publish-container.yml/badge.svg)](https://github.com/Panda-995/ezbookkeeping-pixel/actions/workflows/publish-container.yml)
[![镜像架构](https://img.shields.io/badge/image-amd64%20%7C%20arm64-176b5b)](https://github.com/Panda-995/ezbookkeeping-pixel/pkgs/container/ezbookkeeping-pixel)

这是一个基于 [mayswind/ezbookkeeping](https://github.com/mayswind/ezbookkeeping) 的独立社区优化版本。原项目、原始业务能力与文档由 MaysWind 开发和维护；本仓库保留 MIT License 与原作者版权声明。本版本并非上游官方发行版。

本次重构主要补充：

- 桌面端和移动端统一的像素风设计系统、深色模式、高对比焦点态与减少动态效果支持。
- 账户、子账户余额和交易的明确编辑入口，补齐空状态快捷操作。
- MCP 对账户、余额、交易、分类、标签、标签组的完整认证 CRUD。
- 金额调整通过“余额修改交易”留痕，不直接篡改账户余额。
- 交易支持标签、附件、隐藏金额、地理位置与批量 REST 操作。
- 可复用 Agent Skill，以及 PowerShell / Shell API 命令脚本。
- Docker Compose、完整环境变量清单、SQLite 持久化、健康检查。
- GitHub Actions 自动测试并构建公开的 `linux/amd64`、`linux/arm64` 镜像。
- 现有“导入交易”流程兼容原版 ezBookkeeping `v0.1.0` 起导出的 CSV / TSV 数据。

## 导入原版 ezBookkeeping 数据

无需安装迁移插件，也无需使用单独的迁移工具。打开桌面端“交易详情”，点击“导入”，在现有文件类型中选择“原版 ezBookkeeping 数据导出文件”，即可导入原项目导出的 CSV 或 TSV。

兼容范围：

- `v0.1.0–v0.4.1`：自动兼容只有分钟精度的交易时间和旧 `Comment` 备注列。
- `v0.5.0` 及更高版本：兼容当前 `Description`、地理位置、时区、账户币种和标签格式。
- 兼容收入、支出、转账、余额修改和跨币种转账；导入前仍可在原有预览页面修改账户、分类和标签映射。
- 同一兼容解析器同时用于网页导入、REST API 和 `transaction-import` 命令行操作。

原版 CSV / TSV 导出文件只包含交易及其引用的账户、分类和标签，不包含用户密码、应用设置或交易图片；这是原项目导出格式本身的范围。

## Docker 部署

仓库内的 `compose.yaml` 直接使用已构建好的公开双架构镜像：

```bash
docker compose pull
docker compose up -d
```

浏览器访问 `http://localhost:18088/`。镜像标签固定为：

```bash
docker pull ghcr.io/panda-995/ezbookkeeping-pixel:latest
```

数据通过 `./data` 和 `./storage` 绑定挂载到宿主机，不使用 Docker 命名卷。完整端口、目录、反向代理、备份和 198 个应用环境变量说明见 [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)。

## MCP 与 Agent

Compose 默认开启 API Token；如需启用 MCP，可在 `compose.yaml` 的 `environment` 中增加 `EBK_MCP_ENABLE_MCP=true`。所有调用仍必须使用当前用户签发的 Bearer Token，并受数据归属和服务层校验约束。

- MCP 地址：`https://你的域名/mcp`
- REST API：`https://你的域名/api/v1`
- Agent Skill：`skills/ezbookkeeping/`
- 完整接入与安全说明：[MCP_AGENT_GUIDE.md](MCP_AGENT_GUIDE.md)

## 文档导航

- [产品与设计说明](PRODUCT.md)
- [原项目分析](ORIGINAL_PROJECT_ANALYSIS.zh_CN.md)
- [Docker 部署与环境变量](DOCKER_DEPLOYMENT.md)
- [MCP 与 Agent 使用指南](MCP_AGENT_GUIDE.md)
- [后续功能路线图](FUTURE_FEATURES.md)
- [来源与署名](ATTRIBUTION.md)
- [原项目中文文档](https://ezbookkeeping.mayswind.net/zh_Hans/)

## 本地验证

```bash
npm ci
npx vue-tsc --noEmit
npx eslint .
npm test
npm run build
go test ./...
```

## 许可证

使用 MIT License。详细来源、修改范围和非官方声明见 [ATTRIBUTION.md](ATTRIBUTION.md)。
