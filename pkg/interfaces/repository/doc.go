/*
Package repository defines interactions with different data sources and storage
solutions. Either a file, a memory region or an external service, if data needs
to be loaded or sent away, a repository is the abstraction we need in place.

What lives here

The usecases and domain layers surely define interfaces for repositories that
will provide and store the business objects they operate on. Here is where the
implementations of those repositories reside.

Rules of the road

If any of the words: read, save, store, send, get, update, are on the table,
it's extremely likely that you need a Repository.

The representation chosen by a repository can be different from the actual
domain object. This is perfectly ok, as repositories must be optimized for data
retrieval or storage, or both. However, repository internal objects must NOT
reach the internal layers. They must be converted to business entities first.
*/
package repository
