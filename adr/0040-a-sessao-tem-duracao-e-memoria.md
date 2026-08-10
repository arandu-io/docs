# ADR 0040 — A sessão tem duração, memória de confirmação e um índice por sujeito

**Status:** Aceita. Data: 10/08/2026.

## Contexto

O `SessionStore` tinha um único `ttl`, fixado no construtor e aplicado a todo
cookie. Nem `Start` nem `Rotate` recebiam duração, e o `SessionBackend` era
Get/Put/Delete por id, sem índice.

Três coisas dependiam disso e nenhuma podia existir:

**O "manter conectado" era um controle morto.** O `ui` publica
`<input type="checkbox" name="remember">` na tela de entrada desde sempre. Nada
lia o campo, nada atribuía `AuthPage.Remember`, e não havia forma na API por onde
pedir sessão mais longa — a caixa nem sobrevivia a um login rejeitado. Todo
projeto que rodou o kit publicou um controle que não faz nada.

**Não havia como pedir a senha de novo.** O `ConfirmsPasswords` do Laravel grava
`auth.password_confirmed_at` na sessão para que uma ação sensível possa exigir a
senha outra vez dentro de uma janela. Aqui não havia onde gravar, e por isso a
tela `auth/passwords/confirm.kyse.go` era publicada, compilada, registrada — e
inalcançável.

**O reset de senha não conseguia derrubar as outras sessões.** O próprio handler
do kit dizia em comentário que um reset que não desconecta as outras sessões
deixa quem forçou o reset ainda conectado. E não conseguia fazê-lo: sem índice
por sujeito, não havia como achar as sessões de uma conta.

## Decisão

**Duração por sessão, uma função só.** `Start` e `Rotate` ganharam opção
variádica — `security.Remember(true)` —, não uma segunda função ao lado. Uma
função parametrizada é uma forma; `StartFor()` ao lado de `Start()` seriam duas
(REGRA 9). Quem não passa opção nenhuma continua compilando sem edição, que é o
que mantém os goldens do `ui`, o `examples` e o `arandu` intactos.

Sessão lembrada é **mais longa, não eterna**. E o `MaxAge` do cookie e a validade
do registro são escritos do mesmo valor: cookie que sobrevive ao registro é
sessão que parece funcionar e não funciona.

**O carimbo de confirmação falha para o lado seguro.** Sessão gravada por binário
antigo, sem o campo, é **não confirmada** — nunca o contrário. É o mesmo
argumento que o `Subject.Verified` já carregava, e ele está citado no código para
não ser reinventado.

E quando o backend aceita o carimbo e não o guarda, o `Confirm` devolve
`ErrConfirmationNotStored` em vez de reportar uma confirmação que não existe.
Sem isso a tela de senha pediria a senha de novo imediatamente, para sempre, e
pareceria bug de formulário.

**O `SessionBackend` ganhou `DeleteSubject`.** Ele é interface exportada, então
acrescentar método é mudança quebradora para qualquer implementação de terceiro —
assumida aqui, porque a alternativa é o reset de senha continuar não conseguindo
fazer a única coisa que um reset precisa fazer. Tenant vazio ou sujeito vazio é
erro em toda implementação: uma varredura sem tenant é uma varredura que apaga a
sessão de outro cliente.

## Consequências

- O `kv` ainda não persiste `Remembered` nem `PasswordConfirmedAt`: os campos são
  mais novos que a versão do framework que o `go.mod` dele exige. Até o bump,
  `Confirm` sobre `kv` devolve `ErrConfirmationNotStored` em vez de mentir.
  **Pendência de release**, não de desenho.
- `Confirm` correndo com `DestroyOthers` pode ressuscitar uma sessão: o `Get` do
  primeiro cai antes do delete e o `Put` depois. É estreito — chegar lá exige a
  senha que acabou de ser trocada — e fechar exige compare-and-set no backend, que
  é um segundo método numa interface exportada na mesma release. Registrado aqui
  em vez de emendado às pressas.
