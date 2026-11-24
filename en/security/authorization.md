# Policies and Grant

- [A policy turns identity into bounded authority](#a-policy-turns-identity-into-bounded-authority)
- [Carry a trusted Subject to the use case](#carry-a-trusted-subject-to-the-use-case)
- [Authorize the operation and pass its Grant](#authorize-the-operation-and-pass-its-grant)
- [Keep the decision beside the protected action](#keep-the-decision-beside-the-protected-action)
- [Refuse without revealing another tenant](#refuse-without-revealing-another-tenant)
- [Never reconstruct authority from request input](#never-reconstruct-authority-from-request-input)
- [Prove denial and tenant isolation at storage](#prove-denial-and-tenant-isolation-at-storage)
- [Identity starts the boundary and Grant finishes it](#identity-starts-the-boundary-and-grant-finishes-it)

<a name="a-policy-turns-identity-into-bounded-authority"></a>
## A policy turns identity into bounded authority

Authentication produces `security.Subject`: a trusted tenant and actor identity. Authorization asks a typed Policy whether that subject may perform one operation on one candidate, then returns `security.Grant` for the protected repository call. A role string alone is not authority.

<a name="carry-a-trusted-subject-to-the-use-case"></a>
## Carry a trusted Subject to the use case

Load Subject through session or token middleware and build policies with the collaborators required for the decision. Do not load a tenant-owned record through an unscoped repository merely to decide whether the caller may see it.

- Use the Subject supplied by the authentication boundary.
- Define the operation and candidate type explicitly.
- Keep tenant-aware lookup inside the policy or authorized repository.
- Make guest and denied outcomes deliberate.

<a name="authorize-the-operation-and-pass-its-grant"></a>
## Authorize the operation and pass its Grant

Call `security.Authorize` with the policy, subject, action, and candidate. Pass only the returned Grant to a repository whose method requires it, making the authorization input impossible to omit at the call site. The changes are a literal map of the columns this action may write, read from a typed request, so no submitted key reaches the update. The adapter must still derive tenant scope with `data.Tenant(grant)`.

```go example
const InvoiceUpdate security.Action = "invoice.update"

type Invoice struct{ ID string }

type UpdateInvoiceRequest struct{ Reference, DueOn string }

type Invoices interface {
	Update(ctx context.Context, grant security.Grant, id string, fields map[string]any) error
}

func Update(ctx context.Context, policy security.Policy[Invoice], subject security.Subject, invoices Invoices, invoice Invoice, req UpdateInvoiceRequest) error {
	grant, err := security.Authorize(ctx, policy, subject, InvoiceUpdate, invoice)
	if err != nil {
		return err
	}

	return invoices.Update(ctx, grant, invoice.ID, map[string]any{
		"reference": req.Reference,
		"due_on":    req.DueOn,
	})
}
```

<a name="keep-the-decision-beside-the-protected-action"></a>
## Keep the decision beside the protected action

Use a small policy per domain resource, map operations to stable named actions, and keep controller transport separate. Reauthorize background and queued work from durable provenance rather than serializing a Grant from the original request. A Grant issued while an operator views somebody else's account is read-only: `Grant.ReadOnly` answers true, and every repository, cache, disk, and queue write refuses it with `auth.ErrReadOnly`, so leave the action out of such a page rather than letting it fail.

- Authorize before protected reads and writes.
- Pass Grant through service and repository seams.
- Return the same public absence for hidden foreign records.
- Reconstruct trusted Subject for asynchronous work.

<a name="refuse-without-revealing-another-tenant"></a>
## Refuse without revealing another tenant

A missing Subject is unauthenticated, a policy refusal is forbidden or intentionally hidden as not found, and a repository outage is operational. Do not collapse them into an empty record or continue with a zero Grant.

> [!WARNING]
> A zero or foreign Grant must match no tenant-owned row. Treat an adapter that ignores Grant as a contract failure even if controller tests remain green.

<a name="never-reconstruct-authority-from-request-input"></a>
## Never reconstruct authority from request input

> [!WARNING]
> Never accept tenant, role, Subject, or Grant from a form, path, query, token claim without verification, job payload, or MCP arguments. Never cache an authorization result beyond the subject, action, resource, and policy state it actually covered.

Authentication and signed URLs prove different facts. A valid session or signature still needs the policy decision for the requested application action.

<a name="prove-denial-and-tenant-isolation-at-storage"></a>
## Prove denial and tenant isolation at storage

Run a shared repository contract with allowed, denied, zero, and foreign Grants. Exercise each policy action, guest behavior, cross-tenant identifiers, cancellation, and a background job that rebuilds identity through the trusted application path.

- Allowed Subject receives a bounded Grant.
- Denied Subject reaches no protected adapter operation.
- Foreign tenant ids reveal no existence.
- Every concrete repository enforces Grant.

<a name="identity-starts-the-boundary-and-grant-finishes-it"></a>
## Identity starts the boundary and Grant finishes it

Authentication supplies Subject, sessions transport browser identity, contracts keep Grant in adapter signatures, and controllers remain transport-only.

- [Authentication](/docs/{{version}}/security/authentication)
- [Session](/docs/{{version}}/http/session)
- [Contracts and adapters](/docs/{{version}}/architecture/contracts-adapters)
- [Controllers](/docs/{{version}}/http/controllers)
