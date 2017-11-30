modules/gorm
===============

[Gorm](http://jinzhu.me/gorm) module

## Activation
```ini
module.gorm = github.com/revel/modules/orm/gorm
```

## Drivers

* sqlite3
* postgres
* mysql

## Configuration file

```ini
# Database config
db.autoinit=true # default=true
db.driver=sqlite # mysql, postgres, sqlite3
db.host=localhost  # Use db.host /tmp/app.db is your driver is sqlite
db.user=dbuser
db.name=dbname
db.password=dbpassword

```

## Example usage with transactions
```go
package controllers

import (
    "github.com/revel/revel"
    gormc "github.com/revel/modules/gorm/orm/app/controllers"
)

type App struct {
    gormc.TxnController
}

type Toy struct {
    Name string
}

func (c App) Index() revel.Result {
    c.Txn.LogMode(true)
    c.Txn.AutoMigrate(&Toy{})
    c.Txn.Save(&Toy{Name: "Fidget spinner"})

    return c.Render()
}
```

## Example usage without transactions
```go
package controllers

import (
    "github.com/revel/revel"
    gormc "github.com/revel/modules/gorm/orm/app/controllers"
)

type App struct {
    gormc.Controller
}

type Toy struct {
    Name string
}

func (c App) Index() revel.Result {
    c.DB.LogMode(true)
    c.DB.AutoMigrate(&Toy{})
    c.DB.Save(&Toy{Name: "Fidget spinner"})

    return c.Render()
}
```
