package main

// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

type HookResponse struct {
	Msg string
}

const (
	GIT_STATUS_FAILURE string = "failure"
	GIT_STATUS_PENDING string = "pending"
	GIT_STATUS_SUCCESS string = "success"
)

type PyLintConfig struct {
	Github struct {
		// can be 0 to ignore admin
		AdminId       int64  `env:"PYLINTGO_GITHUB_ADMINID"         envDefault:"0"`
		IntegrationID int64  `env:"PYLINTGO_GITHUB_INTEGRATIONID"   envDefault:"0"`
		Secret        string `env:"PYLINTGO_GITHUB_SECRET"          envDefault:"dummy"`
		KeyPath       string `env:"PYLINTGO_GITHUB_KEYPATH"         envDefault:"/keys/key.pem"`
	}
	Database struct {
		Path string `env:"PYLINTGO_DB_PATH"         envDefault:"test.db"`
	}
	Port       int    `env:"PYLINTGO_PORT"        envDefault:"4444"`
	PublicPort int    `env:"PYLINTGO_PUBLICPORT"  envDefault:"9097"`
	Url        string `env:"PYLINTGO_STATEURL"    envDefault:"http://subdomain.domain.com/"`
	Name       string `env:"PYLINTGO_NAME"        envDefault:"PyLinter"`
}
