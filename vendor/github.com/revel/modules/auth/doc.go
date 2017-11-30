/*	A basic user authentication module for Revel

 list of concerns:
- Separating out the interface and driver
- Removing DB/Storage dependency
- UUID as default identifier?
- how to deal with password/secret or generally, method of authorization
- default {views,controllers,routes} for register/login/logut ?
- reset password in most basic ?
- activation (and other features) in a second / more sophisticated driver
- filter for checking that user is authenticated


I think a driver is made up of 2 parts
data prep and data storage
register and password reset are part of data prep
as is the auth hash method
they don't care how the data is stored

then there is the data store

perhaps each auth user model should instantiate 2 drivers instead of 1?
one for data prep components and one for storage
so the security driver and the storage driver

*/
package auth
