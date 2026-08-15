# ADR 0018-teste-contra-banco-real — O subconjunto portável só vale se for testado nos três

**Status:** Aceita. Fecha uma porta que a ADR 0009 (um schema, três engines)
tinha deixado aberta sem que ninguém percebesse.

## Contexto

A ADR 0009 promete que um schema serve SQLite, PostgreSQL e MySQL, e a promessa
é levada a sério no código: os tipos são o subconjunto portável, o `Rebind`
converte `?` em `$N`, os drivers ficam em compartimentos.

Todo teste do projeto rodava em SQLite.

Isso não foi preguiça — foi uma escolha defensável: `modernc.org/sqlite` é Go
puro, não precisa instalar nada, e o CI fica rápido. O custo apareceu numa
auditoria adversarial: **o MySQL nunca funcionou**.

`aru migrate` cria a própria tabela de controle antes de qualquer coisa, e ela
declarava `id TEXT PRIMARY KEY`. O MySQL guarda `TEXT` fora da página e recusa a
coluna numa chave sem prefixo de tamanho:

```
Error 1170 (42000): BLOB/TEXT column 'id' used in key specification
                    without a key length
```

Primeira instrução da primeira migration. Em qualquer projeto, desde sempre. E o
mesmo formato se repetia em cinco lugares: outbox, users, jobs, a tabela de
controle e o template do gerador — mais o `UNIQUE (tenant_id, email)`, que tem o
mesmo problema.

Nada pegou porque nada nunca falou com um MySQL. O SQLite aceita `TEXT PRIMARY
KEY` sem reclamar, então 251 verificações verdes diziam a verdade sobre um
engine e nada sobre os outros dois.

Dois outros achados têm a mesma raiz:

- `decimal` era `REAL`, que no PostgreSQL tem **quatro bytes**. Um `float64`
  volta arredondado em cerca de sete dígitos — na leitura, em silêncio, sem erro
  em lugar nenhum. No SQLite, `REAL` é IEEE de oito bytes e o round-trip é
  exato.
- `date` trunca no PostgreSQL e não trunca no SQLite, então o mesmo código
  devolve valores diferentes por engine, e uma comparação entre duas datas é
  verdadeira num e falsa no outro.

## Decisão

**Toda afirmação sobre o subconjunto portável é verificada contra os três
engines, ou não é uma afirmação que este framework faz.**

Concretamente:

1. `github.com/arandu-io/database/conformance` é uma suíte só, chamada pelos
   três compartimentos. Ela afirma o subconjunto em si: a tabela de controle, um
   schema no formato que o `aru make:module` emite (chave de texto, `UNIQUE`
   composto, índice), round-trip de `float64`, round-trip de timestamp em UTC, e
   rollback.

2. **Sem servidor, ela pula — não passa.** `go test ./...` continua funcionando
   numa máquina sem nada instalado, e o `plans/checklist.sh` mostra "pendente"
   em vez de "ok". Um script que diz "ok" quando nada rodou é o sinal falso que
   este projeto já encontrou vezes demais.

3. `plans/testar-com-bancos-reais.sh` roda tudo de uma vez, com as credenciais
   fora do script (REGRA 6).

4. `data.KeyText` é a grafia de uma coluna de texto que participa de chave.
   `TEXT` fica para conteúdo livre que ninguém indexa. `aru make:module` segue a
   mesma divisão, que a própria DSL já fazia entre `string` ("texto curto, até
   uma linha") e `text` ("texto longo") e que o gerador achatava.

5. `data.Day` normaliza uma data para meia-noite UTC na escrita, para os três
   engines guardarem a mesma coisa.

## O que isso fecha

O CI ganha serviços de PostgreSQL e MySQL. Isso o deixa mais lento e mais caro
que rodar só SQLite, e é o preço de a ADR 0009 ser verdade em vez de ser uma
intenção.

Não abre exceção: não existe "teste que só roda em SQLite porque o resto é
difícil". Se uma afirmação sobre SQL não dá para verificar nos três, ela sai da
documentação.

## O que isso não fecha

A suíte cobre o subconjunto portável, não o comportamento de cada driver. Um bug
do pgx continua sendo um bug do pgx.

E ela não substitui os testes em SQLite: aqueles rodam em toda máquina, em todo
commit, e são o que dá retorno rápido. A conformidade é a segunda rede, não a
primeira.

## Referências

- ADR 0009 — um schema, três engines
- ADR 0014 — database em compartimentos
- `docs/26-revisao-externa.md` — a revisão que motivou o hábito de auditar
- `plans/testar-com-bancos-reais.sh`
