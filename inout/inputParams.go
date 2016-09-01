package inout

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JREAMLU/core/global"
	"github.com/JREAMLU/core/guid"
	"github.com/beego/i18n"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
	"github.com/pquerna/ffjson/ffjson"
)

type Header struct {
	Source      []string `json:"Source" valid:"Required"`
	Version     []string `json:"Version" `
	SecretKey   []string `json:"Secret-Key" `
	RequestID   []string `json:"Request-ID" valid:"Required"`
	ContentType []string `json:"Content-Type" valid:"Required"`
	Accept      []string `json:"Accept" valid:"Required"`
	Token       []string `json:"Token" `
	IP          []string `json:"Ip" valid:"Required"`
}

type Result struct {
	CheckRes  map[string]string
	RequestID string
	Message   string
}

func InputParams(r *context.BeegoInput) (http.Header, []byte) {
	rawMetaHeader := r.Context.Request.Header
	rawDataBody := r.RequestBody

	js, _ := json.Marshal(rawMetaHeader)

	//记录参数日志
	beego.Trace("input params header :" + string(js))
	beego.Trace("input params body :" + string(rawDataBody))

	return rawMetaHeader, rawDataBody
}

func InputParamsNew(r *context.Context) map[string]interface{} {
	r.Request.ParseForm()

	headerMap := r.Request.Header
	header, _ := json.Marshal(headerMap)

	body := r.Input.RequestBody

	cookiesSlice := r.Request.Cookies()
	cookies, _ := json.Marshal(cookiesSlice)

	querystrMap := r.Request.Form
	querystr, _ := json.Marshal(querystrMap)

	beego.Trace("input params header" + string(header))
	beego.Trace("input params body" + string(body))
	beego.Trace("input params cookies" + string(cookies))
	beego.Trace("input params querystr" + string(querystr))

	data := make(map[string]interface{})
	mu.Lock()
	data["header"] = header
	data["body"] = body
	data["cookies"] = cookies
	data["querystr"] = querystr
	data["headermap"] = headerMap
	data["cookiesslice"] = cookiesSlice
	data["querystrmap"] = querystrMap
	mu.Unlock()

	return data
}

/**
 *	@auther		jream.lu
 *	@intro		入参验证
 *	@logic
 *	@todo		返回值
 *	@meta		meta map[string][]string	   rawMetaHeader
 *	@data		data []byte 					rawDataBody 签名验证
 *	@data		data ...interface{}	切片指针	rawDataBody
 *	@return 	返回 true, metaMap, error
 */
func InputParamsCheck(data map[string]interface{}, stdata ...interface{}) (result Result, err error) {
	headerRes, err := HeaderCheck(data)
	if err != nil {
		return headerRes, err
	}

	//DataParams check
	// valid := validation.Validation{}
	//
	// for _, val := range stdata {
	// 	is, err := valid.Valid(val)
	// 	if err != nil {
	// 		// handle error
	// 		beego.Trace(i18n.Tr(global.Lang, "outputParams.SYSTEMILLEGAL") + err.Error())
	// 	}
	//
	// 	if !is {
	// 		for _, err := range valid.Errors {
	// 			beego.Trace("input params body check : " + i18n.Tr(global.Lang, "outputParams.DATAPARAMSILLEGAL") + err.Key + ":" + err.Message)
	// 			result.CheckRes = nil
	// 			result.RequestID = metaCheckResult.CheckRes["request-id"]
	// 			result.Message = i18n.Tr(global.Lang, "outputParams.DATAPARAMSILLEGAL") + " " + err.Key + ":" + err.Message
	// 			return result, errors.New(i18n.Tr(global.Lang, "outputParams.DATAPARAMSILLEGAL"))
	// 		}
	// 	}
	// }
	//
	// //sign check
	// err = sign.ValidSign(rawDataBody, beego.AppConfig.String("sign.secretKey"))
	// if err != nil {
	// 	result.CheckRes = nil
	// 	result.RequestID = metaCheckResult.CheckRes["request-id"]
	// 	result.Message = err.Error()
	// 	// return result, err
	// }

	return headerRes, nil
}

// func InputParamsCheck(meta map[string][]string, rawDataBody []byte, data ...interface{}) (result Result, err error) {
// 	//MetaHeader check
// 	metaCheckResult, err := MetaHeaderCheck(meta)
// 	if err != nil {
// 		return metaCheckResult, err
// 	}
//
// 	//DataParams check
// 	valid := validation.Validation{}
//
// 	for _, val := range data {
// 		is, err := valid.Valid(val)
//
// 		//日志
//
// 		//检查参数
// 		if err != nil {
// 			// handle error
// 			beego.Trace(i18n.Tr(global.Lang, "outputParams.SYSTEMILLEGAL") + err.Error())
// 		}
//
// 		if !is {
// 			for _, err := range valid.Errors {
// 				beego.Trace("input params body check : " + i18n.Tr(global.Lang, "outputParams.DATAPARAMSILLEGAL") + err.Key + ":" + err.Message)
// 				result.CheckRes = nil
// 				result.RequestID = metaCheckResult.CheckRes["request-id"]
// 				result.Message = i18n.Tr(global.Lang, "outputParams.DATAPARAMSILLEGAL") + " " + err.Key + ":" + err.Message
// 				return result, errors.New(i18n.Tr(global.Lang, "outputParams.DATAPARAMSILLEGAL"))
// 			}
// 		}
// 	}
//
// 	//sign check
// 	err = sign.ValidSign(rawDataBody, beego.AppConfig.String("sign.secretKey"))
// 	if err != nil {
// 		result.CheckRes = nil
// 		result.RequestID = metaCheckResult.CheckRes["request-id"]
// 		result.Message = err.Error()
// 		// return result, err
// 	}
//
// 	return metaCheckResult, nil
// }

/**
 * header参数验证
 * 将header 放入map 返回
 *
 * @meta 	meta  map[string][]string 	header信息 map格式
 */
func HeaderCheck(data map[string]interface{}) (result Result, err error) {
	var h Header
	ffjson.Unmarshal(data["header"].([]byte), &h)

	rid := GetRequestID()
	if len(h.RequestID) > 0 && h.RequestID[0] != "" {
		rid = h.RequestID[0]
	}

	result.CheckRes = nil
	result.Message = ""
	result.RequestID = rid

	ct, err := HeaderParamCheck(h.ContentType, "Content-Type")
	if err != nil {
		ct.RequestID = rid
		return ct, err
	}

	at, err := HeaderParamCheck(h.Accept, "Accept")
	if err != nil {
		at.RequestID = rid
		return at, err
	}

	valid := validation.Validation{}

	is, err := valid.Valid(&h)

	if err != nil {
		beego.Trace(
			i18n.Tr(
				global.Lang,
				"outputParams.SYSTEMILLEGAL") + err.Error(),
		)
		result.Message = i18n.Tr(global.Lang, "outputParams.SYSTEMILLEGAL")

		return result, err
	}

	if !is {
		for _, err := range valid.Errors {
			beego.Trace(
				i18n.Tr(
					global.Lang,
					"outputParams.METAPARAMSILLEGAL") + err.Key + ":" + err.Message)
			result.Message = i18n.Tr(
				global.Lang,
				"outputParams.METAPARAMSILLEGAL") + " " + err.Key + ":" + err.Message

			return result, errors.New(
				i18n.Tr(
					global.Lang,
					"outputParams.METAPARAMSILLEGAL",
				),
			)
		}
	}

	var headerMap = make(map[string]string)
	for key, val := range data["headermap"].(http.Header) {
		headerMap[key] = val[0]
	}
	headerMap["request-id"] = rid
	result.CheckRes = headerMap

	return result, nil
}

//HeaderParamCheck 验证header固定信息
func HeaderParamCheck(h []string, k string) (result Result, err error) {
	if h[0] != beego.AppConfig.String(k) {
		message := ""
		switch k {
		case "Content-Type":
			message = i18n.Tr(
				global.Lang,
				"outputParams.CONTENTTYPEILLEGAL",
			)
		case "Accept":
			message = i18n.Tr(
				global.Lang,
				"outputParams.ACCEPTILLEGAL",
			)
		}

		result.CheckRes = nil
		result.Message = message
		return result, errors.New(message)
	}

	return result, nil
}

//request id增加
func GetRequestID() string {
	var requestID bytes.Buffer
	requestID.WriteString(beego.AppConfig.String("appname"))
	requestID.WriteString("-")
	requestID.WriteString(guid.NewObjectId().Hex())
	return requestID.String()
}
