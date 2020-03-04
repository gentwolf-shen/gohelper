package gohttp

func Get(url string) *GoHttpRequest {
	return new(MethodGet, url)
}

func Post(url string) *GoHttpRequest {
	return new(MethodPost, url)
}

func Put(url string) *GoHttpRequest {
	return new(MethodPut, url)
}

func Patch(url string) *GoHttpRequest {
	return new(MethodPatch, url)
}

func Delete(url string) *GoHttpRequest {
	return new(MethodDelete, url)
}

func Head(url string) *GoHttpRequest {
	return new(MethodHead, url)
}

func Options(url string) *GoHttpRequest {
	return new(MethodOptions, url)
}
