package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fzzp/gotk"
)

const (
	Development = "development"
	Production  = "production"
)

type GxDuration struct {
	time.Duration
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (p *GxDuration) UnmarshalJSON(input []byte) error {
	// 去除可能存在的双引号
	s := strings.Trim(string(input), `"`)
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	p.Duration = d
	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
func (p *GxDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Duration.String())
}

// Config API服务基础配置
type Config struct {
	Mode               string     `json:"mode" validate:"required,oneof=development production"`
	Host               string     `json:"host" validate:"required,ip"`                    // API 服务主机
	Port               int        `json:"port" validate:"required,min=3000,max=65535"`    // API 服务端口
	DSN                string     `json:"dsn" validate:"required"`                        // 数据库链接
	MaxOpenConns       int        `json:"maxOpenConns" validate:"required,min=5,max=100"` // 数据库最大链接数
	MaxIdleConns       int        `json:"maxIdleConns" validate:"required,min=5,max=100"` // 数据库链接最大空闲数
	MaxIdleTime        GxDuration `json:"maxIdleTime" validate:"required"`                // 数据库链接最大空闲时间
	MaxLifetime        GxDuration `json:"maxLifetime" validate:"required"`                // 数据库连接最大存活时间
	LuckFile           string     `json:"luckFile" validate:"required"`                   // 幸运id保存文件
	TokenIssuer        string     `json:"tokenIssuer" validate:"required"`                // token 签发主体
	TokenSecret        string     `json:"tokenSecret" validate:"required,min=32"`         // token 密钥
	TokenExpire        GxDuration `json:"tokenExpire" validate:"required"`                // token 有效时间
	RefreshTokenExpire GxDuration `json:"refreshTokenExpire" validate:"required"`         // refreshToken 有效时间
	AllowFiles         []string   `json:"allowFiles" validate:"required"`
	LogLevel           string     `json:"logLevel" validate:"required,oneof=DEBUG INFO WARN ERROR"`
}

func (c *Config) Println() {
	buf, err := json.MarshalIndent(c, " ", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf))
}

// loadConfig 解析配置文件
func loadConfig(path string, conf *Config) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf, conf); err != nil {
		return err
	}
	if err := gotk.CheckStruct(conf); err != nil {
		return err
	}
	return nil
}

// NewConfig 根据配置文件路径加载配置, 并绑定到conf上
func NewConfig(path string, conf *Config) error {
	return loadConfig(path, conf)
}
