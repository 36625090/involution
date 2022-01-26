package authorities

import "time"

type AuthorizationPolicy string

const (
	AuthorizationPolicyDeny  AuthorizationPolicy = "deny"
	AuthorizationPolicyAllow AuthorizationPolicy = "allow"
)

type AuthType string

const (
	AuthTypeJwt AuthType = "jwt"
	AuthTypeRedis AuthType = "redis"
)

//Settings for the application authorization
//authorization {
//  pkcs8_private_key = "MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAOVxpmJr4ELzX67oQl8YCrHPk61sRwESc8kAFDm9PwrY/Wd/PqBVsCQUFYBmo5dSukdJ/ZkeyXqA9pArnlqn/G42EVUjPPNURiex4W6LbSHXr/96Wt/0Ov7d+8ETkmLUZ+QsdB+9S6CrkG9pfhdUKLBoJ/YPujOhDBQvWNQSnXzXAgMBAAECgYAeTQ8LKnH4hYmaYMP7KQKojuBS49zQsG4oGmGRaoO73AJDO9O6evaDHT/lsChkoKFHLudV5HH5QrTNP2VvVYYJjAcslxVchQssuagplZtbjuixNPfv2ey9qPXafHMbdPZy97uZTZkaxQ0aMNpFOGKk/m5KOXTt8lhsZBKmpb9IqQJBAO72peFpUdCWW0Fvy4Xw9VSZq09EHHForxu6YHRu4sdAXoasLf8vmoIfHBsD87Tat01K6pxw1YaBhDry9Zkr4LMCQQD1zUKMoa9YVYDA3ty8R9DAmkYoguhAV3Sm2cf1jIF/p5kazja+L6c2BGk5sxM/AG/rLMS04vw4lPO8s2boPv1NAkEAj+Q3eKc5m7eeFaYi0HGK2Ll7vUxPMD8QCktNH29R4RcylDeDrwDUMfxXqTDVBBcbf1BYO4F6IfdFT1XTa7tPHwJBAImvDkYEE1ohmttueqqkd5RLVl0+5qWT123Ws6EhsTA2SxauyA9EVh913RNK8c7qicZr70t7kdiH5veeblhNYEkCQQDrSM+LzGB2CipariZdInt/Jkp5YVlPy6Xf8D6DUxmuSgYJSbuWrtP8dAeQuZ48gEuZZbsjjNw/ngfaXxnPHt/4"
//  pkcs1_public_key = "MIGJAoGBAOVxpmJr4ELzX67oQl8YCrHPk61sRwESc8kAFDm9PwrY/Wd/PqBVsCQUFYBmo5dSukdJ/ZkeyXqA9pArnlqn/G42EVUjPPNURiex4W6LbSHXr/96Wt/0Ov7d+8ETkmLUZ+QsdB+9S6CrkG9pfhdUKLBoJ/YPujOhDBQvWNQSnXzXAgMBAAE="
//  timeout = 86400000
//}
type Settings struct {

	AuthType AuthType `json:"auth_type" hcl:"auth_type"`

	//PKCS8 ciphertext block
	PKCS8PrivateKey string `hcl:"pkcs8_private_key" json:"pkcs8_private_key"`
	//PKCS1 ciphertext block
	PKCS1PublicKey string `hcl:"pkcs1_public_key" json:"pkcs1_public_key"`

	//Timeout token timeout
	Timeout time.Duration `hcl:"timeout" json:"timeout"`

	//AnonMethods anonymous methods
	AnonMethods []string `hcl:"anon_methods" json:"anon_methods"`

	//Disabled
	DefaultPolicy AuthorizationPolicy `hcl:"default_policy" json:"default_policy" default:"deny"`

}