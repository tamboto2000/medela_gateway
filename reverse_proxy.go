package medelagateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxy(
	address *url.URL,
	req *http.Request,
	resmod http.RoundTripper,
) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(address)
	p.Director = func(r *http.Request) {
		r.URL = address
		r.Method = req.Method
		r.Header = req.Header
		r.Body = req.Body
		r.GetBody = req.GetBody
		r.ContentLength = req.ContentLength
		r.TransferEncoding = req.TransferEncoding
		r.Close = req.Close
		r.Host = address.Host
		r.Form = req.Form
		r.PostForm = req.PostForm
		r.MultipartForm = req.MultipartForm
		r.Trailer = req.Trailer
	}

	// p.ModifyResponse = resmod
	p.Transport = resmod

	return p
}
