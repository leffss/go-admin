package jwt

import (
	"time"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/leffss/go-admin/pkg/setting"
)

var jwtSecret []byte
var privateKey string
var jwtExpireTime time.Duration
//var rw sync.RWMutex		// 读写锁，和 Mutex 有一点区别

func Setup() {
	// jwtSecret 为程序写死密钥+配置文件设置的密钥
	privateKey = "9[ap)%O;"
	appSetting := setting.GetAppSetting()
	jwtSecret = []byte(privateKey + appSetting.JwtSecret)
	jwtExpireTime = appSetting.JwtExpireTime
}

func SetJwtSecret(secret string) {
	// 通过动态改变 jwtSecret 可以使前面分发的 token 全部失效
	//rw.Lock()	// 写锁，动态改变 jwtSecret 时打开，会降低性能
	//defer rw.Unlock()
	jwtSecret = []byte(secret)
}

type Claims struct {
	// token 不能保存密码等敏感信息，因为是可以随便解密的
	jwtGo.StandardClaims
	Id uint `json:"id"`
	UserName string `json:"username"`
}

// GenerateToken generate tokens used for auth
func GenerateToken(id uint, username string) (string, time.Time, error) {
	//rw.RLock()	// 读锁，动态改变 jwtSecret 时打开，会降低性能
	//defer rw.RUnlock()
	expireTime := time.Now().Add(jwtExpireTime)
	claims := Claims{
		jwtGo.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
		id,
		username,
	}
	tokenClaims := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, expireTime, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	//rw.RLock()	// 读锁，动态改变 jwtSecret 时打开，会降低性能
	//defer rw.RUnlock()
	tokenClaims, err := jwtGo.ParseWithClaims(token, &Claims{}, func(token *jwtGo.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
