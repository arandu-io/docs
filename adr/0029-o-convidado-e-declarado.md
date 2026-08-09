# ADR 0029-o-convidado-e-declarado — leitura pública passa pela policy

**Status:** Aceita.

## Contexto

`security.Authorize` recusava sujeito anônimo **antes** de consultar a policy:

```go
if s.ID == "" {
    return Grant{}, fmt.Errorf("%w: anonymous subject on %s", ErrForbidden, a)
}
```

Isso está certo pelo motivo certo — sujeito vazio é quase sempre sessão que não
carregou, e responder pergunta de autorização sobre ninguém é como um furo se
abre por acidente. Mas deixava o produto sem **nenhuma** forma de dizer "qualquer
pessoa pode ler um post publicado".

O único caminho era `security.SystemGrant`, que **pula a policy inteira**. Uma
página pública servida com o instrumento que uma tarefa agendada usa. O blog do
`examples` virou ferramenta de autoria com login, o sitemap listava rascunho, e
o SEO — canonical, Open Graph, sitemap — não servia para nada, porque o
rastreador caía no login.

## Decisão

**`security.Guest(tenant)` devolve um sujeito anônimo declarado, e a policy
decide sobre ele.**

A recusa continua. O que muda é que ela olha um marcador **não exportado** que só
`Guest` escreve:

```go
if s.ID == "" && !s.guest {
    return Grant{}, ...
}
```

- `Subject{}` continua recusado. Sessão esquecida não vira requisição pública.
- `Guest(tenant)` chega na policy, e a policy responde.
- Policy que não fala de convidado **recusa** — toda policy gerada nega tudo,
  então o padrão é fechado.
- `IsGuest()` deixa a regra explícita no arquivo onde ela mora.

O tenant é obrigatório e é da aplicação, da configuração. A REGRA 14 não é
suspensa porque ninguém entrou.

### Nada foi afrouxado

A autorização continua num lugar só, o `Grant` continua o único caminho até o
repositório, e a decisão continua sendo da Policy. O que existia antes era
**pior**: `SystemGrant` não consulta policy nenhuma.

### Duas coisas que o `examples` descobriu ao usar

**Ação nomeada para a listagem pública.** `PostList` pagina tudo, rascunho
incluído. Abrir ela para convidado com um filtro seria uma permissão cujo
significado depende de quem pergunta. `PostPublicList` tem query própria — e é a
query que torna o rascunho inalcançável, não apenas não listado. Filtro aplicado
depois da leitura é filtro que alguém remove arrumando o código, e a linha já
está em memória quando ele remove.

**`PostView` é perguntada duas vezes.** Uma sobre o valor zero — permissão para
olhar, que é o que produz o `Grant` que o repositório exige — e outra sobre a
linha que voltou. Decidir só uma vez é como uma policy que lê `p.PublishedAt` é
escrita e nunca consultada: o campo em que ela ramifica está sempre vazio.

## Consequências

- **Mudança compatível**: `apidiff` reporta `Guest: added` e `Subject.IsGuest:
  added`. Nenhuma aplicação existente muda de comportamento, porque nenhuma
  chama `Guest`.
- **O `examples` perdeu o único system grant.** O sitemap lê como convidado, pela
  mesma policy, então não pode listar o que a aplicação recusaria servir.
- **O site do `arandu-io` fica possível.** Era a mesma decisão bloqueando os dois.
- O que uma policy **não** pode fazer continua não podendo: `HasRole` responde
  falso para convidado, e não há id para comparar dono.
