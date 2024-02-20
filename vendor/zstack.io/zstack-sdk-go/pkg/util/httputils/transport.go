package httputils

import "net/http"

type transport struct {
	check func(*http.Request) (func(resp *http.Response), error)
	ts    *http.Transport
}

func (self *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var respCheck func(resp *http.Response) = nil
	var err error
	if self.check != nil {
		respCheck, err = self.check(req)
		if err != nil {
			return nil, err
		}
	}
	resp, err := self.ts.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if respCheck != nil {
		respCheck(resp)
	}
	return resp, nil
}

func GetCheckTransport(ts *http.Transport, check func(*http.Request) (func(resp *http.Response), error)) http.RoundTripper {
	ret := &transport{ts: ts, check: check}
	return ret
}
