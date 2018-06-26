package admin

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type BaseController struct {
	beego.Controller
	RedisConn redis.Conn
}

type JSONStruct struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
	Msg    string      `json:"msg"`
}

func (c *BaseController) RedisInit() {
	var err error
	redisConf := fmt.Sprintf("%s:%s", beego.AppConfig.String("REDIS_HOST"), beego.AppConfig.String("REDIS_PORT"))
	c.RedisConn, err = redis.Dial("tcp", redisConf)
	if err != nil {
		c.Data["json"] = JSONStruct{"error", 1001, "", "redis连接失败"}

		c.StopJson()
	}
}

func (c *BaseController) StopJson() {
	c.ServeJSON()
	c.Finish()
	c.StopRun()
}