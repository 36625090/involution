authorization {
  #验证类型可选 jwt redis 参考authorities包实现
  auth_type = "jwt"
  pkcs8_private_key = "MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAOVxpmJr4ELzX67oQl8YCrHPk61sRwESc8kAFDm9PwrY/Wd/PqBVsCQUFYBmo5dSukdJ/ZkeyXqA9pArnlqn/G42EVUjPPNURiex4W6LbSHXr/96Wt/0Ov7d+8ETkmLUZ+QsdB+9S6CrkG9pfhdUKLBoJ/YPujOhDBQvWNQSnXzXAgMBAAECgYAeTQ8LKnH4hYmaYMP7KQKojuBS49zQsG4oGmGRaoO73AJDO9O6evaDHT/lsChkoKFHLudV5HH5QrTNP2VvVYYJjAcslxVchQssuagplZtbjuixNPfv2ey9qPXafHMbdPZy97uZTZkaxQ0aMNpFOGKk/m5KOXTt8lhsZBKmpb9IqQJBAO72peFpUdCWW0Fvy4Xw9VSZq09EHHForxu6YHRu4sdAXoasLf8vmoIfHBsD87Tat01K6pxw1YaBhDry9Zkr4LMCQQD1zUKMoa9YVYDA3ty8R9DAmkYoguhAV3Sm2cf1jIF/p5kazja+L6c2BGk5sxM/AG/rLMS04vw4lPO8s2boPv1NAkEAj+Q3eKc5m7eeFaYi0HGK2Ll7vUxPMD8QCktNH29R4RcylDeDrwDUMfxXqTDVBBcbf1BYO4F6IfdFT1XTa7tPHwJBAImvDkYEE1ohmttueqqkd5RLVl0+5qWT123Ws6EhsTA2SxauyA9EVh913RNK8c7qicZr70t7kdiH5veeblhNYEkCQQDrSM+LzGB2CipariZdInt/Jkp5YVlPy6Xf8D6DUxmuSgYJSbuWrtP8dAeQuZ48gEuZZbsjjNw/ngfaXxnPHt/4"
  pkcs1_public_key = "MIGJAoGBAOVxpmJr4ELzX67oQl8YCrHPk61sRwESc8kAFDm9PwrY/Wd/PqBVsCQUFYBmo5dSukdJ/ZkeyXqA9pArnlqn/G42EVUjPPNURiex4W6LbSHXr/96Wt/0Ov7d+8ETkmLUZ+QsdB+9S6CrkG9pfhdUKLBoJ/YPujOhDBQvWNQSnXzXAgMBAAE="
  timeout = 86400000
  anon_methods = ["account.user.login"]
  #默认的身份校验策略，可选值有 allow deny （未配置视为deny）
  default_policy = "deny"
}

#参考transport
transport{
  sign_type = "md5"
  #默认的签名验证策略，可选值有 allow deny, allow时忽略验签和签名
  default_policy = "deny"
  sign_keys {
    #global 为返回数据包签名
    global = "521004524ef99ad954ad93c3f91c82fd"
    #user1对应客户端请求header里的X-Client-ID的值，然后取到密钥用来对客户端数据验签
    user1  = "d41d8cd98f00b204e9800998ecf8427e"
  }
}

xorm {
  driver = "mysql"
  master = "root:123456@tcp(127.0.0.1:3306)/hytest?charset=utf8"
  slaves = []
  max_conn = 36000
  max_idle = 360000
  show_sql = true
}

redis{
  addrs = ["127.0.0.1:6379"]
  key_prefix = "my_prefix"
}

