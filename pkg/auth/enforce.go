package auth

import (
	goerrors "errors"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/golang-jwt/jwt/v5"

	"gz-dango/pkg/errors"
)

const (
	SubKey = "sub"
	ObjKey = "obj"
	ActKey = "act"
)

// AuthEnforcer 管理身份验证令牌和授权权限
// 它提供线程安全的操作来存储和检索用户声明、角色权限和URL访问控制映射
type AuthEnforcer struct {
	// jwt的密钥
	key []byte

	// jwt的黑名单缓存
	blacklist BlacklistManager

	// enforcer 用于访问控制
	enforcer *casbin.Enforcer
}

// NewAuthEnforcer 创建一个新的认证缓存实例
// 返回初始化后的AuthCache指针
func NewAuthEnforcer(enforcer *casbin.Enforcer, key string) *AuthEnforcer {
	return &AuthEnforcer{
		enforcer: enforcer,
		key:      []byte(key),
	}
}

// SetBlacklist 设置黑名单缓存
func (a *AuthEnforcer) SetBlacklist(manager BlacklistManager) {
	a.blacklist = manager
}

func (a *AuthEnforcer) AddToBlacklist(token string, duration time.Duration) error {
	if a.blacklist != nil {
		return a.blacklist.AddToBlacklist(token, 24*time.Hour)
	}
	return nil
}

// Authentication 验证给定令牌的有效性并返回相应的用户认证信息
// token：待验证的JWT令牌字符串
// 返回用户认证信息和可能发生的错误（如无效令牌、已过期等）
func (c *AuthEnforcer) Authentication(token string) (*UserClaims, *errors.Error) {
	// 检查黑名单
	if c.blacklist != nil {
		blacklisted, err := c.blacklist.IsBlacklisted(token)
		if err != nil {
			return nil, errors.FromError(err)
		}
		if blacklisted {
			return nil, ErrTokenRevoked
		}
	}

	// 解析token
	parsedToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return c.key, nil
	})
	if err != nil {
		if goerrors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired.WithCause(err)
		}
		return nil, ErrInvalidToken.WithCause(err)
	}

	// 验证token有效性
	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	// 类型断言获取claims
	claims, ok := parsedToken.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// Authorization 检查特定角色是否具有对某个HTTP方法和URL路径组合的访问权限
// roleId：请求访问的角色
// url：请求的目标URL路径
// method：HTTP请求方法（GET/POST等）
// 返回是否有访问权限的布尔结果
func (c *AuthEnforcer) Authorization(role, url, method string) (bool, *errors.Error) {
	ok, err := c.enforcer.Enforce(role, url, method)
	if err != nil {
		return false, errors.FromError(err)
	}
	return ok, nil
}

// AddPolicies 批量添加授权策略规则
// rules: 要添加的策略规则列表，每个规则是一个字符串切片
// 返回值: 如果添加成功返回nil，否则返回相应的错误信息
func (c *AuthEnforcer) AddPolicy(sub, obj, act string) error {
	if _, err := c.enforcer.AddPolicy(sub, obj, act); err != nil {
		return err
	}
	return nil
}

// // RemovePolicies 批量移除授权策略规则
// // rules: 要移除的策略规则列表，每个规则是一个字符串切片
// // 返回值: 如果移除成功返回nil，否则返回相应的错误信息
func (c *AuthEnforcer) RemovePolicy(sub, obj, act string) error {
	if _, err := c.enforcer.RemovePolicy(sub, obj, act); err != nil {
		return err
	}
	return nil
}

// AddGroupPolicies 批量添加用户组策略规则
// 返回值: 如果添加成功返回nil，否则返回相应的错误信息
func (c *AuthEnforcer) AddGroupPolicy(sub, obj string) error {
	if _, err := c.enforcer.AddGroupingPolicy(sub, obj); err != nil {
		return err
	}
	return nil
}

// RemoveGroupPolicies 批量移除用户组策略规则
// 返回值: 如果移除成功返回nil，否则返回相应的错误信息
func (c *AuthEnforcer) RemoveGroupPolicy(index int, value string) error {
	if _, err := c.enforcer.RemoveFilteredGroupingPolicy(index, value); err != nil {
		return err
	}
	return nil
}
