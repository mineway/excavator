package manager

import (
	"github.com/ermos/annotation/parser"
	"github.com/ermos/httpchecker"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler struct {}

type Manager struct {
	HTTP 	struct {
		Method 	string
		RequestURI 	string
	}
	Param		map[string]interface{}
	Query 		map[string]interface{}
	Payload 	map[string]interface{}
	annotation parser.API
}

// New initialize a new Manager object with default configuration
func New (r *http.Request, a parser.API, ps httprouter.Params) (m *Manager, status int, err error) {
	m = &Manager{
		annotation: a,
	}

	params := make(map[string]string)
	for _, p := range ps {
		params[p.Key] = p.Value
	}

	res, err := httpchecker.Check(r, a, params)
	if err != nil {
		return m, res.Status, err
	}

	m.Param = res.Params
	m.Query = res.Queries
	m.Payload = res.Payload

	m.HTTP.Method = r.Method
	m.HTTP.RequestURI = r.RequestURI

	return
}