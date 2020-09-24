/*
Package handlers groups together all the code that talks to the World Wide Web.
Here we can decode/encode JSON requests over HTTP, read URL parameters, check
credentials for authorization, read/modify http headers. That information must
be transformed into business language in order to call the inner layers.

What lives here

Handlers and Presenters. A handler carries along a request, dealing with the
protocol on behalf of the user case. A presenter takes a business object and
constructs an adequate representation for the outside world.

Rules of the road

Incoming requests must be mapped to business objects. Those objects will be fed
to the inner layers. Outgoing responses must come on business terms and must be
presented as the contracts demands.

Under no circumstances an external objects enters, unfiltered, to the inner
layers, nor inner layers have a word on outside world representation.
*/
package handlers
