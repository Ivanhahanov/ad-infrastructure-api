package config

type (
	Config struct {
		TerraformProjectPath string `yaml:"terraform_project_path"`
		SshKeys              string `yaml:"ssh_keys"`
		AdminMachine         bool   `yaml:"admin_machine"`
		Network              string `yaml:"network"`
		CheckerPassword      string `yaml:"checker_password"`
		RoundInterval        string `yaml:"round_interval"`
		Teams                Teams  `yaml:"teams"`
		Users                Users  `yaml:"users"`

		App  App  `yaml:"app"`
		HTTP HTTP `yaml:"http"`
		Log  Log  `yaml:"log"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Resources struct {
		Memory int `yaml:"memory"`
		VCPU   int `yaml:"vcpu"`
	}

	Teams struct {
		Number    int       `yaml:"number"`
		Resources Resources `yaml:"resources"`
	}

	Users struct {
		Number    int       `yaml:"number"`
		Resources Resources `yaml:"resources"`
	}
)
