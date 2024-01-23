package authorization

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var JWT_SINGING_KEY string = "3kq9p5d6s8v2y5x1z7n4m8c0q2r8t6w9"

// TokenClaims JWT 令牌包含一个名为 "id" 的字段，用于存储用户 ID
type TokenClaims struct {
	jwt.StandardClaims
	ID string `json:"id"`
}

// AuthInterceptor 认证拦截器，对以authorization为头部，形式为`bearer token`的Token进行验证
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, " %v", err)
	}
	//使用context.WithValue添加了值后，可以用Value(key)方法获取值
	newCtx := context.WithValue(ctx, tokenInfo.ID, tokenInfo)
	//log.Println(newCtx.Value(tokenInfo.ID))
	return newCtx, nil
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID string) (string, error) {
	claims := &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			// 设置令牌有效期等
		},
		ID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SINGING_KEY))
}

// parseToken 解析并验证 JWT 令牌
func parseToken(tokenString string) (*TokenClaims, error) {
	// 替换 YOUR_SIGNING_KEY 为你的密钥
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SINGING_KEY), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
