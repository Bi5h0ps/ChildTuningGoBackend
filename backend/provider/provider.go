package provider

var HttpClientProvider HttpClient

func init() {
	HttpClientProvider = NewHttpClient()
}
