package platform

// AbstractCrawler all media platforms abstract interface, the specific platform needs to implement it
type AbstractCrawler interface {

	// InitConfig  methods receive some parameters and assigning values to the properties of a structure.
	InitConfig(loginType string)

	// Start method is the main process implemented in the specific platform.
	Start()

	// Search method will use some keywords to find for content on corresponding platforms.
	search()
}

// AbstractLogin extract some login methods into an interface, and specific subclasses can implement the methods.
type AbstractLogin interface {

	// begin method is the main process implemented on the logging in
	begin()

	// loginByQrcode use qrcode logging in the specific platform
	loginByQrcode()

	// loginByQrcode use cookies logging in the specific platform
	loginByCookies()
}
