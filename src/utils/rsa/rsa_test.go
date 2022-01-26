package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"
)

func TestEncodeKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader,1024)
	if err != nil {
		t.Fatal(err)
		return
	}
	block,err := x509.MarshalPKCS8PrivateKey(key)
	t.Log(base64.StdEncoding.EncodeToString(block))
	block = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	t.Log(base64.StdEncoding.EncodeToString(block))
	t.Log(EncodePrivateKey(key))
	t.Log(EncodePublicKey(&key.PublicKey))
	t.Log(key.PublicKey.Size() - 11)


}
/*
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAKwaR3rfYNc6fD/d
+gNcn1aGy8O4/UQze1d2RLvcevmIPVxtJYkauqIX77zxfVGIXCB4Y7cv4vtVrfVJ
lj4LZcVc5KnnYZkIFU+zmoXdudUqSzrTXq1HI88/hp0k9Mr97KrNPinT/Hbd+alr
bLwbQaXfbXtqPucAxKoOixjcBTZHAgMBAAECgYA8p/5tZfFBqhFEiT2mlaxq2JNU
ZgyNTv+3sa1D8M8+xy+pNaa3Db6dhoYuN4aNh9vAbe3nEG+VWXs4KjlToLy7IGAZ
nOzWZkYRLPH4JICLvN6RFAYIwyPjS8wLXzBVLh6HKZoe8zfICBJPFIHbSLo7CRiX
8OCYuq4NuLJu0lZmcQJBAMTGrd/ESrXrBj8vhXu4iAgHZ5K9tqcQ36oac0FL0yMW
GEPFp4VCW53EIh/PRSvUitHMXwSSKPFvNG4tnIPN6SMCQQDf5ovg4n6s/31bd9Lk
YV4vscVvrP20nXLyBNrrONipY1v75tjrR8IZXqgXto8tySSuUV3noWt0I1Y1UmK7
3NqNAkEAu0OcuyRSOVhGZKFz9e8CPinVzpePGOT9BiQP5WcksvJW+0BCEZa6G6VJ
GF6npHOr/Mby8iWqo0HCcswjdGfkYwJAXSzrHw4Cm2nDODJYQBRJBt4bBMtf1S8E
q7TbibHhcDRaDi1WLitxme8rUpr1YJ9pNWXFB2TEe9NMx+neDsHs7QJAbNh0EPc6
tY9hsKouLrL0vVz1vP2JaxpGAo/seShata7JLR1FpyVt1YRMYPfjYG+XGwLb/4l0
8/tS0s8SwE3ajQ==


*/