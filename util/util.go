package util

// 复制map
func CopyMap(m map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// map多加一层
func WrapMap(m map[string]string) map[string][]string {
	nm := make(map[string][]string)
	for k, v := range m {
		nm[k] = []string{v}
	}
	return nm
}
