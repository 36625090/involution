package authorities

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"github.com/36625090/involution/utils"
	"testing"
)

func TestNewAuthorized(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader,2048)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(key.Size())

	pri, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
		return
	}
	pub := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	settings := Settings{
		PKCS8PrivateKey: base64.StdEncoding.EncodeToString(pri),
		PKCS1PublicKey:  base64.StdEncoding.EncodeToString(pub),
		Timeout:         3,
		AuthType: AuthTypeRedis,
	}

	handler, err := NewJwtTokenHandler(&settings)
	if err != nil {
		t.Fatal(err)
		return
	}

	auth, err := NewAuthorization(&settings, handler)
	if err != nil {
		t.Fatal(err)
		return
	}

	authed := NewAuthorized("12321","liping", map[string]interface{}{})
	authed.AccountRoles = []string{"admin", "manager"}
	authed.Principal["location"] = "shanghai"

	jwt, err := handler.GenerateToken(&authed)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(jwt)
	authorized, err := auth.Authentication(context.Background(), jwt)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(utils.JSONPrettyDump(authorized))
}
