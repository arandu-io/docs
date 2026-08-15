# ADR 0035-a-conexao-e-uma-url — uma variável por serviço

**Status:** Aceita.

## Contexto

O banco vinha de **sete** variáveis — `DB_CONNECTION`, `DB_HOST`, `DB_PORT`,
`DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE` — e mais `DATABASE_URL`, que
**sobrescrevia todas** quando presente.

Isso é duas formas de dizer onde o banco está, que a REGRA 9 recusa. E a falha
que produz é da pior espécie: metade sobrescrita e metade não, aplicação num
host de uma variável e num banco de outra, com **cada valor individualmente
correto**. Nada falha; conecta no lugar errado.

O mail tinha o mesmo desenho, com sete variáveis próprias e a credencial
espalhada: trocar de provedor deixava `MAIL_KEY` para trás numa variável que o
transporte novo não lê.

O Redis já era `REDIS_URL`, uma só.

## Decisão

**Uma URL por serviço, e o esquema é o driver.**

```
DATABASE_URL=postgres://user:password@host:5432/dbname
DATABASE_URL=mysql://user:password@host:3306/dbname
DATABASE_URL=sqlite://database/database.sqlite

MAIL_URL=log://
MAIL_URL=smtp://user:password@smtp.example.com:587
MAIL_URL=smtps://user:password@smtp.example.com
MAIL_URL=resend://re_xxxxxxxx
MAIL_URL=sendgrid://SG.xxxxxxxx

REDIS_URL=redis://127.0.0.1:6379
```

Treze variáveis viraram duas. Uma string também **viaja**: toda plataforma
gerenciada entrega exatamente isso, um compose passa numa linha, e mudar de uma
para outra é copiar uma coisa em vez de seis.

### O que fica de fora da URL, e por quê

O **pool** (`DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`) e
o **remetente** (`MAIL_FROM_ADDRESS`, `MAIL_FROM_NAME`). Não são identidade da
conexão: são quantas conexões abrir e quem a aplicação é. Não mudam quando o
banco ou o provedor muda. Uma forma de dizer *onde*, outra de dizer *quanto* —
perguntas diferentes, não duas formas da mesma.

### Quatro coisas que a mudança obrigou

**As variáveis antigas são recusadas, não ignoradas.** Ignorar é a falha
silenciosa que isso existe para remover: um `.env` cheio de `DB_PASSWORD` e uma
aplicação conectada em outro lugar. A mensagem carrega a URL a escrever.

**O `DSN()` é construído por driver.** Devolvia a URL crua quando havia uma — o
que está certo para o pgx e errado para os outros dois: o `go-sql-driver` toma
`user:pass@tcp(host:port)/db`, sem esquema, e o `modernc/sqlite` toma um
caminho. Qualquer um dos dois falharia com erro de parse que não nomeia nem o
driver nem a variável.

**Os pragmas do SQLite são escritos à mão**, não por `url.Values.Encode`, que
percent-encoda os parênteses. O driver decodifica de volta, então funciona — e
DSN que ninguém lê num log é DSN que ninguém confere, que é a razão de
`foreign_keys`, `WAL` e `busy_timeout` merecerem ser explícitos.

**`MAIL_URL` inválida derruba o boot.** Cair para o transporte de log seria a
falha que ninguém enxerga: a aplicação sobe, as mensagens vão para o log, e a
primeira pessoa a notar é o cliente que não recebeu.

## Consequências

- Mudança **quebrando**. Todo `.env` existente é recusado no boot, com a linha
  a escrever. É deliberado: aceitar os dois formatos seria manter as duas formas
  que esta ADR existe para eliminar.
- `smtps://` e `smtp://` nomeiam a diferença entre TLS implícito na 465 e
  STARTTLS na 587 — uma porta e um handshake, dentro de um nome só.
- Senha com `@ / : ? #` viaja percent-encoded, como em qualquer URL. Um parser
  que não decodifica conecta com a senha errada e reporta "authentication
  failed" sobre uma senha correta; há teste fixando isso.
- **`FILESYSTEM_*` continua destrinchado** e é o próximo da mesma forma:
  `STORAGE_URL=r2://key:secret@account/bucket`. Não foi feito aqui para não
  misturar duas migrações num commit.
