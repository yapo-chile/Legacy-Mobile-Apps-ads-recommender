/*
Package infrastructure is where all the outside world communications resides.
We don't speak application language here. In fact, we don't care what the
application is, at this level.

Infrastructure layer abstracts away the outside world and allows for the
application to be tested in isolation. It also allows to minimize the effort
of replacing one external component by another, as all possible interactions
from anywhere else are located only on the outermost layer. No business logic
must be impacted by such change.

What lives here

Any code that communicates at low level with the outside world belongs to the
infrastructure layer. Some examples are: Environment variables, handling files,
database drivers, network connections.

Rules of the road

If you didn't write it and you're talking to it, you need code here.
If you need to know where is it, you need code here.
*/
package infrastructure
