# ADR 0042 — O link de recuperação de senha é assinado e não tem tabela

**Status:** Aceita. Data: 10/08/2026.
**Estende:** ADR 0032, que decidiu isto para o link de verificação e registrou o
reset como dívida — *"migrar aquele fluxo é outra mudança"*. É esta.

## Contexto

O reset de senha do starter kit guardava os tokens num mapa de pacote:

```go
var resets = struct { sync.Mutex; byHash map[string]reset }
```

Oito defeitos saíram do mesmo lugar, e o mapeamento contra o `laravel/ui` os
achou juntos:

1. Um restart invalidava todo link em voo, e com duas réplicas o link só
   funcionava na que o emitiu.
2. Sem limite de envio.
3. Mandava e-mail para endereço que nunca foi consultado. A guarda era
   `if email != ""` e o envio era incondicional — relé de spam saindo do
   domínio da própria aplicação. O Laravel chama `getUser()` e não manda nada.
4. O mapa crescia sem limite: cada endereço digitado inseria uma entrada que só
   saía quando aquele token exato fosse apresentado. Entrada expirada nunca era
   varrida.
5. Não apagava o token anterior: N envios, N links vivos.
6. O token não carregava tenant. Sob o resolvedor por host — que o
   `framework/modules/auth/module.go` chama o padrão da Fase 2 — um token emitido
   no tenant A e postado no host do tenant B trocava a senha da conta de B.
7. Consumia o token **antes** de validar a senha: senha curta queimava o link, e
   a re-renderização devolvia um token que já não existia.
8. Não abria sessão e não derrubava as outras — o próprio comentário do handler
   dizia que um reset que não desconecta as outras sessões deixa quem forçou o
   reset ainda conectado.

## Decisão

O token é assinado pelo `security.Signer` e **não é guardado em lugar nenhum**.
O payload carrega quatro campos, cada um com o comprimento na frente:

```
tenant | id | e-mail normalizado | impressão digital do hash da senha
```

O prefixo de comprimento não é enfeite: sem ele um valor contendo o separador
move a fronteira do próximo, que é exatamente o bug que o `signature.go` existe
para impedir — e lá já existe um teste nomeado por isso.

**A impressão digital é o que torna o link de uso único sem tabela.** É
SHA-256 do hash armazenado, truncado em 128 bits. Trocar a senha muda o hash,
que muda a impressão, e todo link emitido antes para de conferir. Não é a hash e
não é segredo: é função de mão única de um valor que já é função de mão única da
senha. Um link lido numa caixa de entrada não diz nada ao leitor além de se a
senha mudou desde então.

Isso fecha 1, 4, 5 e metade do 8 de uma vez, e transforma o 6 num campo.

O resto: o envio consulta a conta antes (3) com resposta **idêntica** nos dois
casos, porque a propriedade anti-enumeração é deliberada e tem teste; a
validação inteira acontece antes de qualquer consumo (7); o reset derruba as
outras sessões da conta pelo `DeleteSubject` que a ADR 0040 acrescentou (8).

**O orçamento de envio é contado à parte do de entrada** (2), e isso é a
decisão em vez de um descuido: quem acabou de errar a senha cinco vezes é
exatamente quem clica em "esqueci minha senha", e um contador compartilhado
recusaria a recuperação no único momento em que ela é pedida.

## Consequências

- Não existe tabela `password_resets`, e a migration equivalente do Laravel
  continua fora de escopo — agora por construção e não por dívida.
- Um link fica válido até a senha mudar ou o TTL expirar, o que vier primeiro.
  Não há como revogar um link individualmente sem trocar a senha; é o preço de
  não ter tabela, e é o mesmo preço que a ADR 0032 já aceitou na verificação.
- 128 bits truncados são muito mais do que o necessário para que nenhuma senha
  escolhível produza a mesma impressão, e curtos o bastante para o link caber
  numa linha. Link que o cliente de e-mail quebrou é link que chega quebrado, e
  este fluxo já foi mordido por isso.
