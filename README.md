# cloud-batch

批量操作多云服务，依赖 cloudpods 操作多公有云

## swagger url

```
http://0.0.0.0:5140/swagger/index.html
```

[api文档](./docs/api.md)
## 常用操作流程

1. Base URL: /api/v1 
1. 使用 POST /auth/login 获取 token
1. 点`Authorize` 或 锁标志, 看到 `ApiKeyAuth`填写 token
1. 使用 POST /batch/servers 批量创建主机
1. 使用 GET /servers 查询主机
1. 使用 DELETE /batch/servers 批量删除主机