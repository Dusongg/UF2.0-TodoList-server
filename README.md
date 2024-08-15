# 安装

## 1.1 Docker

1. `git pull git@github.com:Dusongg/UF2.0-TodoList-server.git`或`git pull https://github.com/Dusongg/UF2.0-TodoList-server.git`

2. 进入项目文件，在docker-compose.yml文件目录下，运行`docker compose up`

3. 查看是否运行成功：`grpcurl -plaintext -d '{"name": "dusong"}' localhost:8001 notification.Service/SayHello`

![image-20240815111932773](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240815111932773.png)

- 所有API

  ![image-20240815112412028](https://typora-dusong.oss-cn-chengdu.aliyuncs.com/image-20240815112412028.png)

