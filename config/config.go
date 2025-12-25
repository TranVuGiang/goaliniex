package config

import (
	"fmt"
	"os"
)

type Config struct {
	BaseURL        string
	PartnerCode    string
	SecretKey      string
	PrivateKeyPath string
	PrivateKey     string
	PublicKeyPath  string
	PublicKey      string
}

func LoadFromEnv() (Config, error) {
	cfg := Config{
		BaseURL:        os.Getenv("ALIX_BASE_URL"),
		PartnerCode:    os.Getenv("ALIX_PARTNER_CODE"),
		SecretKey:      os.Getenv("ALIX_SECRET_KEY"),
		PrivateKeyPath: os.Getenv("ALIX_PRIVATE_KEY_PATH"),
		PublicKeyPath:  os.Getenv("ALIX_PUBLIC_KEY_PATH"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	privateKey, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load private key from %s: %w", cfg.PrivateKeyPath, err)
	}
	cfg.PrivateKey = string(privateKey)

	if cfg.PublicKeyPath != "" {
		publicKey, err := os.ReadFile(cfg.PublicKeyPath)
		if err != nil {
			return Config{}, fmt.Errorf("failed to load public key from %s: %w", cfg.PublicKeyPath, err)
		}
		cfg.PublicKey = string(publicKey)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	missing := []string{}

	if c.BaseURL == "" {
		missing = append(missing, "ALIX_BASE_URL")
	}
	if c.PartnerCode == "" {
		missing = append(missing, "ALIX_PARTNER_CODE")
	}
	if c.SecretKey == "" {
		missing = append(missing, "ALIX_SECRET_KEY")
	}
	if c.PrivateKeyPath == "" {
		missing = append(missing, "ALIX_PRIVATE_KEY_PATH")
	}

	if len(missing) > 0 {
		return fmt.Errorf(
			"aliniex config missing required env: %v",
			missing,
		)
	}

	return nil
}
