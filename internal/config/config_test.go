package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	dir, err := os.Getwd()
	fmt.Println(dir)
	path := filepath.Join("testdata", "config_test.yml")
	cfg, err := Load(path)
	if err != nil {
		t.Fail()
	}
	expected := Config{
		Db: DB{
			Name:           "wallet_test",
			Host:           "localhost",
			Port:           "5432",
			User:           "test_user",
			Password:       "test_user_pass",
			MaxConnections: 10,
		},
		HttpPort: "1234",
	}

	assert.Equal(t, expected, *cfg)
}

func TestReadFromEnvVariables(t *testing.T) {
	// read only from env variable
	_ = os.Setenv("DB_USER", "real_user")
	_ = os.Setenv("DB_PASSWORD", "7239023")
	_ = os.Setenv("DB_MAX_CONNECTIONS", "50")
	_ = os.Setenv("HTTP_PORT", "5001")
	cfg, err := Load("")
	if err != nil {
		t.Fail()
	}
	expected := Config{
		Db: DB{
			Name:           "",
			Host:           "",
			Port:           "",
			User:           "real_user",
			Password:       "7239023",
			MaxConnections: 50,
		},
		HttpPort: "5001",
	}
	assert.Equal(t, expected, *cfg)
}

func TestReadFromYmlFileAndEnvVariables(t *testing.T) {
	_ = os.Setenv("DB_USER", "real_user")
	_ = os.Setenv("DB_PASSWORD", "7239023")
	_ = os.Setenv("DB_MAX_CONNECTIONS", "50")
	_ = os.Setenv("HTTP_PORT", "5001")
	path := filepath.Join("testdata", "config_test.yml")
	cfg, err := Load(path)
	if err != nil {
		t.Fail()
	}
	expected := Config{
		Db: DB{
			Name:           "wallet_test",
			Host:           "localhost",
			Port:           "5432",
			User:           "real_user",
			Password:       "7239023",
			MaxConnections: 50,
		},
		HttpPort: "5001",
	}
	assert.Equal(t, expected, *cfg)
}
