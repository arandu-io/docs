# Cifrado

- [Cifrado autenticado con la clave de aplicación](#authenticated-encryption-with-the-application-key)
- [Lea una vez al arrancar y cifre un valor tipado](#parse-once-at-boot-and-encrypt-a-typed-value)
- [Cifre strings y detecte payloads](#encrypt-strings-and-detect-payloads)
- [Los fallos de descifrado permanecen indistinguibles](#decryption-failures-stay-indistinguishable)
- [Pruebe ida y vuelta y payloads hostiles](#test-round-trips-and-hostile-payloads)
- [Claves, encrypter y firmante](#keys-encrypter-and-signer)
- [El cifrado no sustituye hashing ni ocultación de datos sensibles](#encryption-does-not-replace-hashing-or-redaction)

<a name="authenticated-encryption-with-the-application-key"></a>
## Cifrado autenticado con la clave de aplicación

El cifrado usa AES-256-GCM, uniendo confidencialidad y detección de alteración. Encrypter guarda una copia de la clave de 32 bytes y permite uso concurrente. Los helpers genéricos llevan valores compatibles con JSON; los métodos string llevan texto sin procesar en el formato portátil de payload.

<a name="parse-once-at-boot-and-encrypt-a-typed-value"></a>
## Lea una vez al arrancar y cifre un valor tipado

`ParseKey` acepta la representación generada con prefijo base64 o exactamente 32 bytes sin procesar. Material inválido o ausente impide el arranque. Construya el encrypter con cifrado explícito y use tipos genéricos correspondientes en ida y vuelta para obtener el tipo esperado.

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
## Cifre strings y detecte payloads

`EncryptString` y `DecryptString` evitan codificación JSON de strings cuando la interoperabilidad exige texto sin procesar. `AppearsEncrypted` reconoce el envoltorio sin descifrar. El componente expone `PreviousKeys`, pero el bootstrap publicado aún no conecta claves antiguas al encrypter de extremo a extremo; esta guía no promete un procedimiento de rotación de clave de aplicación.

<a name="decryption-failures-stay-indistinguishable"></a>
## Los fallos de descifrado permanecen indistinguibles

Un envoltorio malformado, texto cifrado alterado, clave incorrecta, tag inválida y payload no compatible envuelven `ErrDecrypt`. No exponga mensajes distintos al cliente: la diferencia se vuelve un oráculo de falsificación. `ErrEncrypt` indica fallo de configuración o serialización JSON; detenga la operación en vez de reintentar con otro formato.

<a name="test-round-trips-and-hostile-payloads"></a>
## Pruebe ida y vuelta y payloads hostiles

Pruebe ida y vuelta de structs representativas, strings Unicode, texto vacío y valores cero. Altere bytes del envoltorio, use otra clave, elimine campos y confirme `ErrDecrypt` en cada rama externa. Una prueba de rotación de aplicación aún es necesaria antes de enseñar composición de claves antiguas en despliegue.

<a name="keys-encrypter-and-signer"></a>
## Claves, encrypter y firmante

`GenerateKey` produce una clave segura para configuración; `ParseKey` valida el material; Supported verifica el par clave/cifrado; `NewEncrypter` construye el servicio. Encrypt, Decrypt, `EncryptString` y `DecryptString` procesan payloads. `NewSigner` usa la misma clave de aplicación para firmas y enlaces con expiración sin exponer detalles internos del cifrado.

<a name="encryption-does-not-replace-hashing-or-redaction"></a>
## El cifrado no sustituye hashing ni ocultación de datos sensibles

Las contraseñas usan hashing unidireccional porque no necesitan recuperarse. Los tokens que solo exigen integridad usan Signer. Los logs y JSON aún necesitan eliminación explícita de datos sensibles: un valor cifrado puede copiarse o repetirse aunque oculte el texto plano. Mantenga `APP_KEY` fuera de structs serializadas, query, diagnósticos y control de versiones.
