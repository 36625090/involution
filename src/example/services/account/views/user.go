package views

type Status string

type User struct {
	Mobile     string `json:"mobile" name:"手机号" validate:"required,len=11" example:"13800000000"`
	Source     string `json:"source" name:"注册渠道" validate:"required"`
	OpenId     string `json:"open_id" name:"微信open_id"`
	UnionId    string `json:"union_id" name:"微信union_id"`
	VerifyCode string `json:"verify_code" name:"短信验证码(微信登录传000000)" example:"1111" validate:"required"`
}

type SmsCodeArgs struct {
	Mobile string `json:"mobile" name:"手机号" validate:"required"`
	Kind   string `json:"kind" name:"获取验证码的类型(login register common forget)" validate:"required"`
}
type WxCodeLoginArgs struct {
	Name string `json:"name"  name:"微信公众号拼音首字母" example:"eg. wlkj" validate:"required"`
	Code string `json:"code"  name:"微信公众号获取到的code" validate:"required"`
}

type WxCodeLoginReply struct {
	Authorization string `json:"authorization"  name:"token" validate:"required"`
}

type BuyPermitArgs struct {
	ProductId int    `json:"product_id" name:"产品ID" validate:"required"`
	SkuId     string `json:"sku_id" name:"SKU ID" validate:"required"`
}
type BuyPermitReply struct {
	BuyLimited    bool   `json:"buy_limited" name:"购买是否被限制"`
	LimitedReason string `json:"limited_reason" name:"限制原因"`
}

type LoginReply struct {
	Token string `json:"token" name:"token"`
	Principal interface{} `json:"principal"`
}