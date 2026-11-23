# Hashing

- [Senhas são verificadas, nunca descriptografadas](#passwords-are-verified-never-decrypted)
- [Calcule o hash na fronteira final da escrita](#hash-at-the-final-write-boundary)
- [Verifique e atualize o hash sem expor texto puro](#verify-and-rehash-without-exposing-plaintext)
- [Separe hashes malformados de erros de dependência](#keep-malformed-hashes-and-dependency-errors-distinct)
- [Teste caminhos corretos, incorretos e de atualização](#test-correct-wrong-and-upgrade-paths)
- [Contratos de hashing publicados](#released-hashing-contracts)
- [Aplique hashing em fluxos completos de identidade](#apply-hashing-inside-complete-identity-flows)

<a name="passwords-are-verified-never-decrypted"></a>
## Senhas são verificadas, nunca descriptografadas

O contrato de hashing transforma senha em representação adaptativa unidirecional e compara a senha enviada sem expor o valor armazenado. Criptografia e assinatura tratam segredo reversível e integridade; nenhuma é armazenamento de senha.

<a name="hash-at-the-final-write-boundary"></a>
## Calcule o hash na fronteira final da escrita

Calcule `hashing.Make` sobre texto puro aceito e compare com `hashing.Check`. O pacote grava Argon2id PHC canônico e verifica formatos legados publicados sem expor nem recuperar texto puro.

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
## Verifique e atualize o hash sem expor texto puro

Use Check no login e só após sucesso use `NeedsRehash` para substituir valor legado ou antigo pelo serviço autorizado de conta. Parâmetros Argon2id canônicos são compilados em :hesapeVersion, não configurados por requisição ou ambiente.

- Compare pelo contrato do hasher.
- Atualize o hash só após verificação bem-sucedida.
- Limite tentativas de senha com política explícita de armazenamento.
- Destrua outras sessões após recuperar a senha.

<a name="keep-malformed-hashes-and-dependency-errors-distinct"></a>
## Separe hashes malformados de erros de dependência

Senha errada é recusa normal de autenticação. Hash armazenado malformado, falha de recursos de hashing ou persistência indisponível são erros operacionais; nunca devem virar login válido nem resposta que enumere usuários.

> [!WARNING]
> Não compare hashes como strings, não os registre como identificadores inofensivos, não reutilize digests de tokens como hashes de senha nem reduza custo durante requisição para evitar latência. Limite concorrência de autenticação e meça configurações de produção.

<a name="test-correct-wrong-and-upgrade-paths"></a>
## Teste caminhos corretos, incorretos e de atualização

Use senhas corretas e erradas conhecidas, dados codificados malformados, custos atuais e obsoletos, casos Unicode, carga concorrente de login e falha de escrita do rehash. Confirme ausência de texto puro e hashes codificados em resposta, JSON, slog e eventos.

- Senha correta passa; senha errada é recusada.
- `NeedsRehash` muda apenas após sucesso.
- Armazenamento malformado nunca autentica.
- Nenhum segredo chega a logs ou recursos.

<a name="released-hashing-contracts"></a>
## Contratos de hashing publicados

O hasher publicado expõe Make, Check e `NeedsRehash` pela fronteira de segurança. Policy de senha da aplicação, persistência autorizada, rate limit e invalidação de sessões são colaboradores separados.

| Superfície | Contrato publicado |
| --- | --- |
| Make | produz um hash adaptativo codificado |
| Check | verifica o texto puro enviado |
| `NeedsRehash` | detecta configuração codificada obsoleta |

<a name="apply-hashing-inside-complete-identity-flows"></a>
## Aplique hashing em fluxos completos de identidade

Autenticação chama Check, recuperação grava novo hash e destrói sessões, rate limit limita tentativas online e criptografia documenta proteção reversível de dados.

- [Autenticação](/docs/{{version}}/security/authentication)
- [Recuperação de senha](/docs/{{version}}/security/password-reset)
- [Limitação de requisições](/docs/{{version}}/services/rate-limiting)
- [Criptografia](/docs/{{version}}/security/encryption)
