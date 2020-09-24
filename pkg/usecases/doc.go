/*
Package usecases is where the interaction between business objects takes place.
This is the glue code that get things done by coordinating domain objects,
abiding by their rules.

What lives here

Interactors are the main kind of objects to be found here. An interactor may
hold reference to one or many repositories from where to fetch resources, then
use those resources for query or modification and ensure they are properly
written back.

The suggested implementation is having an interface with the usecase name and
a single function. Then, implementing this interface with a suitable interactor
that can be mocked on the test code of the outer layers.

Rules of the road

The only package we can (and should) import from is the domain layer. Every
data transformation of domain entities must done via their public API. If that
is not possible then it can mean two things. Either the domain dissalows it and
cannot be performed, or the domain is incomplete and needs extension.
*/
package usecases
