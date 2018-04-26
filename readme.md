## beego Steps

- 1、创建golang及beego环境
- 2、...

```
cd /Users/guilinli/go/src
git clone http://git.onemore.cc/backend/go-demo.git
cd ./jeedev-api

docker-compose up -d
docker ps -a
docker exec -ti 857fed7f95f5 bash
go build jeedev-api
./jeedev-api

http://localhost:8080/v1/app

http://localhost:8080/swagger
```

## 方法2
```
bee generate docs
bee run watchall true
```