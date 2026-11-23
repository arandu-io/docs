# Criptografia

- [Criptografia autenticada com a chave da aplicação](#authenticated-encryption-with-the-application-key)
- [Leia uma vez no boot e criptografe um valor tipado](#parse-once-at-boot-and-encrypt-a-typed-value)
- [Criptografe strings e detecte payloads](#encrypt-strings-and-detect-payloads)
- [Falhas de descriptografia permanecem indistinguíveis](#decryption-failures-stay-indistinguishable)
- [Teste ida e volta e payloads hostis](#test-round-trips-and-hostile-payloads)
- [Chaves, encrypter e assinador](#keys-encrypter-and-signer)
- [Criptografia não substitui hashing nem ocultação de dados sensíveis](#encryption-does-not-replace-hashing-or-redaction)

<a name="authenticated-encryption-with-the-application-key"></a>
## Criptografia autenticada com a chave da aplicação

Criptografia usa AES-256-GCM, unindo confidencialidade e detecção de adulteração. Encrypter guarda cópia da chave de 32 bytes e permite uso concorrente. Helpers genéricos carregam valores compatíveis com JSON; métodos de string carregam texto bruto no formato portátil de payload.

<a name="parse-once-at-boot-and-encrypt-a-typed-value"></a>
## Leia uma vez no boot e criptografe um valor tipado

`ParseKey` aceita a representação gerada com prefixo base64 ou exatamente 32 bytes brutos. Material inválido ou ausente impede o boot. Construa o encrypter com cifra explícita e use tipos genéricos correspondentes na ida e volta para obter o tipo esperado.

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
## Criptografe strings e detecte payloads

`EncryptString` e `DecryptString` evitam codificação JSON de strings quando a interoperabilidade exige texto bruto. `AppearsEncrypted` reconhece o envelope sem descriptografar. O componente expõe `PreviousKeys`, mas o bootstrap publicado ainda não conecta chaves antigas ao encrypter de ponta a ponta; este guia não promete um procedimento de rotação da chave da aplicação.

<a name="decryption-failures-stay-indistinguishable"></a>
## Falhas de descriptografia permanecem indistinguíveis

Envelope malformado, texto cifrado alterado, chave errada, tag inválida e payload não suportado envolvem `ErrDecrypt`. Não exponha mensagens distintas ao cliente: a diferença vira um oráculo de falsificação. `ErrEncrypt` indica falha de configuração ou serialização JSON; interrompa a operação em vez de repetir com outro formato.

<a name="test-round-trips-and-hostile-payloads"></a>
## Teste ida e volta e payloads hostis

Teste ida e volta de structs representativas, strings Unicode, texto vazio e valores zero. Altere bytes do envelope, use outra chave, remova campos e confirme `ErrDecrypt` em cada ramo externo. Um teste de rotação da aplicação ainda é necessário antes de ensinar composição de chaves antigas em implantação.

<a name="keys-encrypter-and-signer"></a>
## Chaves, encrypter e assinador

`GenerateKey` produz chave segura para configuração; `ParseKey` valida o material; Supported verifica par chave/cifra; `NewEncrypter` constrói o serviço. Encrypt, Decrypt, `EncryptString` e `DecryptString` tratam payloads. `NewSigner` usa a mesma chave da aplicação para assinaturas e links com expiração sem expor internos da criptografia.

<a name="encryption-does-not-replace-hashing-or-redaction"></a>
## Criptografia não substitui hashing nem ocultação de dados sensíveis

Senhas usam hashing unidirecional porque não precisam ser recuperadas. Tokens que só exigem integridade usam Signer. Logs e JSON ainda precisam de remoção explícita de dados sensíveis: valor criptografado pode ser copiado ou repetido mesmo ocultando texto puro. Mantenha `APP_KEY` fora de structs serializadas, query, diagnósticos e versionamento.
