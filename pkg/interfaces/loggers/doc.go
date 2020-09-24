/*
Package loggers is where all the events from the codebase converge. The goal is
to decouple the event's source from the actual handling, as both will need to
change by different reasons. The logger will decide how to deal with each
event, freeing the reporter from such responsibility.

What lives here

Types that require events to be logged should declare an interface that covers
every possible event that can occur within its scope. Those events must be
reported on business terms, and the implementation of such interfaces must live
on this package.

Rules of the road

If you are using syslog, a logging library, or printing to stdout/stderr
directly, rather declare a logger interface and implement it here.
*/
package loggers
