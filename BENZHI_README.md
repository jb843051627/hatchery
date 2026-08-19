# hatchery

孵化场监控系统 — Go 后端 + 前端页面。

## 构建

```bash
docker build -t benzhi/hatchery:latest -f benzhi.Dockerfile .
```

## 运行

```bash
docker run -p 8080:8080 -e HATCHERY_DB=/data/hatchery.db -v ./data:/data benzhi/hatchery:latest
```
