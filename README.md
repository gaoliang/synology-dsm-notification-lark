# Synology DSM notification lark

Make Synology DSM notification webhook compatible with lark webhook

将群晖的通知消息转发到飞书

DockerHub: https://hub.docker.com/r/gaoliang/synology-dsm-notification-lark

### Usage
1. start the conatainer 
```bash
docker run -d -p 10080:8080 -e LARK_WEBHOOK_URL=https://replace.with.your.lark.custom.bot.webbhook.url gaoliang/synology-dsm-notification-lark
```

> If you setup your bot with [signature verification](https://open.larksuite.com/document/uAjLw4CM/ukTMukTMukTM/bot-v3/use-custom-bots-in-a-group#348211be), you can pass your keystore through `LARK_SECRET` param.

2. config DSM notification webhook, set webhook url to POST http://localhost:10080/lark and add a `content` field in http body.
