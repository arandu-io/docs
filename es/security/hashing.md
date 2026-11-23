# Hashing

- [Las contraseñas se verifican, nunca se descifran](#passwords-are-verified-never-decrypted)
- [Calcule el hash en la frontera final de escritura](#hash-at-the-final-write-boundary)
- [Verifique y actualice el hash sin exponer texto plano](#verify-and-rehash-without-exposing-plaintext)
- [Separe hashes malformados de errores de dependencia](#keep-malformed-hashes-and-dependency-errors-distinct)
- [Pruebe recorridos correctos, incorrectos y de actualización](#test-correct-wrong-and-upgrade-paths)
- [Contratos de hashing publicados](#released-hashing-contracts)
- [Aplique hashing en flujos completos de identidad](#apply-hashing-inside-complete-identity-flows)

<a name="passwords-are-verified-never-decrypted"></a>
## Las contraseñas se verifican, nunca se descifran

El contrato de hashing transforma la contraseña en una representación adaptativa unidireccional y compara la enviada sin exponer el valor almacenado. El cifrado y la firma tratan secreto reversible e integridad; ninguno es almacenamiento de contraseña.

<a name="hash-at-the-final-write-boundary"></a>
## Calcule el hash en la frontera final de escritura

Calcule `hashing.Make` sobre texto plano aceptado y compare con `hashing.Check`. El paquete guarda Argon2id PHC canónico y verifica formatos legados publicados sin exponer ni recuperar texto plano.

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
## Verifique y actualice el hash sin exponer texto plano

Use Check en login y solo tras éxito use `NeedsRehash` para sustituir un valor legado o antiguo mediante el servicio autorizado de cuenta. Los parámetros Argon2id canónicos se compilan en :hesapeVersion, no se configuran por petición ni entorno.

- Compare mediante el contrato del hasher.
- Actualice el hash solo tras una verificación correcta.
- Limite intentos de contraseña con política explícita de almacenamiento.
- Destruya otras sesiones tras recuperar la contraseña.

<a name="keep-malformed-hashes-and-dependency-errors-distinct"></a>
## Separe hashes malformados de errores de dependencia

Una contraseña incorrecta es un rechazo normal de autenticación. Un hash almacenado malformado, un fallo de recursos de hashing o persistencia indisponible son errores operativos; nunca deben convertirse en login válido ni respuesta que enumere usuarios.

> [!WARNING]
> No compare hashes como strings, no los registre como identificadores inocuos, no reutilice digests de tokens como hashes de contraseña ni reduzca coste durante una petición para evitar latencia. Limite concurrencia de autenticación y mida la configuración de producción.

<a name="test-correct-wrong-and-upgrade-paths"></a>
## Pruebe recorridos correctos, incorrectos y de actualización

Use contraseñas correctas e incorrectas conocidas, datos codificados malformados, costes actuales y obsoletos, casos Unicode, carga concurrente de login y fallo de escritura de rehash. Confirme ausencia de texto plano y hashes codificados en respuesta, JSON, slog y eventos.

- La contraseña correcta pasa; la incorrecta se rechaza.
- `NeedsRehash` cambia solo tras el éxito.
- El almacenamiento malformado nunca autentica.
- Ningún secreto llega a logs ni recursos.

<a name="released-hashing-contracts"></a>
## Contratos de hashing publicados

El hasher publicado expone Make, Check y `NeedsRehash` mediante la frontera de seguridad. La policy de contraseña de aplicación, la persistencia autorizada, el rate limit y la invalidación de sesiones son colaboradores separados.

| Superficie | Contrato publicado |
| --- | --- |
| Make | produce un hash adaptativo codificado |
| Check | verifica el texto plano enviado |
| `NeedsRehash` | detecta configuración codificada obsoleta |

<a name="apply-hashing-inside-complete-identity-flows"></a>
## Aplique hashing en flujos completos de identidad

La autenticación llama a Check, la recuperación guarda un hash nuevo y destruye sesiones, rate limit limita intentos online y el cifrado documenta protección reversible de datos.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Recuperación de contraseña](/docs/{{version}}/security/password-reset)
- [Limitación de peticiones](/docs/{{version}}/services/rate-limiting)
- [Cifrado](/docs/{{version}}/security/encryption)
