package test

import (
	"bytes"
	"github.com/vuuvv/mapstructure"
	"github.com/vuuvv/viper"
	"os"
	"testing"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type DbManager struct {
	Master   *DbConfig
	Replicas []*DbConfig
	Max      int
	Idle     int
}

func TestYaml(t *testing.T) {
	config := `
gboot:
    db:
        max: 20
        idle: 2000
        master:
            host: 192.168.1.50
            port: 3306
            username: root
            password: 123456
        replicas:
            - host: 192.168.1.51
              port: 3306
              username: root
              password: 123456
            - host: 192.168.1.51
              port: 3306
              username: root
              password: 123456
            - host: 192.168.1.51
              port: 3306
              username: root
              password: 123456
`
	_ = os.Setenv("GBOOT_DB_MAX", "300")
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(config)))
	if err != nil {
		t.Fatalf("read config error: %#v", err)
	}

	result := &DbManager{}
	err = viper.UnmarshalKey("gboot.db", &result, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.SystemEnvironmentHookFunc("gboot.db"),
		)))

	if err != nil {
		t.Fatalf("unmarshar error: %#v", err)
	}

	if result.Max != 300 {
		t.Fatalf("wanted: %d, actual: %d", 300, result.Max)
	}

	//max := viper.Get("gboot.db.max")
	//values := viper.AllSettings()
	//t.Log(max, values)
}
