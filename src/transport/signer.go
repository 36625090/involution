package transport

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/36625090/involution/utils"
)

var errInvalidSign = errors.New("invalid signature")
var errInvalidHeaderSignKey = errors.New("headers[X-Client-ID] cannot be empty")
var errInvalidSignKey = errors.New("invalid signature key id")
var errInvalidDefaultSignKey = errors.New("invalid signature global default key")

type Signer interface {
	Sign(keyId string, resp Codec) (string, error)
	Verify(keyId, sign string, req Codec) error
}

type signer struct {
}

func (m *signer) build(req Codec) *bytes.Buffer {
	params := req.Map()
	var buf bytes.Buffer
	for _, v := range req.Keys() {
		if v == "sign" {continue}
		if utils.IsNil(params[v]){continue}
		buf.WriteString(v)
		buf.WriteString(fmt.Sprintf("%v", params[v]))
	}
	return &buf
}

type md5Signer struct {
	signer
	settings *Settings
}

func NewMD5Signer(settings *Settings) Signer {
	return &md5Signer{settings: settings}
}

//Sign 签名方法
//keyId keys id
func (m *md5Signer) Sign(keyId string, resp Codec) (string, error) {
	if m.settings.DefaultPolicy == SignPolicyAllow {
		return "", nil
	}

	key := m.settings.SignKeys[keyId]
	if "" == key{
		return "", errInvalidDefaultSignKey
	}

	buf := m.build(resp)
	buf.WriteString(key)
	return m.md5(buf.Bytes()), nil
}

//Verify 签名校验方法
func (m *md5Signer) Verify(keyId, sign  string, req Codec) error {
	if m.settings.DefaultPolicy == SignPolicyAllow {
		return nil
	}

	if "" == keyId{
		return errInvalidHeaderSignKey
	}

	key := m.settings.SignKeys[keyId]
	if "" == key{
		return errors.New(errInvalidSignKey.Error() + ": "+keyId)
	}

	buf := m.build(req)
	buf.WriteString(key)

	if m.md5(buf.Bytes()) != sign{
		return errInvalidSign
	}
	return nil
}

func (m *md5Signer) md5(data []byte) string {
	md := md5.New()
	md.Write(data)
	return hex.EncodeToString(md.Sum(nil))
}

