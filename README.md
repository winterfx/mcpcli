# MCP-CLI

MCP-CLI 是一个用于管理和调用 MCP 服务器的命令行工具。

## 功能特性

- 管理多个 MCP 服务器
- 交互式命令行界面
- 支持查看和调用服务器工具(Tools)
- 支持查看服务器提示(Prompts)和资源(Resources)
- 支持 tool的调用

## 安装
```bash
go install github.com/winterfx/mcp-cli@latest
```

## 配置

创建配置文件 `~/.mcp-cli.json`:

```json
{
  "mcpServers": {
    "server1": {
      "command": "/path/to/server",
      "args": ["--arg1", "--arg2"],
      "env": {
        "KEY": "value"
      }
    }
  }
}
```

## 使用方法

基本命令:

```bash
# 列出所有服务器
mcp-cli server list

# 检查特定服务器
mcp-cli server inspect -n server1
```

交互式命令:
```bash
> tools     # 显示可用工具
> prompts   # 显示可用提示
> resources # 显示可用资源
> call tool-name {"param": "value"}  # 调用工具
> help      # 显示帮助信息
> exit      # 退出
```

## 许可证

MIT License
