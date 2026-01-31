// 文件路径: support/config.go
package support

import (
	"fmt"
	"strings"

	"Unipay/src/exceptions"
)

// Config 配置管理类
type Config struct {
	config map[string]interface{}
}

// NewConfig 创建新配置实例
func NewConfig(config map[string]interface{}) *Config {
	if config == nil {
		config = make(map[string]interface{})
	}
	return &Config{
		config: config,
	}
}

// Get 获取配置值
// key: 配置键，支持点号分隔的多维访问，如 "alipay.app_id"
// defaultValue: 可选默认值
func (c *Config) Get(key string, defaultValue ...interface{}) interface{} {
	// 如果key为空，返回整个配置
	if key == "" {
		return c.config
	}

	// 分割键
	keys := strings.Split(key, ".")
	config := c.config

	// 逐层访问
	for i, segment := range keys {
		value, ok := config[segment]
		if !ok {
			// 未找到，返回默认值
			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return nil
		}

		// 如果是最后一层，返回值
		if i == len(keys)-1 {
			return value
		}

		// 继续深入下一层
		subMap, ok := value.(map[string]interface{})
		if !ok {
			// 下一层不是map，无法继续访问
			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return nil
		}
		config = subMap
	}

	return config
}

// GetString 获取字符串配置
func (c *Config) GetString(key string, defaultValue ...string) string {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string, defaultValue ...int) int {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	default:
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
}

// GetFloat 获取浮点数配置
func (c *Config) GetFloat(key string, defaultValue ...float64) float64 {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0.0
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0.0
	}
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		lower := strings.ToLower(v)
		return lower == "true" || lower == "1" || lower == "yes" || lower == "on"
	case int, int8, int16, int32, int64:
		// 非零值视为true
		return fmt.Sprintf("%v", v) != "0"
	case float32, float64:
		// 非零值视为true
		return fmt.Sprintf("%v", v) != "0"
	default:
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
}

// GetMap 获取map配置
func (c *Config) GetMap(key string, defaultValue ...map[string]interface{}) map[string]interface{} {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	if m, ok := value.(map[string]interface{}); ok {
		return m
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// GetSlice 获取切片配置
func (c *Config) GetSlice(key string, defaultValue ...[]interface{}) []interface{} {
	value := c.Get(key)
	if value == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	if slice, ok := value.([]interface{}); ok {
		return slice
	}

	if slice, ok := value.([]string); ok {
		result := make([]interface{}, len(slice))
		for i, v := range slice {
			result[i] = v
		}
		return result
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// Set 设置配置值
// key: 配置键，支持点号分隔，最多支持三维
// value: 配置值
func (c *Config) Set(key string, value interface{}) error {
	if key == "" {
		return exceptions.NewInvalidArgumentException("Invalid config key.")
	}

	keys := strings.Split(key, ".")

	// 只支持最多三维数组（与PHP版本保持一致）
	if len(keys) > 3 {
		return exceptions.NewInvalidArgumentException("Invalid config key.")
	}

	switch len(keys) {
	case 1:
		c.config[keys[0]] = value
	case 2:
		if c.config[keys[0]] == nil {
			c.config[keys[0]] = make(map[string]interface{})
		}
		if subMap, ok := c.config[keys[0]].(map[string]interface{}); ok {
			subMap[keys[1]] = value
		} else {
			return exceptions.NewInvalidArgumentException(
				fmt.Sprintf("Cannot set value for key '%s': parent is not a map", key),
			)
		}
	case 3:
		// 确保第一层存在
		if c.config[keys[0]] == nil {
			c.config[keys[0]] = make(map[string]interface{})
		}

		firstMap, ok := c.config[keys[0]].(map[string]interface{})
		if !ok {
			return exceptions.NewInvalidArgumentException(
				fmt.Sprintf("Cannot set value for key '%s': first level is not a map", key),
			)
		}

		// 确保第二层存在
		if firstMap[keys[1]] == nil {
			firstMap[keys[1]] = make(map[string]interface{})
		}

		secondMap, ok := firstMap[keys[1]].(map[string]interface{})
		if !ok {
			return exceptions.NewInvalidArgumentException(
				fmt.Sprintf("Cannot set value for key '%s': second level is not a map", key),
			)
		}

		secondMap[keys[2]] = value
	}

	return nil
}

// Has 检查配置是否存在
func (c *Config) Has(key string) bool {
	return c.Get(key) != nil
}

// Delete 删除配置
func (c *Config) Delete(key string) error {
	return c.Set(key, nil)
}

// All 获取所有配置
func (c *Config) All() map[string]interface{} {
	return c.config
}

// Merge 合并配置
func (c *Config) Merge(config map[string]interface{}) {
	c.mergeMaps(c.config, config)
}

// mergeMaps 递归合并maps
func (c *Config) mergeMaps(dest, src map[string]interface{}) {
	for key, srcVal := range src {
		if destVal, ok := dest[key]; ok {
			// 如果dest中已存在该key
			if srcMap, ok := srcVal.(map[string]interface{}); ok {
				// src的值是map
				if destMap, ok := destVal.(map[string]interface{}); ok {
					// dest的值也是map，递归合并
					c.mergeMaps(destMap, srcMap)
				} else {
					// dest的值不是map，直接替换
					dest[key] = srcVal
				}
			} else {
				// src的值不是map，直接替换
				dest[key] = srcVal
			}
		} else {
			// dest中不存在该key，直接添加
			dest[key] = srcVal
		}
	}
}

// Clone 克隆配置
func (c *Config) Clone() *Config {
	newConfig := make(map[string]interface{})
	c.cloneMap(c.config, newConfig)
	return NewConfig(newConfig)
}

// cloneMap 递归克隆map
func (c *Config) cloneMap(src, dest map[string]interface{}) {
	for key, value := range src {
		if subMap, ok := value.(map[string]interface{}); ok {
			newSubMap := make(map[string]interface{})
			c.cloneMap(subMap, newSubMap)
			dest[key] = newSubMap
		} else if slice, ok := value.([]interface{}); ok {
			// 克隆slice
			newSlice := make([]interface{}, len(slice))
			copy(newSlice, slice)
			dest[key] = newSlice
		} else {
			dest[key] = value
		}
	}
}

// Filter 过滤配置（只保留指定前缀的配置）
func (c *Config) Filter(prefix string) *Config {
	filtered := make(map[string]interface{})
	c.filterMap(c.config, filtered, prefix, "")
	return NewConfig(filtered)
}

// filterMap 递归过滤map
func (c *Config) filterMap(src, dest map[string]interface{}, prefix, currentPath string) {
	for key, value := range src {
		fullPath := key
		if currentPath != "" {
			fullPath = currentPath + "." + key
		}

		if strings.HasPrefix(fullPath, prefix) {
			// 去掉前缀
			newKey := strings.TrimPrefix(fullPath, prefix)
			if newKey != "" && newKey[0] == '.' {
				newKey = newKey[1:]
			}

			if subMap, ok := value.(map[string]interface{}); ok {
				newDest := make(map[string]interface{})
				c.filterMap(subMap, newDest, prefix, fullPath)
				if len(newDest) > 0 {
					dest[newKey] = newDest
				}
			} else {
				dest[newKey] = value
			}
		} else if subMap, ok := value.(map[string]interface{}); ok {
			// 继续深入查找
			c.filterMap(subMap, dest, prefix, fullPath)
		}
	}
}
