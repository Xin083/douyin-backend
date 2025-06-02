package interf

import "time"

type YmlConfigInterface interface {
	ConfigFileChangeListen()
	Clone(fileName string) YmlConfigInterface
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetBool(keyName string) bool
	GetInt(keyName string) int
	GetInt32(keyName string) int32
	GetInt64(keyName string) int64
	GetFloat32(keyName string) float32
	GetFloat(keyName string) float64
	GetDuration(keyName string) time.Duration
	GetStringSlice(keyName string) []string
	GetStringMap(keyName string) map[string]interface{}
	GetStringMapString(keyName string) map[string]string
	GetStringMapBool(keyName string) map[string]bool
	GetStringMapInt(keyName string) map[string]int
	GetStringMapInt32(keyName string) map[string]int32
	GetStringMapInt64(keyName string) map[string]int64
	GetStringMapFloat32(keyName string) map[string]float32
	GetStringMapFloat64(keyName string) map[string]float64
}
