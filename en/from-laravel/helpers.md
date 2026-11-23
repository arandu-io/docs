# Why helpers are deliberately small

- [Helpers are package functions](#helpers-are-package-functions)
- [Translate each helper family](#translate-each-helper-family)
- [The package catalogue](#the-package-catalogue)
- [What deliberately has no helper](#what-deliberately-has-no-helper)
- [Where each family is documented](#where-each-family-is-documented)

> [!WARNING]
> **Current limitations.** The str and pagination functions are consumed by tagged releases: the framework derives field labels with `str.Headline` and the starter's note service returns a `pagination.Page`. The support, number and collections functions are released and tested in Hesape :hesapeVersion without a tagged first-party application consumer, so verify a chosen function in the consuming application. No released Arandu module validates Brazilian documents.

<a name="helpers-are-package-functions"></a>
## Helpers are package functions

Arandu has no global helper namespace. Nothing like collect, `data_get`, now or url is visible in every file, nothing is registered at boot, and no helper reads the request or a container behind the caller's back. A helper is an exported function of a named package, called with that package's name: `str.Slug`, `number.CurrencyFromCents`, `support.Blank`.

Look in the standard library first. slices, maps, strings, strconv, `path/filepath`, `net/url`, time and sync answer most helper questions with no Arandu import, and a Go reader recognises them on sight. Reach for a Hesape package only when the function carries a convention several packages must agree on: a slug, a field label, a currency amount, a page of results.

The official helpers are the exported functions of five Hesape packages: support, str, number, collections with its arr subpackage, and pagination. This page is the index from the helper a reader already knows to the call that replaces it. It is not a library of its own, and it adds no function.

<a name="translate-each-helper-family"></a>
## Translate each helper family

Read each row left to right. The middle column is the answer whenever it fits; the right column is the package function to call when the standard library has no equivalent or would repeat a shared convention by hand. A dash means nothing in Hesape replaces the helper, on purpose.

### Arrays and objects

| Helper | Go first | Hesape package function |
| --- | --- | --- |
| `Arr::first`, `Arr::last`, head, last | `s[0]`, `s[len(s)-1]`, `slices.IndexFunc` | `arr.First`, `arr.Last`, `collections.Head` |
| `Arr::where`, `Arr::reject`, `Arr::partition` | a for loop, `slices.DeleteFunc` on a copy | `arr.Where`, `arr.Reject`, `arr.Partition` |
| `Arr::sort`, `Arr::sortDesc` | `slices.SortFunc`, `slices.SortStableFunc` | `arr.Sort`, `arr.SortDesc` |
| `Arr::only`, `Arr::except` | a typed struct; `maps.Clone` followed by delete | `arr.Only`, `arr.Except` |
| `Arr::get`, `data_get`, `Arr::dot`, `Arr::undot` | a field of a typed struct; `encoding/json` into that struct | `arr.Get`, `arr.DataGet`, `arr.Dot`, `arr.Undot`, for a `map[string]`any the application does not own |
| `Arr::wrap`, `Arr::join` | `[]T{v}`, `strings.Join` | `arr.Wrap`, `arr.Join` |
| collect, pluck, `keyBy`, `groupBy` | a slice and a for loop; `slices.Sorted(maps.Keys(m))` | `collections.Collect`, `collections.Map`, `collections.Pluck`, `collections.KeyBy`, `collections.GroupBy`, `collections.Reduce` |

### Strings

| Helper | Go first | Hesape package function |
| --- | --- | --- |
| `Str::contains`, `Str::startsWith`, `Str::endsWith` | `strings.Contains`, `strings.HasPrefix`, `strings.HasSuffix` | - |
| `Str::lower`, `Str::upper`, trim, `Str::replace` | `strings.ToLower`, `strings.ToUpper`, `strings.TrimSpace`, `strings.ReplaceAll` | - |
| `Str::before`, `Str::after` | `strings.Cut` | `str.Before`, `str.After` |
| `Str::slug` | - | `str.Slug`, with the separator passed explicitly |
| `Str::headline`, `Str::camel`, `Str::snake`, `Str::kebab`, `Str::studly` | - | `str.Headline`, `str.Camel`, `str.Snake`, `str.Kebab`, `str.Studly` |
| `Str::limit`, `Str::mask` | byte slicing, which can split a rune | `str.Limit`, `str.Mask` |
| `Str::plural` | - | `str.Plural`, English inflection only |
| `Str::uuid`, `Str::ulid`, `Str::random` | `crypto/rand.Text` for an unguessable token | `str.UUID`, `str.ULID`, `str.Random` |

### Numbers

| Helper | Go first | Hesape package function |
| --- | --- | --- |
| `Number::format` | `strconv.FormatFloat` | `number.Format` |
| `Number::currency` | - | `number.CurrencyFromCents` for stored minor units, `number.Currency` for a fractional amount |
| `Number::percentage` | `fmt.Sprintf` | `number.Percentage` |
| `Number::fileSize` | - | `number.FileSize` |
| `Number::abbreviate`, `Number::forHumans` | - | `number.Abbreviate`, `number.ForHumans` |
| `Number::ordinal`, `Number::spell` | - | `number.Ordinal`, `number.Spell` |
| `Number::clamp` | `min(max(v, low), high)` | `number.Clamp` |

### Paths and URLs

| Helper | Go first | Hesape package function |
| --- | --- | --- |
| `base_path`, `app_path`, `storage_path`, `public_path`, `resource_path`, `database_path`, `lang_path`, `config_path` | `filepath.Join(root, ...)` with the root read from configuration and passed to the constructor that needs it; `embed.FS` for files compiled into the binary | - |
| url, `secure_url` | `url.Parse`, `url.JoinPath`, `url.Values` | - |
| route, `to_route`, action | - | `ctx.URL` and `ctx.RedirectRoute`, methods of the request context rather than package functions |
| paginate, `simplePaginate`, `cursorPaginate` | - | the Model terminals Paginate and `SimplePaginate`; `pagination.NewPage`, `pagination.NewLengthAwarePage`, `pagination.NewCursorPage` and `pagination.OptionsFrom` for a result that does not come from a Model |

### Miscellaneous and other utilities

| Helper | Go first | Hesape package function |
| --- | --- | --- |
| blank, filled | s == "", `len(xs)` == 0, p == nil when the type is known | `support.Blank`, `support.Filled` when it is not |
| tap, transform, with | ordinary statements | `support.Tap`, `support.Transform` |
| now, today, Carbon | `time.Now`, or a clock the constructor receives when a test must fix time | - |
| once | `sync.OnceValue` | - |
| retry, Sleep | a loop that honours `ctx.Done()` and a `time.Timer` | - |
| env, config | the typed configuration struct read once at boot | - |
| app, resolve | a constructor argument a provider in app/Providers wires | - |
| request, auth, session, old | the `*http.Context` the handler receives: `ctx.Query`, `ctx.Input`, `ctx.User`, `ctx.Old` | - |
| e | Kyse and `html/template` escape by default; `html.EscapeString` | - |
| dd, dump, Benchmark | `fmt.Printf("%#v\n", v)`, the debug console, `testing.B` | - |
| Lottery, Pipeline, Uri | `math/rand/v2`, plain function composition, `net/url` | - |

<a name="the-package-catalogue"></a>
## The package catalogue

Call each function by its package name. The packages are imported where they are used and nowhere else, so a file's import block says which conventions it depends on.

| Package | Functions named on this page |
| --- | --- |
| `github.com/arandu-io/hesape/support` | Blank, Filled, Tap, Transform |
| `github.com/arandu-io/hesape/str` | Slug, Headline, Camel, Snake, Kebab, Studly, Before, After, Limit, Mask, Plural, UUID, ULID, Random |
| `github.com/arandu-io/hesape/number` | Format, Currency, `CurrencyFromCents`, Percentage, `FileSize`, Abbreviate, `ForHumans`, Ordinal, Spell, Clamp |
| `github.com/arandu-io/hesape/collections` | Collect, Map, Pluck, `KeyBy`, `GroupBy`, Reduce, Head, and the Filter method of Collection |
| `github.com/arandu-io/hesape/collections/arr` | First, Last, Where, Reject, Partition, Sort, `SortDesc`, Only, Except, Get, `DataGet`, Dot, Undot, Wrap, Join |
| `github.com/arandu-io/hesape/pagination` | `NewPage`, `NewLengthAwarePage`, `NewCursorPage`, `OptionsFrom` |

```go example
func HelperExamples(byID map[string]int, root, baseURL, invoiceID string, read, current int, requestURL *url.URL) error {
	// The standard library first.
	ids := slices.Sorted(maps.Keys(byID))
	exports := filepath.Join(root, "storage", "exports")
	link, err := url.JoinPath(baseURL, "invoices", invoiceID)
	if err != nil {
		return err
	}

	// A Hesape package function where the convention is shared.
	slug := str.Slug("Josefin Sans Café", "-")        // "josefin-sans-cafe"
	label := str.Headline("email_notification_sent")  // "Email Notification Sent"
	amount := number.CurrencyFromCents(123450, "USD") // "$1,234.50"
	empty := support.Blank("   ")                     // true
	page := pagination.NewPage(read, 20, current, pagination.OptionsFrom(requestURL))

	_, _, _, _, _, _, _, _ = ids, exports, link, slug, label, amount, empty, page
	return nil
}
```

> [!NOTE]
> A function that changes the element type, such as `collections.Map`, is a package function rather than a method because a Go method cannot declare its own type parameter. That is a property of the language, not a second style to choose between.

<a name="what-deliberately-has-no-helper"></a>
## What deliberately has no helper

Three things look alike and are not interchangeable. `ctx.URL` belongs to the request, because it needs the router and the current host. `number.CurrencyFromCents` is pure and can be called anywhere. An HTTP client, a mailer or a clock is a collaborator the constructor receives. Turning one into another for the sake of a shorter call hides where the value comes from.

- Configuration and environment values are fields of the typed configuration struct, read once at boot. No helper reads the environment later.
- Services are constructor arguments a provider in app/Providers wires over the drivers. No helper resolves a dependency by name.
- The current user, the session and the submitted input come from the `*http.Context` the handler receives, and from no package-level function.
- Path roots come from configuration. No package knows where the application is installed.

> [!WARNING]
> `number.UseLocale`, `number.WithLocale`, `number.UseCurrency`, `str.UseLanguage`, `str.CreateUUIDsUsing` and `str.FreezeUUIDs` change a default shared by the whole process, and `WithLocale` does so for every goroutine while its callback runs. Call them once at boot or inside a test, never while serving a request; for one request, pass the currency code to `number.Currency` or `number.CurrencyFromCents` instead of changing the default.

National documents are not core helpers. Hesape names no CPF, CNPJ or CEP: the mask package applies and removes any shape the application declares, such as 000.000.000-00, and deliberately does not check the digits, because that arithmetic belongs to whoever owns the document. number writes BRL amounts with pt-BR separators when the process locale is pt-BR, set once at boot. No released Arandu module validates Brazilian documents yet, so an application that needs the check keeps one implementation of it rather than a copy per feature.

<a name="where-each-family-is-documented"></a>
## Where each family is documented

The strings guide explains when str is worth an import and where its ASCII folding stops. The collections guide covers Collection, its lazy form and the arr helpers. The pagination guide shows the Model terminals that build pages for you. support and number have no guide of their own; their package documentation is the reference.

- [Strings](/docs/{{version}}/architecture/strings)
- [Collections](/docs/{{version}}/architecture/collections)
- [Pagination](/docs/{{version}}/database/pagination)
- [support package](https://pkg.go.dev/github.com/arandu-io/hesape/support)
- [number package](https://pkg.go.dev/github.com/arandu-io/hesape/number)
