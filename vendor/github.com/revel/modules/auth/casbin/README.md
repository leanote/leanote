modules/casbin
======

Casbin authz is an authorization middleware for [Revel](https://github.com/revel/revel), it's based on [https://github.com/casbin/casbin](https://github.com/casbin/casbin).

## Simple Example

```Go
package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/casbin/casbin"
	"github.com/revel/revel"
	"github.com/revel/modules/auth/casbin"
)

var adapter = casbinauth.NewAdapter(params)
var enforcer = casbin.NewEnforcer("authz_model.conf", adapter)
var casbinModule = casbinauth.NewCasbinModule(enforcer)

var testFilters = []revel.Filter{
	casbinModule.AuthzFilter,
	func(c *revel.Controller, fc []revel.Filter) {
		c.RenderHTML("OK.")
	},
}

func main() {
	r, _ := http.NewRequest("GET", "/dataset1/resource1", nil)
    	r.SetBasicAuth("alice", "123")
    	w := httptest.NewRecorder()
    	context := revel.NewGoContext(nil)
    	context.Request.SetRequest(r)
    	context.Response.SetResponse(w)
    	c := revel.NewController(context)
    
    	testFilters[0](c, testFilters)
}
```

## Documentation

The authorization determines a request based on ``{subject, object, action}``, which means what ``subject`` can perform what ``action`` on what ``object``. In this plugin, the meanings are:

1. ``subject``: the logged-on user name
2. ``object``: the URL path for the web resource like "dataset1/item1"
3. ``action``: HTTP method like GET, POST, PUT, DELETE, or the high-level actions you defined like "read-file", "write-blog"


For how to write authorization policy and other details, please refer to [the Casbin's documentation](https://github.com/casbin/casbin).
