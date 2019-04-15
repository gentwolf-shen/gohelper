package memcache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"

	mc "github.com/bradfitz/gomemcache/memcache"
	"github.com/gentwolf-shen/gohelper/config"
)

var (
	client            *mc.Client
	keyPrefix         = "mem_"
	defaultExpiration = int32(60)
)

func InitFromConfig(cfg config.CacheConfig) {
	Init(cfg.Prefix, cfg.Expiration, cfg.Hosts)
}

func Init(prefix string, expiration int32, hosts []string) {
	keyPrefix = prefix
	defaultExpiration = expiration
	Connect(hosts...)
}

func Connect(hosts ...string) {
	client = mc.New(hosts...)
	client.Timeout = 200 * time.Millisecond
}

func GetConn() *mc.Client {
	return client
}

func SetByte(name string, bytes []byte, args ...int32) error {
	return client.Set(&mc.Item{Key: keyPrefix + name, Value: bytes, Expiration: getExpire(args...)})
}

func GetByte(name string) ([]byte, error) {
	item, err := client.Get(keyPrefix + name)
	if err != nil {
		return nil, err
	}

	return item.Value, nil
}

func Set(name string, value interface{}, args ...int32) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(&mc.Item{Key: keyPrefix + name, Value: b, Expiration: getExpire(args...)})
}

func Get(name string, value interface{}) error {
	item, err := client.Get(keyPrefix + name)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, value)
}

func GetObject(key string, obj interface{}) error {
	item, err := client.Get(keyPrefix + key)
	if err != nil {
		return err
	}

	return decode(item.Value, obj)
}

func SetObject(key string, obj interface{}, args ...int32) error {
	b, err := encode(obj)
	if err != nil {
		return err
	}

	return client.Set(&mc.Item{Key: keyPrefix + key, Value: b, Expiration: getExpire(args...)})
}

func Delete(key string) error {
	return client.Delete(keyPrefix + key)
}

func Increment(key string, delta uint64) (uint64, error) {
	return client.Increment(keyPrefix+key, delta)
}

func Decrement(key string, delta uint64) (uint64, error) {
	return client.Decrement(keyPrefix+key, delta)
}

func encode(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(b []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(v)
}

func getExpire(args ...int32) int32 {
	expire := defaultExpiration
	if len(args) > 0 {
		expire = args[0]
	}

	return expire
}
