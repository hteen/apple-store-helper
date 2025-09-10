package model

type Product struct {
	Title string
	Type  string
	Code  string
}

var TypeCode = map[string]string{
	"iphoneair":      "A",
	"iphone17":       "D",
	"iphone17pro":    "A",
	"iphone17promax": "A",
}
