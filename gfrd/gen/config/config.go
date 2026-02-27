package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 生成器配置
type Config struct {
	Database  DatabaseConfig  `yaml:"database"`
	Generator GeneratorConfig `yaml:"generator"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

// GeneratorConfig 生成器配置
type GeneratorConfig struct {
	Backend  BackendConfig  `yaml:"backend"`
	Frontend FrontendConfig `yaml:"frontend"`
	Features FeaturesConfig `yaml:"features"`
}

// BackendConfig 后端配置
type BackendConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Output    string `yaml:"output"`
	Package   string `yaml:"package"`
	LayerMode string `yaml:"layerMode"` // simple / standard
	WithTest  bool   `yaml:"withTest"`
	WithDoc   bool   `yaml:"withDoc"`
}

// FrontendConfig 前端配置
type FrontendConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Output     string `yaml:"output"`
	APIOutput  string `yaml:"apiOutput"`
	ViewOutput string `yaml:"viewOutput"`
	TypeScript bool   `yaml:"typescript"`
}

// FeaturesConfig 功能配置
type FeaturesConfig struct {
	List        bool `yaml:"list"`
	Add         bool `yaml:"add"`
	Edit        bool `yaml:"edit"`
	Delete      bool `yaml:"delete"`
	View        bool `yaml:"view"`
	Export      bool `yaml:"export"`
	Import      bool `yaml:"import"`
	BatchDelete bool `yaml:"batchDelete"`
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig 保存配置文件
func SaveConfig(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Driver: "mysql",
		},
		Generator: GeneratorConfig{
			Backend: BackendConfig{
				Enabled:   true,
				Output:    "./server",
				Package:   "github.com/gfrd/server",
				LayerMode: "simple",
				WithTest:  true,
				WithDoc:   true,
			},
			Frontend: FrontendConfig{
				Enabled:    true,
				Output:     "./web/src",
				APIOutput:  "api",
				ViewOutput: "views",
				TypeScript: true,
			},
			Features: FeaturesConfig{
				List:        true,
				Add:         true,
				Edit:        true,
				Delete:      true,
				View:        true,
				Export:      false,
				Import:      false,
				BatchDelete: true,
			},
		},
	}
}
