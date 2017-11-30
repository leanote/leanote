# NewRelic
This module is the for [NewRelic](http://newrelic.com) it wraps the Go
http server and inserts a filter to track each requests.

###App.conf
- **server.engine** You must set this to `newrelic` in order to use this server engine
- **server.newrelic.license** The license key for NewRelic
- **server.newrelic.addfilter** Defauls `true` When true a filter is inserted
into `revel.Filters` at position 2. This filter logs every request.

###Other Notes
There is a controller called RelicController which contains a single method
to return a `newrelic.Application` object. If you want to access this object
you can extend that class.
 
