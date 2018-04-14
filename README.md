# GitHub webhooks cather in the rye!

[![Build Status](https://travis-ci.org/mxssl/go-github-webhook-to-tg.svg?branch=master)](https://travis-ci.org/mxssl/go-github-webhook-to-tg)

Simple golang app that catch webhooks from GitHub and send them to telegram group or as personal message.

![logo](https://raw.githubusercontent.com/mxssl/go-github-webhook-to-tg/master/img/C.png)

To run this app you need 3 things:
*  Configure webhooks in your GitHub repo setting
*  Telegram bot id from [BotFather](https://t.me/BotFather)
*  ChatID - send `/getgroupid` to [myidbot](https://t.me/myidbot)

#### How to run this app.

Build container

`docker build -t go-github-webhook-to-tg -f ./docker/Dockerfile .`

Then you have two options. 

Run with pure docker

```
docker container run \
  --name=bot \
  --rm \
  --publish 9191:80 \
  --env TGTOKEN="your_token" \
  --env CHATID="your_chatid" \
  --detach \
  go-github-webhook-to-tg
```

Convenient way is to use [docker-compose](https://docs.docker.com/compose/)

Edit this file `docker-compose.yml`

```
version: '3.4'

services:
  go-github-webhook-to-tg:
    image: go-github-webhook-to-tg
    environment:
      - TGTOKEN="your_token"
      - CHATID="your_chatid"
    ports:
      - "9191:80"
```

And then just use this command

`docker-compose up -d`

To check health of the app you can use this link:
[http://your_ip:your_port/health](http://your_ip:your_port/health)

If you see "Go!", it means that everything is OK.
