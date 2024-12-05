package env

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/joho/godotenv"
)

type MasterEnv struct {
	ServerHost                  string `env:"server_host"`
	ServerPort                  string `env:"server_port"`
	MessageQueueHost            string `env:"mq_sub_host"`
	MessageQueuePort            string `env:"mq_sub_port"`
	MessageQueueTransportPubKey string `env:"mq_transport_pub_key"`
	MessageQueueTransportPemKey string `env:"mq_transport_pem_key"`
	BoltDBPath                  string `env:"bolt_db_path"`
}

type NodeEnv struct {
	MessageQueueHost            string `env:"mq_sub_host"`
	MessageQueuePort            string `env:"mq_sub_port"`
	MessageQueueTransportPubKey string `env:"mq_transport_pub_key"`
	MessageQueueInviteCode      string `env:"mq_invite_code"`
}

// If not nil, it means env has been already loaded.
var nodeEnv *NodeEnv
var masterEnv *MasterEnv

// Loads environment variables for master server.
// If error is not nil, it means neither .env on the working directory or ${CFG_DIR}/master.env could be lookup.
func Master() (MasterEnv, error) {
	if masterEnv != nil {
		return *masterEnv, nil
	}

	err := loadEnv("master.env")
	if err != nil {
		return MasterEnv{}, err
	}

	masterEnv = fromEnv(MasterEnv{})
	return *masterEnv, err
}

// Loads environment variables for node agent.
// If error is not nil, it means neither .env on the working directory or ${CFG_DIR}/node.env could be lookup.
func Node() (NodeEnv, error) {
	if nodeEnv != nil {
		return *nodeEnv, nil
	}

	err := loadEnv("node.env")
	if err != nil {
		return NodeEnv{}, err
	}

	nodeEnv = fromEnv(NodeEnv{})
	return *nodeEnv, err
}

func Save[T MasterEnv | NodeEnv](
	env T,
) error {
	d, err := config.Dir()
	if err != nil {
		return err
	}

	if reflect.TypeOf(env) == reflect.TypeFor[NodeEnv]() {
		return setEnv(env, filepath.Join(d, "node.env"))
	} else {
		return setEnv(env, filepath.Join(d, "master.env"))
	}
}

func loadEnv(
	filename string,
) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	d, err := config.Dir()
	if err == nil {
		err = godotenv.Load(filepath.Join(d, filename))
	}

	if err != nil {
		return err
	}

	return nil
}

func fromEnv[T MasterEnv | NodeEnv](
	env T,
) *T {
	v := reflect.ValueOf(env)
	addr := reflect.ValueOf(&env).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		k := f.Tag.Get("env")
		addr.Field(i).SetString(os.Getenv(k))
	}

	return &env
}

func setEnv[T MasterEnv | NodeEnv](
	env T,
	path string,
) error {

	m := map[string]string{}
	v := reflect.ValueOf(env)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		k := f.Tag.Get("env")
		v := v.Field(i)

		m[k] = v.String()
	}

	return godotenv.Write(m, path)
}
