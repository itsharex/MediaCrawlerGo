package platform

import "net/http"

// AbstractCrawler all media platforms abstract interface, the specific platform needs to implement it
type AbstractCrawler interface {

	// InitConfig  methods receive some parameters and assigning values to the properties of a structure.
	InitConfig(loginType string)

	// Start method is the main process implemented in the specific platform.
	Start()

	// Search method will use some keywords to find for content on corresponding platforms.
	Search()
}

// AbstractLogin extract some login methods into an interface, and specific subclasses can implement the methods.
type AbstractLogin interface {

	// Begin method is the main process implemented on the logging in
	Begin()

	// LoginByQrcode use qrcode logging in the specific platform
	LoginByQrcode()

	// LoginByCookies use cookies logging in the specific platform
	LoginByCookies()

	// CheckLoginState Asynchronous polling check login status.
	CheckLoginState()
}

// AbstractClient all platform client abstract client
type AbstractClient interface {
	// Get http method
	Get(url string, headers map[string]string) (*http.Response, error)

	// Post http method
	Post(url string, body []byte, headers map[string]string) (*http.Response, error)

	// PreHeaders process request headers
	PreHeaders(params ...any) map[string]string
}
