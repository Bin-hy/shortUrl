# 短链接服务 GO实践

思路：
v1 版本使用 MYSQL 主键ID 生成短链接，使用MYSQL查询短链

v2 版本使用 Redis 自增ID 生成短链接，使用布隆过滤器优化DB查询短链的性能

v3 版本使用 雪花算法生成 唯一ID 作为短链接，使用布隆过滤器优化DB查询短链的性能

## 项目结构

- `cmd/url-service`: 主命令
- `internal/biz`: 业务逻辑
- `internal/conf`: 配置
- `internal/data`: 数据访问

## 依赖
Redis
MySQL

### 依赖启动
```bash
docker compose up -d
```

## 项目启动

```
go run cmd/url-service/main.go
```

