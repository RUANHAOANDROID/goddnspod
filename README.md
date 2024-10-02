##### dnspod ddns
![1727886447464](https://github.com/user-attachments/assets/8b77d040-8d10-42af-b066-9a31e0caa098)




- 通过反查公网IPV4, 实时绑定到dnspod
- 支持IPV6公网，可配置的

添加config.yml文件
``` yaml
# 配置文件
userAgent: Hao88 DDNS/0.1Alpha(52927295@qq.com) #按照dnspod要求的格式
tokenId: xxxxxxxx #tokenID
loginToken: xxxxxxxxxxxxxxxxxxxxx #token
domain: xxx.cloud #域名
subDomain: xxx #子域名
timer: 5m0s #间隔时间 分钟
support: v6 #支持的IP版本，v4或v6
```    
挂载卷
```yaml
volumes:
   - /mnt/disk1/appdata/ddnspodgo:/app/config.yml
```

docker compose

``` yaml
version: "3"

services:
  ddnspod-go:
    image: hao88/ddnspodgo
    container_name: ddnspodgo
    environment:
      - USER_UID=1000
      - USER_GID=1000
    restart: always
    volumes: 
      - /mnt/disk1/appdata/ddnspodgo/config.yml:/app/config.yml
```
