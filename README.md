# VX Replacer
Simple discord bot to automatically replace the following links automatically:
```
twitter/x.com -> vxtwitter.com
bsky.app/.social -> bskx.app
instagram.com -> ddinstagram.com 
```

## Setup

Add your discord token into the `.env` file which you can create by copying the `.env.example` and editing.
```bash
$ cp .env.example .env
$ nano .env
```
Start via docker compose:
```bash
$ docker compose up -d
```
Logs will be visable inside of the /logs directory.
