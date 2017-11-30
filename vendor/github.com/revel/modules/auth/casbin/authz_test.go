package casbinauthz

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/revel/revel"
	"github.com/revel/modules/orm/gorm/app"
)

var adapter = NewAdapter(DefaultDbParams())
var enforcer = casbin.NewEnforcer("authz_model.conf", adapter)
var casbinModule = NewCasbinModule(enforcer)

var testFilters = []revel.Filter{
	casbinModule.AuthzFilter,
	func(c *revel.Controller, fc []revel.Filter) {
		c.RenderHTML("OK.")
	},
}

func DefaultDbParams() gormdb.DbInfo {
	params := gormdb.DbInfo{}
	params.DbDriver = "mysql"
	params.DbHost = "(localhost:3306)"
	params.DbUser = "root"
	params.DbPassword = ""
	params.DbName = "casbin"
	return params
}

func testRequest(t *testing.T, user string, path string, method string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	r.SetBasicAuth(user, "123")
	w := httptest.NewRecorder()
	context := revel.NewGoContext(nil)
	context.Request.SetRequest(r)
	context.Response.SetResponse(w)
	c := revel.NewController(context)

	testFilters[0](c, testFilters)

	if c.Response.Status != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", user, path, method, c.Response.Status, code)
	}
}

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	myRes := e.GetPolicy()
	log.Print("Policy: ", myRes)

	if !util.Array2DEquals(res, myRes) {
		t.Error("Policy: ", myRes, ", supposed to be ", res)
	}
}

func initPolicy(t *testing.T) {
	// Because the DB is empty at first,
	// so we need to load the policy from the file adapter (.CSV) first.
	e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")

	a := NewAdapter(DefaultDbParams())
	// This is a trick to save the current policy to the DB.
	// We can't call e.SavePolicy() because the adapter in the enforcer is still the file adapter.
	// The current policy means the policy in the Casbin enforcer (aka in memory).
	err := a.SavePolicy(e.GetModel())
	if err != nil {
		panic(err)
	}
}

func TestBasic(t *testing.T) {
	// Initialize some policy in DB.
	initPolicy(t)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.

	testRequest(t, "alice", "/dataset1/resource1", "GET", 200)
	testRequest(t, "alice", "/dataset1/resource1", "POST", 200)
	testRequest(t, "alice", "/dataset1/resource2", "GET", 200)
	testRequest(t, "alice", "/dataset1/resource2", "POST", 403)
}

func TestPathWildcard(t *testing.T) {
	// Initialize some policy in DB.
	initPolicy(t)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.

	testRequest(t, "bob", "/dataset2/resource1", "GET", 200)
	testRequest(t, "bob", "/dataset2/resource1", "POST", 200)
	testRequest(t, "bob", "/dataset2/resource1", "DELETE", 200)
	testRequest(t, "bob", "/dataset2/resource2", "GET", 200)
	testRequest(t, "bob", "/dataset2/resource2", "POST", 403)
	testRequest(t, "bob", "/dataset2/resource2", "DELETE", 403)

	testRequest(t, "bob", "/dataset2/folder1/item1", "GET", 403)
	testRequest(t, "bob", "/dataset2/folder1/item1", "POST", 200)
	testRequest(t, "bob", "/dataset2/folder1/item1", "DELETE", 403)
	testRequest(t, "bob", "/dataset2/folder1/item2", "GET", 403)
	testRequest(t, "bob", "/dataset2/folder1/item2", "POST", 200)
	testRequest(t, "bob", "/dataset2/folder1/item2", "DELETE", 403)
}

func TestRBAC(t *testing.T) {
	// Initialize some policy in DB.
	initPolicy(t)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.

	// cathy can access all /dataset1/* resources via all methods because it has the dataset1_admin role.
	testRequest(t, "cathy", "/dataset1/item", "GET", 200)
	testRequest(t, "cathy", "/dataset1/item", "POST", 200)
	testRequest(t, "cathy", "/dataset1/item", "DELETE", 200)
	testRequest(t, "cathy", "/dataset2/item", "GET", 403)
	testRequest(t, "cathy", "/dataset2/item", "POST", 403)
	testRequest(t, "cathy", "/dataset2/item", "DELETE", 403)
}
