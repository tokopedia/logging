package logging

var buildhash string

func appVersion() string {
	return buildhash
}

// SetAppVersion see https://github.com/tokopedia/logging/issues/17
func SetAppVersion(v string) {
	buildhash = v
}
