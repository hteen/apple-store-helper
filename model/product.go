package model

type Product struct {
    Title string
    Type string
    Code string
}

var TypeCode = map[string]string{
    "iphone13mini": "D",
    "iphone13": "D",
    
    "iphone13pro": "A",
    "iphone13promax": "A",
}
