# Hashing

- [Passwords are verified, never decrypted](#passwords-are-verified-never-decrypted)
- [Hash at the final write boundary](#hash-at-the-final-write-boundary)
- [Verify and rehash without exposing plaintext](#verify-and-rehash-without-exposing-plaintext)
- [Keep malformed hashes and dependency errors distinct](#keep-malformed-hashes-and-dependency-errors-distinct)
- [Test correct, wrong, and upgrade paths](#test-correct-wrong-and-upgrade-paths)
- [Released hashing contracts](#released-hashing-contracts)
- [Apply hashing inside complete identity flows](#apply-hashing-inside-complete-identity-flows)

<a name="passwords-are-verified-never-decrypted"></a>
## Passwords are verified, never decrypted

The hashing contract turns a password into an adaptive one-way representation and compares a submitted password without exposing the stored value. Encryption and signing solve reversible secrecy and integrity; neither is a password store.

<a name="hash-at-the-final-write-boundary"></a>
## Hash at the final write boundary

Hash accepted plaintext with `hashing.Make` and compare with `hashing.Check`. The package writes canonical Argon2id PHC values and can verify released legacy encodings without exposing or recovering plaintext.

```go example
func CheckPassword(candidate string) error {
	encoded, err := hashing.Make(candidate)
	if err != nil {
		return err
	}

	if err := hashing.Check(candidate, encoded); err != nil {
		return err
	}
	return nil
}
```

<a name="verify-and-rehash-without-exposing-plaintext"></a>
## Verify and rehash without exposing plaintext

Use Check for login, and only after success use `NeedsRehash` to replace a legacy or older encoded value through the authorized account service. Canonical Argon2id parameters are compiled into :hesapeVersion and are not request or environment configuration.

- Compare through the hasher contract.
- Rehash only after a successful check.
- Rate-limit password attempts with an explicit store policy.
- Destroy other sessions after a password reset.

<a name="keep-malformed-hashes-and-dependency-errors-distinct"></a>
## Keep malformed hashes and dependency errors distinct

Wrong password is an ordinary authentication refusal. A malformed stored hash, hashing resource failure, or persistence outage is an operational error and must not become a false successful login or a user-enumerating response.

> [!WARNING]
> Do not compare hashes as strings, log them as harmless identifiers, reuse token digests as password hashes, or lower cost on a request path to avoid latency. Bound authentication concurrency and benchmark production settings.

<a name="test-correct-wrong-and-upgrade-paths"></a>
## Test correct, wrong, and upgrade paths

Use known correct and wrong passwords, malformed encoded data, current and obsolete cost values, Unicode policy cases, concurrent login load, and a rehash write failure. Assert plaintext and encoded hashes are absent from response, JSON, slog, and events.

- Correct password checks and wrong password refuses.
- `NeedsRehash` changes only after success.
- Malformed storage never authenticates.
- No secret reaches logs or resources.

<a name="released-hashing-contracts"></a>
## Released hashing contracts

The released hasher exposes Make, Check, and `NeedsRehash` through the security boundary. Application password policy, authorized persistence, rate limiting, and session invalidation remain separate collaborators.

| Surface | Released contract |
| --- | --- |
| Make | produce an encoded adaptive hash |
| Check | verify submitted plaintext |
| `NeedsRehash` | detect an obsolete encoded configuration |

<a name="apply-hashing-inside-complete-identity-flows"></a>
## Apply hashing inside complete identity flows

Authentication calls Check, password reset writes a new hash and destroys sessions, rate limiting bounds online guesses, and encryption documents reversible data protection.

- [Authentication](/docs/{{version}}/security/authentication)
- [Password reset](/docs/{{version}}/security/password-reset)
- [Rate limiting](/docs/{{version}}/services/rate-limiting)
- [Encryption](/docs/{{version}}/security/encryption)
