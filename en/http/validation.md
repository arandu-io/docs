# Validation

- [Validation narrows input while policy grants authority](#validation-narrows-input-while-policy-grants-authority)
- [Declare the input beside its use case](#declare-the-input-beside-its-use-case)
- [Reject one create request before storage](#reject-one-create-request-before-storage)
- [Put each rule at the boundary it understands](#put-each-rule-at-the-boundary-it-understands)
- [Preserve field errors without leaking causes](#preserve-field-errors-without-leaking-causes)
- [Allowlist fields and expensive work](#allowlist-fields-and-expensive-work)
- [Test the smallest failing value and the boundary](#test-the-smallest-failing-value-and-the-boundary)
- [Continue from valid data to authorized work](#continue-from-valid-data-to-authorized-work)

<a name="validation-narrows-input-while-policy-grants-authority"></a>
## Validation narrows input while policy grants authority

Validation turns untrusted transport values into a typed application input and a field-indexed error set. It checks shape and business-neutral constraints; authorization and database uniqueness remain separate operations with their own failure semantics.

<a name="declare-the-input-beside-its-use-case"></a>
## Declare the input beside its use case

Keep transport names, domain fields, and maximum sizes explicit. Do not add `TenantID` to browser input: the repository derives tenancy from the Grant.

<a name="reject-one-create-request-before-storage"></a>
## Reject one create request before storage

Bind the request into a typed value whose form tags name every accepted field and whose Rules method declares what each field may be. `ctx.Bind` runs the rules over the raw input before it converts anything and returns a failure as a non-empty `validation.Errors` set, which the action adapter turns into a flash and redirect for a page and a 422 problem document for a JSON client. Bind reads the body alone for POST, PUT, PATCH, and DELETE, so a value in the query string can neither satisfy a rule nor reach the stored record. Normalize email, identifiers, and whitespace before persistence uses the accepted value.

```go example
type RegisterRequest struct {
	Email string `form:"email"`
}

func (RegisterRequest) Rules() validation.Rules {
	return validation.Rules{
		"email": "required|email|max:254",
	}
}

func Register(ctx *fhttp.Context) error {
	var req RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	return nil
}
```

<a name="put-each-rule-at-the-boundary-it-understands"></a>
## Put each rule at the boundary it understands

Declare required, type, size, membership, and comparison rules in Rules, as a pipe string or a list per field. What spans fields or asks the data belongs in a Validate method of the request, which the service calls; do not state again in Validate a restriction Rules already declares. Messages come out in the language of the request, and field names stay stable so forms can place each error beside its input.

- Validate before authorization when malformed input should never reach policy.
- Authorize before a protected lookup reveals record existence.
- Return 422 JSON or a flash-and-redirect form result consistently.
- Keep persistence conflicts distinct from syntax errors.

<a name="preserve-field-errors-without-leaking-causes"></a>
## Preserve field errors without leaking causes

Unknown rules, invalid rule arguments, empty returned error sets, and type coercion that hides malformed input are implementation defects. A database outage is never a field validation message.

> [!WARNING]
> Do not use a uniqueness query without the current tenant and authorization boundary. A globally unique-looking validation rule can disclose another tenant's data or reject a value that is valid in the current scope.

<a name="allowlist-fields-and-expensive-work"></a>
## Allowlist fields and expensive work

> [!WARNING]
> Allowlist fields, bound body bytes, collection counts, string lengths, and upload sizes before expensive work. Never accept a tenant selector, storage path, model field, or password hash from submitted input, and keep credentials out of messages, logs, and flashed old input.

<a name="test-the-smallest-failing-value-and-the-boundary"></a>
## Test the smallest failing value and the boundary

Table-test every rule boundary with the actual field name and message key, then send a failing request through normal and enhanced clients. Assert rejected input changes no record and valid input reaches the service in normalized form.

- Required and maximum boundaries are exact.
- Cross-field confirmation uses the accepted field names.
- 422 responses remain rejected responses.
- No protected lookup runs outside tenant scope.

<a name="continue-from-valid-data-to-authorized-work"></a>
## Continue from valid data to authorized work

Request parsing supplies candidate values, controllers choose the use case, policies decide authority, and responses render the resulting field errors without changing their status.

- [Requests](/docs/{{version}}/http/requests)
- [Controllers](/docs/{{version}}/http/controllers)
- [Responses](/docs/{{version}}/http/responses)
