# Personal Finance Tracker

A web service for tracking income and expenses. Transactions come in from a bank CSV export or are added by hand, get sorted into categories by rules the user defines, and come back out as a monthly report with per-category budgets.

- **Student:** Pastore Helena-Jullia, group DMWD
- **Project:** 2, Personal Finance Tracker
- **Language:** Go

## What it does

The service stores transactions and works out spending by month and by category. Amounts are kept as integer cents, never as floating-point numbers, so totals are exact: `0.10` plus `0.20` is always `0.30`.

Categories are assigned by an ordered list of rules (a text to look for in the description, and the category to give). The first matching rule wins, and a transaction that matches nothing is left uncategorised. A bank CSV can be imported in one request. A malformed row is rejected with its line number and the reason, and the remaining rows are still imported.

The service has four routes shared by every project (`/health`, `/version`, `/` and `POST /reset`) and the finance routes below.

Current state: the four shared routes work, and the `money` package that holds amounts exactly is in place. The transaction, import, rules, budget and report routes are still being built.

## How to run it

```
docker run --rm -p 8080:8080 ghcr.io/helenajullia/personal-finance-tracker:latest
```

To build and run without Docker:

```
go run ./cmd/app
go run ./cmd/app --version
```

To run the tests:

```
go vet ./...
go test ./... -race
```

## What you should see

With the service running:

```
curl.exe -i http://localhost:8080/health
```

```
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok","uptime_seconds":3}
```

`http://localhost:8080/` in a browser shows a page with the application name, the student and the project. `POST /reset` answers `204 No Content`.

## Routes

| Route | Method | Purpose |
|---|---|---|
| `/health` | GET | status and uptime |
| `/version` | GET | app name, version, commit, build time |
| `/` | GET | identity page |
| `/reset` | POST | discard all data |
| `/transactions` | POST | add one transaction |
| `/transactions/import` | POST | import a CSV sent as a string |
| `/transactions` | GET | list, filtered by `from`, `to`, `category`, `q` |
| `/rules` | PUT | set the ordered categorisation rules |
| `/budgets` | PUT | set a monthly limit for a category |
| `/reports/monthly` | GET | totals and over-budget categories for `?month=YYYY-MM` |

Errors always have the shape `{"error":"<code>","message":"<explanation>"}`.
