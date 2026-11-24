# Testes de console

- [Testes de console executam a fronteira de comando do projeto](#console-tests-execute-the-project-owned-command-boundary)
- [Despache um comando e inspecione seu efeito](#dispatch-one-command-and-inspect-its-effect)
- [Verifique argumentos, confirmação, efeitos e saída](#assert-arguments-confirmation-effects-and-exit)
- [Um nome registrado não comprova que a operação executa](#a-registered-name-is-not-proof-that-its-operation-runs)
- [Exercite o comando delegado pela Aru e pelo binário](#exercise-the-delegated-command-through-aru-and-the-binary)
- [Superfícies publicadas de testes de console](#released-console-testing-surfaces)
- [Associe comandos aos seus guias de recuperação operacional](#pair-commands-with-their-operational-recovery-guide)

<a name="console-tests-execute-the-project-owned-command-boundary"></a>
## Testes de console executam a fronteira de comando do projeto

Os comandos que conhecem a aplicação são implementados pelo binário do projeto e invocados pela Aru. A fronteira publicada da aplicação é `Application.HandleCommand(ctx, args)`, que responde o código de saída, e `WithConsole` lhe dá buffers no lugar dos streams do processo. Um teste útil verifica a saída, o código de saída e o efeito durável; um inventário de nomes de comando prova apenas que eles podem ser descobertos.

<a name="dispatch-one-command-and-inspect-its-effect"></a>
## Despache um comando e inspecione seu efeito

Defina uma configuração de aplicação isolada, construa a aplicação com o mesmo bootstrap que o binário usa, entregue-lhe buffers com `WithConsole`, chame `HandleCommand` com o comando e os argumentos exatos e inspecione o código de saída ao lado do resultado no banco de dados, na fila, no agendador ou no sistema de arquivos.

```go example
func runSchedule(ctx context.Context, t *testing.T, app *foundation.Application) {
	var out, errOut bytes.Buffer
	code := app.WithConsole(&out, &errOut, strings.NewReader("")).
		HandleCommand(ctx, []string{"schedule:run", "billing.close", "--tenant=acme"})
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut.String())
	}
}
```

<a name="assert-arguments-confirmation-effects-and-exit"></a>
## Verifique argumentos, confirmação, efeitos e saída

Teste argumentos inválidos, padrões contra flags, invocação repetida e o efeito real do adaptador selecionado. `HandleCommand` recebe o contexto que o teste possui e os streams que `WithConsole` definiu, inicializa a aplicação só para um comando que precisa dela e fecha o que abriu antes de retornar. Use um subprocesso só quando a própria fronteira do processo, um sinal ou a saída real do binário, for o contrato.

- Verifique o comando delegado e as flags exatas.
- Inspecione os efeitos persistidos.
- Leia a saída e os erros dos buffers de `WithConsole`.
- Execute comandos destrutivos apenas em dados descartáveis.

<a name="a-registered-name-is-not-proof-that-its-operation-runs"></a>
## Um nome registrado não comprova que a operação executa

Um comando pode aceitar a sintaxe e ainda ignorar um valor de configuração ou usar o padrão errado do construtor. Nomes desconhecidos, marcadores de projeto ausentes, operações de produção recusadas, locks retidos, dependências indisponíveis e escritas parciais precisam do status real de saída e sua evidência. Um teste que apenas lista nomes de comandos não substitui uma jornada que inicia, pausa, repete ou drena um worker; isso exige o armazenamento escolhido e a fronteira de processo.

<a name="exercise-the-delegated-command-through-aru-and-the-binary"></a>
## Exercite o comando delegado pela Aru e pelo binário

Para cada família de comandos, execute um efeito bem-sucedido, chamada malformada, falha de dependência, cancelamento e chamada repetida. Para comandos delegados da Aru, execute a CLI fixada contra um binário temporário do projeto e verifique argumentos encaminhados e código de saída.

- Padrões e substituições chegam às opções do runtime.
- Operações recusadas retornam erro.
- Chamadas repetidas preservam o contrato de idempotência.
- O efeito persistente escolhido é observável.

<a name="released-console-testing-surfaces"></a>
## Superfícies publicadas de testes de console

A Aru é dona da descoberta e do encaminhamento de nível superior; `Application.HandleCommand(ctx, args)` int é a fronteira publicada da aplicação para o comando do projeto, e `WithConsole(out, errOut, in)` define os streams em que ele escreve e de que lê. help e list imprimem a lista de comandos, e um nome desconhecido a imprime e responde 1.

| Superfície | Contrato publicado |
| --- | --- |
| `Application.HandleCommand` | contexto, slice de argumentos e código de saída |
| `Application.WithConsole` | stdout, stderr e stdin de um teste |
| estado durável | principal efeito observável |
| subprocesso | sinais e a saída real do binário |

<a name="pair-commands-with-their-operational-recovery-guide"></a>
## Associe comandos aos seus guias de recuperação operacional

A referência da CLI lista comandos, operações de fila definem efeitos dos workers, migrations definem efeitos no schema e implantação define a execução em produção.

- [Referência da CLI](/docs/{{version}}/cli/reference)
- [Operações de fila](/docs/{{version}}/operations/queues)
- [Migrations](/docs/{{version}}/database/migrations)
- [Implantação](/docs/{{version}}/operations/deployment)
