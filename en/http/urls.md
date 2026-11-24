# URLs and signed URLs

- [Route names are the address source](#route-names-are-the-address-source)
- [Build and redirect by name](#build-and-redirect-by-name)
- [Keep paths, hosts, and signatures separate](#keep-paths-hosts-and-signatures-separate)
- [Do not hide the framework bridge gap](#do-not-hide-the-framework-bridge-gap)
- [Change every signed byte](#change-every-signed-byte)
- [Released URL surfaces and boundary](#released-url-surfaces-and-boundary)
- [Use links at their real trust level](#use-links-at-their-real-trust-level)

<a name="route-names-are-the-address-source"></a>
## Route names are the address source

Named routes are the source of application addresses. Reverse a route through its table instead of concatenating segments, then sign only the canonical path and bounded query values when a link must be tamper-evident.

<a name="build-and-redirect-by-name"></a>
## Build and redirect by name

Generate the named path first and fail on an unknown name or parameter mismatch. A signer can then attach an expiry and signature for a one-purpose link; verification must run before the handler uses any signed query value.

```go example
func SignInvoiceLink(router *fhttp.Router, signer *encryption.Signer, invoiceID string) (string, error) {
	path, err := router.Table().URL("invoices.show", invoiceID)
	if err != nil {
		return "", err
	}

	token := signer.Sign("invoice.show", path, 15*time.Minute)
	verifiedPath, err := signer.Verify("invoice.show", token)
	if err != nil {
		return "", err
	}
	return verifiedPath, nil
}
```

<a name="keep-paths-hosts-and-signatures-separate"></a>
## Keep paths, hosts, and signatures separate

Keep host and scheme decisions in trusted deployment configuration. Generate absolute URLs only after trusted-proxy normalization, preserve canonical parameter order in the signing implementation, and use a short expiry for email actions.

- Name routes when they are registered.
- Handle URL reversal errors in application code.
- Sign the canonical destination, not caller-provided text.
- Verify expiry and signature before any mutation.

<a name="do-not-hide-the-framework-bridge-gap"></a>
## Do not hide the framework bridge gap

Unknown route names, wrong parameter counts, changed queries, expired signatures, and untrusted host data are distinct refusal paths. Never fall back to an unsigned destination when signing fails. `hesape/routing.SignedRoute` takes the route table of `hesape/routing`, not the table of `framework/http.Router`; with the framework router, get the path from `Table().URL` and sign it with `security.Signer`, as in the example above.

> [!WARNING]
> A valid signature proves integrity and expiry, not identity or permission. The handler still loads the current Subject, authorizes the action, and treats a single-use operation as single-use in durable storage.

<a name="change-every-signed-byte"></a>
## Change every signed byte

Round-trip a named URL, sign it, verify it, and then alter one path segment, query value, and expiry. Test absolute URLs behind the actual trusted-proxy topology and assert an attacker-controlled Host cannot become the canonical origin.

- Route parameters are escaped and counted.
- A modified value fails verification.
- An expired link fails verification.
- Authorization still runs after signature success.

<a name="released-url-surfaces-and-boundary"></a>
## Released URL surfaces and boundary

`Route.Name` records stable identity and `Routes.URL` reverses it. The released security signer produces and verifies bounded signed destinations. These APIs do not implement account lookup, authorization, or one-time consumption.

<a name="use-links-at-their-real-trust-level"></a>
## Use links at their real trust level

Routing owns the path template, trusted proxies establish external scheme and host, authorization owns the action, and one-time codes own durable single-use behavior.

- [Routing](/docs/{{version}}/http/routing)
- [Authorization](/docs/{{version}}/security/authorization)
- [Email verification](/docs/{{version}}/security/email-verification)
