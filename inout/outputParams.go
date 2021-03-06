package inout

import (
	"time"

	"github.com/JREAMLU/core/global"
	"github.com/astaxie/beego"

	"github.com/beego/i18n"
)

//logic service
const (
	SUCCESS           = 0
	DATAPARAMSILLEGAL = 10000
	METAPARAMSILLEGAL = 15000
	LOGICILLEGAL      = 20000
	SYSTEMILLEGAL     = 30000
)

// Output output
type Output struct {
	Meta       MetaList    `json:"meta"`
	StatusCode int         `json:"statusCode"`
	Message    interface{} `json:"message"`
	Data       interface{} `json:"data"`
}

// MetaList meta list
type MetaList struct {
	RequestID string    `json:"RequestID"`
	UpdatedAt time.Time `json:"updatedAT"`
	Timezone  string    `json:"timezone"`
}

// Suc success
/**
 *	@auther		jream.lu
 *	@intro		出参成功
 *	@logic
 *	@todo		返回值
 *	@params		params ...interface{}	切片指针
 *	@return 	?
 */
func Suc(data interface{}, requestID string) Output {
	var op Output
	op.Meta.RequestID = requestID
	op.Meta.UpdatedAt = time.Now()
	op.Meta.Timezone = beego.AppConfig.String("Timezone")

	op.StatusCode = SUCCESS
	op.Message = i18n.Tr(global.Lang, "outputParams.SUCCESS")
	op.Data = data

	return op
}

// Fail fail
func Fail(msg interface{}, status string, requestID string) Output {
	var op Output
	op.Meta.RequestID = requestID
	op.Meta.UpdatedAt = time.Now()
	op.Meta.Timezone = beego.AppConfig.String("Timezone")

	switch status {
	case "SUCCESS":
		op.StatusCode = SUCCESS
	case "DATAPARAMSILLEGAL":
		op.StatusCode = DATAPARAMSILLEGAL
	case "METAPARAMSILLEGAL":
		op.StatusCode = METAPARAMSILLEGAL
	case "LOGICILLEGAL":
		op.StatusCode = LOGICILLEGAL
	case "SYSTEMILLEGAL":
		op.StatusCode = SYSTEMILLEGAL
	}

	op.Message = msg
	op.Data = make(map[string]interface{})

	return op
}
