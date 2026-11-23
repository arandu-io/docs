# Encryption

- [Authenticated encryption with the application key](#authenticated-encryption-with-the-application-key)
- [Parse once at boot and encrypt a typed value](#parse-once-at-boot-and-encrypt-a-typed-value)
- [Encrypt strings and detect payloads](#encrypt-strings-and-detect-payloads)
- [Decryption failures stay indistinguishable](#decryption-failures-stay-indistinguishable)
- [Test round trips and hostile payloads](#test-round-trips-and-hostile-payloads)
- [Keys, encrypter, and signer](#keys-encrypter-and-signer)
- [Encryption does not replace hashing or redaction](#encryption-does-not-replace-hashing-or-redaction)

<a name="authenticated-encryption-with-the-application-key"></a>
## Authenticated encryption with the application key

The encryption package uses AES-256-GCM so confidentiality and tamper detection are one operation. Encrypter holds a copied 32-byte key and is safe for concurrent use. Generic helpers carry JSON-compatible values, while string methods carry raw text in the portable payload format.

<a name="parse-once-at-boot-and-encrypt-a-typed-value"></a>
## Parse once at boot and encrypt a typed value

`ParseKey` accepts either the generated base64-prefixed representation or exactly 32 raw bytes. Invalid or missing material stops boot. Construct the encrypter with the explicit cipher, then use matching generic types for round trips so decoding returns the type the caller expects.

```go example
type Secret struct{ Value string }

func RoundTripSecret(appKey, value string) (Secret, error) {
	key, err := encryption.ParseKey(appKey)
	if err != nil {
		return Secret{}, err
	}
	encrypter, err := encryption.NewEncrypter(key, encryption.AES256GCM)
	if err != nil {
		return Secret{}, err
	}

	payload, err := encryption.Encrypt(encrypter, Secret{Value: value})
	if err != nil {
		return Secret{}, err
	}
	return encryption.Decrypt[Secret](encrypter, payload)
}
```

<a name="encrypt-strings-and-detect-payloads"></a>
## Encrypt strings and detect payloads

`EncryptString` and `DecryptString` avoid JSON string encoding when interoperability needs the raw text. `AppearsEncrypted` recognizes the package envelope without attempting decryption. The component exposes `PreviousKeys`, but the released application bootstrap does not yet wire retired keys into its encrypter end to end; this guide therefore does not promise an application key-rotation procedure.

<a name="decryption-failures-stay-indistinguishable"></a>
## Decryption failures stay indistinguishable

A malformed envelope, altered ciphertext, wrong key, invalid tag, and unsupported payload all wrap `ErrDecrypt`. Do not translate those into different client messages, because difference becomes a forgery oracle. `ErrEncrypt` indicates configuration or JSON serialization failure and should stop the operation rather than trigger a retry with another format.

<a name="test-round-trips-and-hostile-payloads"></a>
## Test round trips and hostile payloads

Round-trip representative structs, Unicode strings, empty text, and zero values. Flip bytes in the envelope, use a different key, remove fields, and assert each external branch matches `ErrDecrypt`. An application-level rotation test remains required before this page can teach retired-key deployment wiring.

<a name="keys-encrypter-and-signer"></a>
## Keys, encrypter, and signer

`GenerateKey` emits a configuration-safe key; `ParseKey` validates configured material; Supported checks a key and cipher pair; `NewEncrypter` constructs the service. Encrypt, Decrypt, `EncryptString`, and `DecryptString` handle payloads. `NewSigner` uses the same application key for signatures and expiring links without exposing encryption internals.

<a name="encryption-does-not-replace-hashing-or-redaction"></a>
## Encryption does not replace hashing or redaction

Passwords use one-way password hashing because they never need recovery. Tokens that only need integrity use Signer. Logs and JSON responses still require explicit redaction: an encrypted value can be copied or replayed even when its plaintext is hidden. Keep `APP_KEY` out of serialized structs, query strings, diagnostics, and source control.
