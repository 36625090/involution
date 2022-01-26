package transport

import (
	"testing"
)

func TestMd5Signer_Sign(t *testing.T) {

	settings := Settings{
		Transport: struct {
			SignType      string            `json:"sign_type" hcl:"sign_type" default:"md5"`
			SignKeys      map[string]string `json:"sign_keys" hcl:"sign_keys"`
			DefaultPolicy SignPolicy        `json:"default_policy" hcl:"default_policy"`
		}{
			SignType: "md5",
			SignKeys: map[string]string{
				"default": "a1ede45dbcbbf7d9dc64def29f95beda",
				"user1":   "d41d8cd98f00b204e9800998ecf8427e",
			},
			DefaultPolicy: SignPolicyDeny,
		},
	}

	signer := &signer{}
	md5Sig := &md5Signer{
		settings: &settings,
	}

	req := Request{
		Method:    "account.user.login",
		Data:      "{\"name\":\"account\"}",
		Timestamp: 16123719311,
		Version:   "1.0.0",
		Sign:      "JSLJSLJSL1238210131",
		SignType:  "md5",
	}

	bs := signer.build(&req)
	sign, err := md5Sig.Sign("user1", &req)
	t.Log(sign)
	err = md5Sig.Verify("user1", sign, &req)
	t.Log("verification sign error: ", err)

	res := Response{
		Code:      0,
		Message:   "success",
		Content:   map[string]interface{}{"code": 0},
		TraceID:   "66821e08-1cc4-4ecc-8aa5-a8816208ed81",
		Timestamp: 16123719311,
		Sign:      "",
	}

	bs = signer.build(&res)
	t.Log(string(bs.String()))
	t.Log(md5Sig.Sign(GlobalSignKey, &res))
}
