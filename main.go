package caddy_ja4h

import (
	"net/http"

	ja4h "github.com/WJQSERVER/go-ja4h"                         // 引入 ja4h 包，用于计算 JA4H 指纹
	"github.com/caddyserver/caddy/v2"                           // 引入 Caddy v2 包
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"     // 引入 Caddyfile 配置解析包
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile" // 引入 HTTP Caddyfile 配置解析包
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"         // 引入 Caddy HTTP 模块
)

// 初始化函数，注册中间件模块和 Caddyfile 解析器
func init() {
	caddy.RegisterModule(Middleware{})                                    // 注册 Middleware 模块
	httpcaddyfile.RegisterHandlerDirective("ja4h_header", parseCaddyfile) // 注册 Caddyfile 指令
}

// Middleware 实现一个 HTTP 处理器，计算 JA4H 指纹并将其作为头部添加。
type Middleware struct{}

// CaddyModule 返回 Caddy 模块信息。
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.ja4h_header",                    // 模块 ID
		New: func() caddy.Module { return new(Middleware) }, // 创建新实例的函数
	}
}

// Provision 实现 caddy.Provisioner 接口。
func (m *Middleware) Provision(ctx caddy.Context) error {
	return nil // 目前不需要任何初始化
}

// Validate 实现 caddy.Validator 接口。
func (m *Middleware) Validate() error {
	return nil // 目前不需要任何验证
}

// ServeHTTP 实现 caddyhttp.MiddlewareHandler 接口。
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	j := ja4h.JA4H(r)       // 计算 JA4H 指纹
	r.Header.Add("Ja4h", j) // 将 JA4H 指纹添加到请求头部
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	repl.Set("ja4h", j)
	return next.ServeHTTP(w, r) // 调用下一个处理器
}

// UnmarshalCaddyfile 实现 caddyfile.Unmarshaler 接口。
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil // 目前不需要解析 Caddyfile
}

// parseCaddyfile 从 Caddyfile 中解析令牌并返回新的 Middleware 实例。
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser) // 解析 Caddyfile
	return m, err                            // 返回中间件实例和可能的错误
}

// 接口保护
var (
	_ caddy.Provisioner           = (*Middleware)(nil) // 确保 Middleware 实现了 Provisioner 接口
	_ caddy.Validator             = (*Middleware)(nil) // 确保 Middleware 实现了 Validator 接口
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil) // 确保 Middleware 实现了 MiddlewareHandler 接口
	_ caddyfile.Unmarshaler       = (*Middleware)(nil) // 确保 Middleware 实现了 Unmarshaler 接口
)
