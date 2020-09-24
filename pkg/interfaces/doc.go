/*
Package interfaces is a place where the outside world lies on one side and our
business on the other. Here we transform information from the outside world to
business objects and internal replies to their corresponding outside world
representation.

What lives here

As an intermediary, here is where most translation code dwells. HTTP Handlers
abstract away the WWW and its protocol. Databases speak their own language.
That is allowed at this layer, but not anywhere else. File objects, external
services, anything that affects or is affected by the real world needs to be
separated from our business internals.

Rules of the road

Despite we being able to talk to external services at this layer, we cannot
hardwire any service location. We talk the language, we don't decide who we
talk to.
*/
package interfaces
