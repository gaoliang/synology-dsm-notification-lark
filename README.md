# Synology DSM notification lark

mark Synology DSM notification webhook compatible with lark webhook

https://hub.docker.com/repository/docker/gaoliang/synology-dsm-notification-lark

### Usage
1. start the conatainer 
```bash
docker run -d -p 10080:8080 -e LARK_WEBHOOK_URL=https://replace.with.your.lark.custom.bot.webbhook.url gaoliang/synology-dsm-notification-lark
```

2. config DSM notification webhook, set address to http://localhost:10080/lark and add a `content` field in http body.

Note: DSM’s message test will not take effect because its json format is illegal