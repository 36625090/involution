package transport

type SignPolicy string

const (
	SignPolicyDeny  SignPolicy = "deny"
	SignPolicyAllow SignPolicy = "allow"
)

const GlobalSignKey = "global"

type Settings struct {
	SignType      string            `json:"sign_type" hcl:"sign_type" default:"md5"`
	SignKeys      map[string]string `json:"sign_keys" hcl:"sign_keys"`
	DefaultPolicy SignPolicy        `json:"default_policy" hcl:"default_policy"`
}