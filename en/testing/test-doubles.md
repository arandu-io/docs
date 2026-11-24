# Test doubles and time

- [A test double implements the consumer-owned contract](#a-test-double-implements-the-consumer-owned-contract)
- [Inject one local fake and inspect its call](#inject-one-local-fake-and-inspect-its-call)
- [Prefer capture APIs for released components](#prefer-capture-apis-for-released-components)
- [Do not fake away the behavior being proved](#do-not-fake-away-the-behavior-being-proved)
- [Run the same contract against fake and real adapter](#run-the-same-contract-against-fake-and-real-adapter)
- [Released capture and fake surfaces](#released-capture-and-fake-surfaces)
- [Combine narrow doubles with complete feature journeys](#combine-narrow-doubles-with-complete-feature-journeys)

<a name="a-test-double-implements-the-consumer-owned-contract"></a>
## A test double implements the consumer-owned contract

A fake lives beside the test for the smallest interface the consumer needs. It records typed calls, returns configured outcomes, and respects context and Grant semantics. Global swap registries make parallel tests interfere and are avoided where per-instance capture exists.

<a name="inject-one-local-fake-and-inspect-its-call"></a>
## Inject one local fake and inspect its call

Define a local fake for the service collaborator, inject it through the normal constructor, invoke the subject, and assert the exact context-derived Subject or Grant and typed input rather than matching an unstructured map.

```go example
type invoiceSenderFake struct {
	grant security.Grant
	id    string
	err   error
}

func (f *invoiceSenderFake) Send(_ context.Context, g security.Grant, id string) error {
	f.grant, f.id = g, id
	return f.err
}

func TestFakeRecordsTheCall(t *testing.T) {
	fake := &invoiceSenderFake{}
	var grant security.Grant

	if err := fake.Send(context.Background(), grant, "inv_123"); err != nil {
		t.Fatal(err)
	}
	if fake.id != "inv_123" {
		t.Fatalf("id = %q", fake.id)
	}
}
```

<a name="prefer-capture-apis-for-released-components"></a>
## Prefer capture APIs for released components

Use HTTP client fakes with stray requests prevented, notification Capture per test, captured mail transports, memory adapters only where their contract matches, and explicit clocks or ids when released. Keep each recording owned by one test.

- Fail unexpected calls.
- Record typed inputs and order.
- Respect cancellation and error outcomes.
- Avoid package-global mutable fakes.

<a name="do-not-fake-away-the-behavior-being-proved"></a>
## Do not fake away the behavior being proved

A fake that always succeeds can hide retry, rollback, cancellation, and partial failure. A memory adapter can differ in transaction, ordering, lock, TTL, and cross-process semantics. Never use the double to prove the behavior it intentionally omits. Parser or constructor fakes do not prove that a configured value reaches the live component, so add a feature journey that changes the setting and observes the downstream behavior.

<a name="run-the-same-contract-against-fake-and-real-adapter"></a>
## Run the same contract against fake and real adapter

Run a shared behavioral contract against fake and real adapter where substitution is promised, then add a feature test through bootstrap. Run tests in parallel and with race detection to expose global recordings or unsafe callbacks.

- Unexpected calls fail.
- Context and Grant are preserved.
- Configured failures reach the subject.
- Parallel tests share no recording.

<a name="released-capture-and-fake-surfaces"></a>
## Released capture and fake surfaces

The HTTP client Factory supplies Fake, `FakeSequence`, `StubUrl`, `PreventStrayRequests`, and sent assertions; `notifications.Capture` returns per-test channels and deliveries; other application doubles implement consumer interfaces directly.

| Surface | Released contract |
| --- | --- |
| client Factory fakes | outbound request capture |
| `notifications.Capture` | per-test channel deliveries |
| local interface fake | application service or repository seam |
| shared contract | behavioral parity across adapters |

<a name="combine-narrow-doubles-with-complete-feature-journeys"></a>
## Combine narrow doubles with complete feature journeys

HTTP tests use doubles at external edges, database tests keep concrete storage, contracts define substitution, and the overview chooses the smallest complete boundary.

- [HTTP testing](/docs/{{version}}/testing/http)
- [Database testing](/docs/{{version}}/testing/database)
- [Contracts](/docs/{{version}}/architecture/contracts-adapters)
- [Testing overview](/docs/{{version}}/testing/overview)
