# ADR 0031-transporte-de-provedor-nao-e-submodulo — Resend e SendGrid no core

**Status:** Aceita.

## Contexto

A REGRA 11 manda adapter pesado para submódulo com `go.mod` próprio, porque em
Go não existe dependência opcional: o driver do provedor entra no `go.sum` de
todo mundo que importa o pacote, use ou não.

O `framework/mail` levou essa regra ao pé da letra e escreveu, no comentário dos
transportes:

> Resend and SendGrid are adapters, and they go in submodules of their own —
> `mail/resend` and `mail/sendgrid`

E o `bootstrap/app.go` do `examples` terminava assim:

```go
case appconfig.MailerResend, appconfig.MailerSendGrid:
    panic("mail: " + string(cfg.Mailer) + " needs its submodule: go get ...")
```

Um `panic` no boot, apontando para dois repositórios que **não existem**. O
`config/mail.go` já listava `resend` e `sendgrid` como valores aceitos de
`MAIL_MAILER` desde o começo. Quem colocasse o valor documentado no `.env` não
recebia um erro de configuração: recebia um processo que não sobe.

## Decisão

**Resend e SendGrid ficam no core, em `framework/mail/api.go`, escritos sobre
`net/http`.**

A REGRA 11 continua valendo e não foi contornada — ela foi lida de novo. O que
ela manda para submódulo é o adapter que **arrasta um SDK**. Estes dois não
arrastam nada:

| | |
|---|---|
| Resend | `POST https://api.resend.com/emails`, `Authorization: Bearer`, corpo JSON |
| SendGrid | `POST https://api.sendgrid.com/v3/mail/send`, mesma forma |

O custo dos dois somados, no core, é `net/http` e `encoding/json` — que todo
binário Arandu já linka. O `resend-go` traria árvore transitiva própria em troca
de uma requisição que cabe numa tela.

**Amazon SES continua fora**, e agora por um motivo que se enxerga melhor: é o
único dos três que precisaria de implementação de assinatura de requisição. É
exatamente o caso que a REGRA 11 descreve, e é o que ninguém pediu.

### O que a decisão trouxe junto

**`mail.ErrRetryable`.** 429 e 5xx não são o mesmo evento que endereço
recusado. Tratados igual, um rate limit engole o link de verificação em
silêncio; distinguidos, o job reagenda e o endereço errado falha de vez.

**`Content.TextView`.** O `Content` só tinha `Text string`, um literal Go. O
projeto gerava `mail/password-reset-text.kyse.go`, a view existia na árvore, e
**não havia como enviá-la** — toda mensagem saía só em HTML. Achado por auditoria
ao escrever o e-mail de verificação.

## Consequências

- Some o `panic` no boot. `MAIL_MAILER=resend` com `MAIL_KEY` preenchido envia.
- Dois repositórios a menos para criar. A REGRA 3 continua em treze.
- O core ganha ~250 linhas e **zero** dependências.
- A tag do envelope é reduzida ao que o provedor aceita antes de sair. A Resend
  recusa a **mensagem inteira** por causa de um espaço num rótulo, e o e-mail
  perdido seria o de confirmação de conta.
- Se um provedor futuro exigir SDK, ele vai para submódulo, como sempre. A
  regra não mudou; a leitura ficou correta.
