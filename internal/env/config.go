package env

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	Env      string `yaml:"baseenv"`
	Postgres struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		DB   string `yaml:"database"`
		User string `yaml:"user"`
		Pwd  string `yaml:"password"`
	} `yaml:"postgres"`
	Mysql struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		DB   string `yaml:"database"`
		User string `yaml:"user"`
		Pwd  string `yaml:"password"`
	} `yaml:"mysql"`
	Mongo struct {
	} `yaml:"mongo"`
	Meta struct {
		HTTPPort int `yaml:"http_port"`
		GRPCPort int `yaml:"grpc_port"`
	} `yaml:"meta"`
	Logger struct {
		Level zapcore.Level `yaml:"level"`
	} `yaml:"logger"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		DB   int    `yaml:"db"`
	} `yaml:"redis"`
	Aws struct {
		AccessID     string `yaml:"access_id"`
		AccessSecret string `yaml:"access_secret"`
		S3Bucket     string `yaml:"s3_bucket"`
		S3Folder     string `yaml:"s3_folder"`
		S3Region     string `yaml:"s3_region"`
	} `yaml:"aws"`
	Jaeger struct {
		TracerName         string `yaml:"tracer_name"`
		ProjectEnv         string `yaml:"project_env"`
		LocalAgentHostPort string `yaml:"local_agent_host_port"`
		SamplingServerURL  string `yaml:"sampling_server_url"`
	} `yaml:"jaeger"`
	// EventKafka struct {
	// 	Addresses []string `yaml:"addresses"`
	// 	User      string   `yaml:"user"`
	// 	Password  string   `yaml:"password"`
	// 	Enable    bool     `yaml:"enable"`
	// 	Topics    struct {
	// 		EventStatusTopic string `yaml:"event_status_topic"`
	// 		EventCreateTopic string `yaml:"event_create_topic"`
	// 	} `yaml:"topics"`
	// 	Groups struct {
	// 		EventStatusTopicGroup string `yaml:"event_status_topic_group"`
	// 	} `yaml:"groups"`
	// } `yaml:"event_kafka"`
	// RateLimit struct {
	// 	PostUserTTL   int `yaml:"post_user_ttl"`
	// 	PostUserLimit int `yaml:"post_user_limit"`
	// } `yaml:"rate_limit"`
}

func provideConf(path string) (conf *Conf, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config.yaml err: [%w]", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	configBs, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read config file err; [%w]", err)
	}

	conf = new(Conf)
	err = yaml.Unmarshal(configBs, conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal conf str into conf err:%v", err)
	}

	conf.SetConfValueFromEnv()
	if conf.Logger.Level != zapcore.InfoLevel {
		fmt.Printf("got conf: [%+v]\n", conf)
	}
	return conf, nil
}

func (c *Conf) SetConfValueFromEnv() {
	if c.Redis.Host == "" {
		c.Redis.Host = os.Getenv("redis_host")
	}
	if c.Redis.Port == 0 {
		port, _ := strconv.ParseInt(os.Getenv("redis_port"), 10, 0)
		c.Redis.Port = int(port)
	}
	if c.Redis.DB == 0 {
		redisDB, _ := strconv.ParseInt(os.Getenv("redis_db"), 10, 0)
		c.Redis.DB = int(redisDB)
	}
	if c.Aws.AccessID == "" {
		c.Aws.AccessID = os.Getenv("aws_access_id")
	}
	if c.Aws.AccessSecret == "" {
		c.Aws.AccessSecret = os.Getenv("aws_access_secret")
	}
}
