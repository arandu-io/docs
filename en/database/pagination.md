# Pagination

- [A paginator answers navigation, not only rows](#a-paginator-answers-navigation-not-only-rows)
- [Choose length-aware, simple, or cursor pagination](#choose-length-aware-simple-or-cursor-pagination)
- [Preserve filters and stable ordering](#preserve-filters-and-stable-ordering)
- [Reject malformed input and impossible positions](#reject-malformed-input-and-impossible-positions)
- [Test boundaries and concurrent change](#test-boundaries-and-concurrent-change)
- [Three terminal shapes](#three-terminal-shapes)
- [Pagination is part of the response contract](#pagination-is-part-of-the-response-contract)

> [!WARNING]
> **Current limitations.** Cursor signing, expiry, and deterministic ordering are released, but no direct first-party application journey consumes this pagination path. Keep the ordering unique and test cursor expiry at the application boundary.

<a name="a-paginator-answers-navigation-not-only-rows"></a>
## A paginator answers navigation, not only rows

Returning a slice alone cannot distinguish a full final page from a full page with more data behind it. The pagination terminals return typed items together with the state needed to build stable navigation. They execute the query under the same Grant and tenant scope as Get.

<a name="choose-length-aware-simple-or-cursor-pagination"></a>
## Choose length-aware, simple, or cursor pagination

Paginate runs the data query and a count for total pages. `SimplePaginate` avoids the total count and only determines whether another page exists. `CursorPaginate` uses ordered values encoded in an opaque cursor, avoiding large offsets on a changing or deep result set.

```go
invoices, page, err := models.Invoices(db).
	OrderBy("created_at").
	Paginate(ctx, g, 25, requestedPage, pagination.Options{})
if err != nil { return nil, nil, err }
return invoices, page, nil
```

<a name="preserve-filters-and-stable-ordering"></a>
## Preserve filters and stable ordering

Compose search and tenant-independent filters before the terminal, validate per-page bounds at the request layer, and copy approved query parameters into generated links. Cursor pagination needs a deterministic order with a unique tie-breaker so equal timestamps cannot make a row appear twice or disappear between requests.

<a name="reject-malformed-input-and-impossible-positions"></a>
## Reject malformed input and impossible positions

The released model pagination terminals normalize a page number below one to page 1. Request validation may still reject that input when the HTTP contract chooses to, but the database terminal itself does not return a negative-page error. Unbounded per-page values, malformed cursors, and cursors that cannot be decoded remain request-boundary concerns. A page beyond the end should follow the endpoint's documented empty-page or not-found contract.

<a name="test-boundaries-and-concurrent-change"></a>
## Test boundaries and concurrent change

Cover page values below one and assert they normalize to page 1, then cover empty results, exactly one page, a full final page, one extra record, the maximum permitted page size, and another tenant with matching sort values. For cursors, insert a row between requests and assert the chosen ordering gives the documented continuation without trusting the cursor's internal text.

<a name="three-terminal-shapes"></a>
## Three terminal shapes

Paginate returns the rows and a `*pagination.LengthAwarePage` with total information. `SimplePaginate` returns the rows and a `*pagination.Page` focused on previous and next navigation. `CursorPaginate` accepts a parsed `pagination.Cursor` and returns the rows and a `*pagination.CursorPage`. All three take context, Grant, page size, and options after the query has been composed.

<a name="pagination-is-part-of-the-response-contract"></a>
## Pagination is part of the response contract

HTTP resources decide how items, totals, and links are serialized; the page only supplies the facts. Query validation owns allowed filters and sort columns. Use Chunk or Cursor for server-side batch work where no user-facing navigation is needed, because a page object is overhead when the consumer only processes every row once.
