package services

// runtime flags controlled by UI/user settings
var preferCache bool

func SetPreferCache(v bool) {
	preferCache = v
}

func PreferCacheEnabled() bool {
	return preferCache
}
