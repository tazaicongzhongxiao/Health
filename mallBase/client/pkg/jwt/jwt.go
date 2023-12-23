package jwt

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenExpired     = app.Err(app.Fail, "Token 已过期")
	TokenNotValidYet = app.Err(app.Fail, "Token 未激活")
	TokenMalformed   = app.Err(app.Fail, "这不是 Token")
	TokenInvalid     = app.Err(app.Fail, "无法解析的 Token")
	jwfConf          jwfConfig
)

type jwfConfig struct {
	Secret  string        `mapstructure:"secret"`
	TimeOut time.Duration `mapstructure:"timeout"`
}

// 载荷，可以加一些自己需要的信息
type claims struct {
	Uid      int64  // 用户UID
	Type     int8   // 用户登录类型 1 用户 3总后台
	UserName string // 签名密钥
	Platform int8   // 登录平台
	jwt.StandardClaims
}

func Start() {
	if err := config.Config().Bind(app.CfgName, "jwt", &jwfConf, func() {
		_ = config.Config().Bind(app.CfgName, "jwt", &jwfConf, nil)
	}); err != nil {
		log.Error("JWT配置读取错误", err.Error())
		panic(err)
		return
	}
}

// CreateToken
// @Description: 生成一个token
// @param uid 用户ID
// @param userName
// @param Type 用户类型 1用户 2企业 3管理员
// @Param platform 平台标记
// @param expireDuration
// @return tokenString
// @return err
func CreateToken(uid int64, userName string, Type int8, platform int8, expireDuration time.Duration) (tokenString string, expiresAt int64, err error) {
	// The token content.
	// iss: （Issuer）签发者
	// iat: （Issued At）签发时间，用Unix时间戳表示
	// exp: （Expiration Time）过期时间，用Unix时间戳表示
	// aud: （Audience）接收该JWT的一方
	// sub: （Subject）该JWT的主题
	// nbf: （Not Before）不要早于这个时间
	// jti: （JWT ID）用于标识JWT的唯一ID
	if expireDuration == 0 {
		expireDuration = jwfConf.TimeOut * time.Hour
	}
	expiresAt = time.Now().Add(expireDuration).Unix()
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		Uid:      uid,
		UserName: userName,
		Type:     Type,
		Platform: platform,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt,
		},
	}).SignedString([]byte(jwfConf.Secret))
	return
}

// ParseToken
// @Description: 解析 Token
// @param tokenString
// @return *claims
// @return error
func ParseToken(tokenString string) (*claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwfConf.Secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	} else {
		if val, ok := tokenClaims.Claims.(*claims); ok && tokenClaims.Valid {
			return val, nil
		}
	}
	return nil, TokenInvalid
}

// RefreshToken
// @Description: 更新token
// @param tokenString
// @return string
// @return error
func RefreshToken(tokenString string) (newTokenString string, expiresAt int64, err error) {
	if token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwfConf.Secret), nil
	}); err != nil {
		return "", 0, err
	} else {
		if val, ok := token.Claims.(*claims); ok && token.Valid {
			return CreateToken(val.Uid, val.UserName, val.Type, val.Platform, 0)
		}
		return "", 0, TokenInvalid
	}
}
