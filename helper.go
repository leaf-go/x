package x

type H map[string]interface{}

func (h H) Get(key string) (interface{}, bool) {
	val, exists := h[key]
	return val, exists
}
