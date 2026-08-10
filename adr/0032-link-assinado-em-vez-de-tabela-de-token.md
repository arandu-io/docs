# ADR 0032-link-assinado-em-vez-de-tabela-de-token — verificação de e-mail

**Status:** Aceita.

## Contexto

O `examples` publicava as nove telas do starter kit e **duas delas não tinham
rota**: `auth/register.kyse.go` e `auth/verify.kyse.go` postavam para endereços
que ninguém registrava. A tela existia, o formulário existia, o `POST` dava 404.

Faltava, junto, a razão de ter verificação: `users` não tinha coluna de
confirmação, e nenhuma policy perguntava nada sobre isso.

O `PasswordController` já tinha um mecanismo de token — mapa em memória, chave
hasheada, TTL — e ele diz de si mesmo, no comentário, que está certo para uma
instância e errado para duas. Copiar aquele desenho para a verificação seria
copiar a limitação.

## Decisão

**O link de verificação é assinado, não armazenado.** `security.Signer` emite um
token HMAC sobre a chave da aplicação, carregando o id do usuário e um vencimento.

```
<payload-base64>.<vencimento>.<assinatura>
```

Nada é escrito quando a mensagem sai. Não há tabela, não há limpeza, e a
segunda réplica atrás do balanceador aceita o link que a primeira emitiu.

Duas propriedades fazem isso ser seguro numa URL:

- **O propósito entra na assinatura.** Um token emitido para confirmar endereço
  não vale numa troca de senha, embora a mesma chave assine os dois. E não é
  hipótese remota: os dois são "um link num e-mail com um id dentro".
- **O vencimento entra na assinatura.** Mover a data é mudar o payload, então um
  link vencido não se estende editando a URL.

O propósito é escrito com o **comprimento na frente** (`%d:%s|%s`). Sem isso,
propósito `verify` com corpo `|x.y` e propósito `verify|` com corpo `x.y`
produzem a mesma cadeia de bytes, e um token vale para dois propósitos.

### O que o link explicitamente não faz

**Não é uso único.** Isso exigiria estado, e onde uso único importa — troca de
senha — o estado é a própria senha: o token carrega o hash atual, então usar o
link muda aquilo contra o que ele foi assinado. Para confirmação de endereço,
clicar duas vezes é o caso comum (segundo cliente de e-mail, aba fechada), e
responder com erro se lê como link quebrado.

**Não abre sessão.** Confirma e manda para o login. Link de e-mail que autentica
é sessão que qualquer um com acesso àquela caixa — ou à mensagem encaminhada —
passa a ter.

### O que a verificação passou a valer

`security.Subject` ganhou `Verified`, escrito por `SubjectOf` no login, e
`CommentPolicy` recusa `CommentCreate` sem ele. Ler não precisa; escrever
precisa. É a barreira mais barata que existe contra um comentário automatizado:
conta é grátis, uma ida à caixa postal não é.

`Verified` fica no sujeito e não é lido do banco a cada requisição — a pergunta
é feita por policy, às vezes duas vezes por escrita, e ida ao banco dentro de
verificação de autorização é ida que ninguém enxerga do call site. O custo é que
o valor tem a idade da sessão, e isso falha para o lado seguro: conta nasce não
verificada, então sessão velha só pode ser **mais** restritiva que a verdade.

### Auto-registro passa pela policy

`Service.Register` autoriza com `security.Guest(tenant)` (ADR 0029) sobre
`ActionUserCreate`, e a `UserPolicy` responde olhando o **candidato**:

```go
if s.IsGuest() && len(u.Roles) == 0 { return nil }
```

Convidado cria usuário sem papel; o mesmo convidado pedindo `admin` é recusado
pela mesma linha. Escalação de privilégio pelo formulário de cadastro não é um
bug que se possa introduzir ali — campo novo em `RegisterRequest` continua
chegando nessa verificação.

Fechar o cadastro é apagar essas três linhas. Não existe segundo caminho que crie
usuário a partir de formulário.

## Consequências

- **Migration nova** em `modules/auth`, não edição da primeira: `name` e
  `verified_at`, ambas anuláveis (REGRA 16).
- `security.Signer` é geral. O próximo link em e-mail — descadastro, convite —
  usa o mesmo, com outro propósito.
- ~~O `PasswordController` continua com o mapa em memória. Migrar aquele fluxo é
  outra mudança, e o comentário dele já diz o que custa.~~ **Feito em
  10/08/2026, ADR 0042.** O mapa era a origem de oito defeitos, e o desenho que
  fechou cinco deles de uma vez é o desta ADR aplicado ao reset: token assinado,
  guardado em lugar nenhum, com uma impressão digital do hash da senha no payload
  para ser de uso único sem tabela.
- Teste de feature segue o link **que a aplicação mandou**, lido do transporte
  array. Foi ele que achou o kyse engolindo linha em branco: a URL saía colada
  na frase seguinte, e todo cliente que autolinka produzia token com três
  caracteres a mais.
