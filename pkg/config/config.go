package config

import "os"

// Config holds application configuration
type Config struct {
	Port       string
	DBPath     string
	UploadDir  string
	FilesDir   string
	ImagesDir  string
	StaticDir  string
	ExportsDir string
	BackupsDir string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "/app/data/app.db"),
		UploadDir:  getEnv("UPLOAD_DIR", "/app/uploads"),
		FilesDir:   getEnv("FILES_DIR", "/app/files"),
		ImagesDir:  getEnv("IMAGES_DIR", "/app/images"),
		StaticDir:  getEnv("STATIC_DIR", "/app/static"),
		ExportsDir: getEnv("EXPORTS_DIR", "/app/exports"),
		BackupsDir: getEnv("BACKUPS_DIR", "/app/backups"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
