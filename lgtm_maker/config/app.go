package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/naoina/kocha"
	"github.com/naoina/kocha/log"
)

var (
	AppName   = "myapp"
	AppConfig = &kocha.Config{
		Addr:          kocha.SettingEnv("KOCHA_ADDR", "127.0.0.1:9100"),
		AppPath:       rootPath,
		AppName:       AppName,
		DefaultLayout: "app",
		Template: &kocha.Template{
			PathInfo: kocha.TemplatePathInfo{
				Name: AppName,
				Paths: []string{
					filepath.Join(rootPath, "app", "view"),
				},
			},
			FuncMap: kocha.TemplateFuncMap{},
		},

		// Logger settings.
		Logger: &kocha.LoggerConfig{
			Writer:    os.Stdout,
			Formatter: &log.LTSVFormatter{},
			Level:     log.INFO,
		},

		// Middlewares.
		Middlewares: []kocha.Middleware{
			&kocha.RequestLoggingMiddleware{},
			&kocha.SessionMiddleware{},
			&kocha.FlashMiddleware{},
		},

		// Session settings
		Session: &kocha.SessionConfig{
			Name: "myapp_session",
			Store: &kocha.SessionCookieStore{
				// AUTO-GENERATED Random keys. DO NOT EDIT.
				SecretKey:  "\xa1\xe9}\xe2ц%r&Q|*i\x15p\x84w\x14\x9b$\xe2\xfdZ\xcb\xf5JxAq\a\x1co",
				SigningKey: "SV=\xc6\x1d\xad\xad\tb\xa3Y\x82&p\x82\x8b",
			},

			// Expiration of session cookie, in seconds, from now.
			// Persistent if -1, For not specify, set 0.
			CookieExpires: time.Duration(90) * time.Hour * 24,

			// Expiration of session data, in seconds, from now.
			// Perssitent if -1, For not specify, set 0.
			SessionExpires: time.Duration(90) * time.Hour * 24,
			HttpOnly:       false,
		},

		MaxClientBodySize: 1024 * 1024 * 10, // 10MB
	}

	_, configFileName, _, _ = runtime.Caller(0)
	rootPath                = filepath.Dir(filepath.Join(configFileName, ".."))
)
