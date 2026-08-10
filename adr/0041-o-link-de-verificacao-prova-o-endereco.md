# ADR 0041 — O link de verificação prova o endereço, e o evento vai na mesma transação

**Status:** Aceita. Data: 10/08/2026.

## Contexto

Duas coisas em `framework/modules/auth`, achadas no mapeamento contra o
`laravel/ui`.

**O link assinava só o id.** O `sendVerification` publicado assinava `u.ID`, e o
`MarkVerified` lia a linha por id e carimbava sem comparar endereço nenhum. Então
um link mandado para o endereço que a pessoa tinha na hora continua verificando a
conta **depois** de o endereço mudar. Hoje é latente — o Service não tem método
que troque e-mail —, e a armadilha é essa: quem escrever "trocar meu e-mail" não
tem motivo nenhum para abrir o `signature.go`. O Laravel amarra os dois com
`hash_equals` contra o `sha1` do endereço atual.

**Nada podia assinar um cadastro.** O `framework/events` existe, tem outbox, e
está registrado nos dois bootstraps. `grep -n 'events\.' framework/modules/auth/`
não devolvia nada fora dos testes. O que havia era `col.RecordEvent(...)`, e o
comentário do próprio Collector diz que ele é o ring buffer de `/_arandu/debug` e
é nil fora de desenvolvimento. A forma era a pior das duas: um método chamado
`RecordEvent`, que se lê como dispatch, sendo buffer de debug.

## Decisão

**A assinatura carrega id e endereço**, com o mesmo cuidado que o propósito já
tinha: separador com prefixo de comprimento, senão o payload reconstrói
exatamente o bug que o `signature.go` existe para impedir — e existe um teste
nomeado por isso.

O `MarkVerified` recusa quando o endereço gravado difere do que veio no payload.
**Todo link já em caixa de entrada para de funcionar** na hora em que isto entra.
É o comportamento correto para uma mudança de significado, e é o custo que esta
ADR está registrando: quem tiver link pendente pede outro na tela de reenvio.

**O carimbo é condicional, e o evento sai da mesma escrita.** O `UserRepo.Confirm`
carimba `verified_at` só se ainda for nulo e reporta se foi ele quem carimbou.
Antes, `first` vinha de uma leitura anterior à escrita, já velha na hora de usar —
dois cliques simultâneos produziam duas confirmações e dois eventos.

**Falha ao gravar o evento desfaz a operação.** É o oposto do que uma falha de
envio de e-mail recebe, e de propósito: o outbox existe exatamente para que
"aconteceu" e "foi anunciado" sejam a mesma transação. Se o evento se perde
calado, o resto do sistema nunca fica sabendo e não há registro em lugar nenhum
de que algo faltou. Um rollback é recuperável; um evento perdido não é.

## Consequências

- `auth.Register` passou a não conseguir concluir sem a tabela de outbox. Os dois
  bootstraps publicados registram o `events.NewModule()`, então a tabela viaja —
  mas uma aplicação que remova aquela linha recebe `no such table: outbox` em todo
  cadastro, em tempo de execução, na tela onde um 500 é menos recuperável. Uma
  regra de `aru doctor` cobre isso.
- O `RecordEvent` continua: ele serve ao console de debug, que é outra pergunta.
- O `SetPassword` ainda tem a forma de lost update — lê, monta a linha inteira e
  grava, com a leitura fora da transação. É o mesmo defeito que o `Confirm`
  resolveu, num caminho diferente, e está na lista.
