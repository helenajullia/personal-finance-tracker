# Architecture note

## The shape

The service is a single Go program that listens on port 8080 and answers HTTP requests. It is split into three parts.

- `cmd/app` is the entry point. It builds a `http.ServeMux`, registers the routes and runs an `http.Server` with read, write and idle timeouts. It shuts down cleanly on Ctrl-C or on the signal Docker sends when stopping the container.
- `internal/money` holds `Amount`, the type used for every sum of money. It parses text such as `-12.50`, prints it back, and adds two amounts.
- `internal/store` holds the transactions in memory and hands out ids.
- `internal/version` reports the app name, version, commit and build time, which the Dockerfile stamps into the binary at build time.

A request enters through the mux, is routed by method and path to a handler, and the handler reads or writes the store. A request that adds a transaction carries `amount` as a JSON string. The handler passes it to `money.Parse`, which rejects anything that is not a decimal with exactly two places, and then calls `store.Add`. A listing request calls `store.List`, which returns the transactions newest first, ties broken by ascending id. Errors leave the handler as `{"error":"<code>","message":"<explanation>"}` with the matching status code.

Today only the four routes shared by every project are wired: `/health`, `/version`, `/` and `POST /reset`. The money type and the store exist and are tested, but are not yet connected to any route.

## The data

A transaction has an integer id, a date, an amount, a description and a category. Amounts are `int64` counts of cents, so `-12.50` is stored as `-1250`.

Everything lives in a slice inside the store, guarded by a mutex because the HTTP server handles each request in its own goroutine. Nothing is written to disk, so all data is lost when the process stops. `POST /reset` empties the slice and restarts the ids at 1.

## Three decisions

**Amounts are integer cents, not floats.** The alternative was `float64`, which cannot represent `0.10` exactly, so `0.10 + 0.20` becomes `0.30000000000000004` and a report would be wrong by fractions of a cent. A decimal library was the other option, but it would be a dependency for something an `int64` already does. The cost is that parsing and printing must be written by hand, which is why `money.Parse` accepts only two decimal places and rejects everything else, including a decimal comma.

**Data is kept in memory, not in a database.** The project description does not require data to survive a restart. SQLite would give persistence and queries, but needs a cgo-free driver or a C toolchain in the build, a schema, and a way to make `/reset` clear it. A slice behind a mutex is enough for the filters and reports required, and keeps the Docker image small. The cost is that a restart loses everything.

**Categorisation rules are an ordered list where the first match wins.** The alternative was giving each rule a priority number and picking the highest. That would need a tie-break rule of its own and make behaviour harder to predict. The project states that rule order is the user's and must not be reordered, and a plain list makes that visible: the user reads the rules top to bottom and the first that matches decides.

## What is not built yet

- `POST /transactions` and `GET /transactions`, with the filters `from`, `to`, `category` and `q`.
- `POST /transactions/import`: CSV parsing, per-row rejection with a line number, and detecting a file that was already imported.
- `PUT /rules` and applying rules to a transaction without a category.
- `PUT /budgets` and the `over_budget` list.
- `GET /reports/monthly`, with exact totals by category and the 12-month limit.
- `POST /reset` still does nothing; it must call `store.Reset` once the store is wired in.
- Request validation and the error responses for 400, 404, 409 and 422.
