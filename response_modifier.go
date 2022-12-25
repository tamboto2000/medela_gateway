package medelagateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/google/uuid"
)

type responseModifier struct {
	body        interface{}
	header      http.Header
	statusCode  int
	uuidStr     string
	mergeBody   bool
	mergeHeader bool
	http.RoundTripper
}

func newResponseModifier() *responseModifier {
	return &responseModifier{
		uuidStr:      uuid.NewString(),
		RoundTripper: http.DefaultTransport,
	}
}

func (respMod *responseModifier) RoundTrip(req *http.Request) (*http.Response, error) {
	r, err := respMod.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respMod.statusCode = r.StatusCode
	r.Header.Set("Request-Id", respMod.uuidStr)

	if r.StatusCode > 399 {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		respMod.body = b
		buff := bytes.NewBuffer(b)
		rc := io.NopCloser(buff)
		r.Body = rc

		return r, nil
	}

	if respMod.header != nil {
		if respMod.mergeHeader {
			for k, h := range respMod.header {
				for _, v := range h {
					if !isHeaderCanon(k) {
						r.Header.Add(k, v)
					} else {
						if r.Header.Get(k) == "" {
							r.Header.Set(k, v)
						}
					}
				}
			}

			respMod.header = r.Header
		}
	} else {
		respMod.header = r.Header
	}

	// TODO
	// check response Content-Type header
	if respMod.body != nil {
		// only support JSON merging
		if respMod.mergeBody {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			a := new(interface{})
			if err := json.Unmarshal(b, a); err != nil {
				return nil, err
			}

			respMod.body = mergeData(reflect.ValueOf(respMod.body), reflect.ValueOf(*a))
			b, err = json.Marshal(respMod.body)
			if err != nil {
				return nil, err
			}

			buff := bytes.NewBuffer(b)

			rc := io.NopCloser(buff)
			r.Body = rc
			r.ContentLength = int64(len(b))
			r.Header.Set("Content-Length", strconv.Itoa(len(b)))
		}
	} else {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		a := new(interface{})
		if err := json.Unmarshal(b, a); err != nil {
			return nil, err
		}

		respMod.body = *a
		buff := bytes.NewBuffer(b)
		rc := io.NopCloser(buff)
		r.Body = rc
	}

	return r, nil
}

func mergeData(data, dest reflect.Value) interface{} {
	var res interface{}

	for data.Kind() == reflect.Interface {
		data = data.Elem()
	}

	for dest.Kind() == reflect.Interface {
		dest = dest.Elem()
	}

	// merge destination
	switch dest.Kind() {
	// map[string]interface{}
	case reflect.Map:
		// data to be merged
		switch data.Kind() {
		case reflect.Map:
			dataIter := data.MapRange()
			for dataIter.Next() {
				destMV := dest.MapIndex(dataIter.Key())
				if destMV.IsValid() {
					dest.SetMapIndex(dataIter.Key(), reflect.ValueOf(mergeData(dataIter.Value(), destMV)))
				} else {
					dest.SetMapIndex(dataIter.Key(), dataIter.Value())
				}
			}

			res = dest.Interface()

		default:
			i := 0
			for {
				key := reflect.ValueOf(fmt.Sprintf("extra_data_%d", i))
				destK := dest.MapIndex(key)
				if destK.IsValid() {
					i++

					continue
				}

				dest.SetMapIndex(key, data)
				break
			}

			res = dest.Interface()
		}

	case reflect.Slice:
		switch data.Kind() {
		case reflect.Slice:
			for i := 0; i < data.Len(); i++ {
				dest = reflect.Append(dest, data.Index(i))
			}

			res = dest.Interface()

		default:
			dest = reflect.Append(dest, data)
			res = dest.Interface()
		}

	default:
		switch data.Kind() {
		case reflect.Slice:
			slice := reflect.ValueOf([]interface{}{})
			slice = reflect.Append(slice, dest)
			for i := 0; i < data.Len(); i++ {
				slice = reflect.Append(slice, data.Index(i))
			}

			res = slice.Interface()

		default:
			slice := reflect.ValueOf([]interface{}{})
			slice = reflect.Append(slice, dest, data)
			res = slice.Interface()
		}
	}

	return res
}
