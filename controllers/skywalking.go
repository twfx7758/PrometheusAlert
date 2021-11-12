package controllers

import (
	"PrometheusAlert/model"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
)

type SkywalkingController struct {
	beego.Controller
}

type Skywalking struct {
	ScopeId      int    `json:"scopeid"`
	Name         string `json:"name"`
	Id0          int    `json:"id0"`
	Id1          int    `json:"id1"`
	AlarmMessage string `json:"alarmmessage"`
	StartTime    int    `json:"starttime"`
}

func (c *SkywalkingController) SkywalkingWorkWechat() {
	alert := Skywalking{}
	logsign := "[" + LogsSign() + "]"
	logs.Info(logsign, string(c.Ctx.Input.RequestBody))
	json.Unmarshal(c.Ctx.Input.RequestBody, &alert)
	c.Data["json"] = SendMessageSkywalking(alert, 12, logsign, "", "", "", "", "", "", "", "", "", "", "", "")
	logs.Info(logsign, c.Data["json"])
	c.ServeJSON()
}

//typeid 为0,触发电话告警和钉钉告警, typeid 为1 仅触发dingding告警
func SendMessageSkywalking(message Skywalking, typeid int, logsign, ddurl, wxurl, fsurl, txdx, txdh, hwdx, rlydh, alydx, alydh, email, bddx, groupid string) string {
	Title := beego.AppConfig.String("title")
	var DDtext, RLtext, FStext, WXtext, EmailMessage, titleend string
	//告警级别定义 0 信息,1 警告,2 一般严重,3 严重,4 灾难
	AlertLevel := []string{"信息", "警告", "一般严重", "严重", "灾难"}
	titleend = "故障恢复信息"
	timeobj := time.Unix(int64(message.StartTime), 0)
	date := timeobj.Format("2006-01-02 15:04:05")
	model.AlertsFromCounter.WithLabelValues("skywalking", message.AlarmMessage, "4", "", "resolved").Add(1)
	WXtext = "[" + Title + "skywalking" + message.Name + "]" + "**\n>`告警级别:`" + AlertLevel[4] + "\n`开始时间:`" + date + "\n" + message.AlarmMessage
	PhoneCallMessage = message.AlarmMessage
	//触发email
	if typeid == 1 {
		if email == "" {
			email = beego.AppConfig.String("Default_emails")
		}
		SendEmail(EmailMessage, email, logsign)
	}
	//触发钉钉
	if typeid == 2 {
		if ddurl == "" {
			ddurl = beego.AppConfig.String("ddurl")
		}
		PostToDingDing(Title+titleend, DDtext, ddurl, "", logsign)
	}
	//触发微信
	if typeid == 3 {
		if wxurl == "" {
			wxurl = beego.AppConfig.String("wxurl")
		}
		PostToWeiXin(WXtext, wxurl, "", logsign)
	}

	//取到手机号

	//触发电话告警
	if typeid == 4 {
		if txdh == "" {
			txdh = GetUserPhone(1)
		}
		PostTXphonecall(PhoneCallMessage, txdh, logsign)
	}
	//触发腾讯云短信告警
	if typeid == 5 {
		if txdx == "" {
			txdx = GetUserPhone(1)
		}
		PostTXmessage(PhoneCallMessage, txdx, logsign)
	}
	//触发华为云短信告警
	if typeid == 6 {
		if hwdx == "" {
			hwdx = GetUserPhone(1)
		}
		PostHWmessage(PhoneCallMessage, hwdx, logsign)
	}
	//触发阿里云短信告警
	if typeid == 7 {
		if alydx == "" {
			alydx = GetUserPhone(1)
		}
		PostALYmessage(PhoneCallMessage, alydx, logsign)
	}
	//触发阿里云电话告警
	if typeid == 8 {
		if alydh == "" {
			alydh = GetUserPhone(1)
		}
		PostALYphonecall(PhoneCallMessage, alydh, logsign)
	}
	//触发容联云电话告警
	if typeid == 9 {
		if rlydh == "" {
			rlydh = GetUserPhone(1)
		}
		PostRLYphonecall(PhoneCallMessage, rlydh, logsign)
	}
	//触发飞书
	if typeid == 10 {
		if fsurl == "" {
			fsurl = beego.AppConfig.String("fsurl")
		}
		PostToFS(Title+titleend, FStext, fsurl, "", logsign)
	}
	//触发TG
	if typeid == 11 {
		SendTG(PhoneCallMessage, logsign)
	}
	//触发企业微信消息
	if typeid == 12 {
		SendWorkWechat(beego.AppConfig.String("WorkWechat_ToUser"), beego.AppConfig.String("WorkWechat_ToParty"), beego.AppConfig.String("WorkWechat_ToTag"), WXtext, logsign)
	}
	//触发百度云短信告警
	if typeid == 13 {
		if bddx == "" {
			bddx = GetUserPhone(1)
		}
		PostBDYmessage(PhoneCallMessage, bddx, logsign)
	}
	//触发百度Hi(如流)
	if typeid == 14 {
		if groupid == "" {
			groupid = beego.AppConfig.String("BDRL_ID")
		}
		PostToRuLiu(groupid, RLtext, beego.AppConfig.String("BDRL_URL"), logsign)
	}
	//触发Bark
	if typeid == 15 {
		SendBark(PhoneCallMessage, logsign)
	}
	return "告警消息发送完成."
}
