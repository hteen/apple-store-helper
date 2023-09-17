package model

type Product struct {
	Title string
	Type  string
	Code  string
}

var TypeCode = map[string]string{
	"iphone15":       "A",
	"iphone15pro":    "A",
	"iphone15promax": "A",
}
