# Processos

- [Execute um vetor de argumentos, nunca uma frase de shell](#execute-an-argument-vector-never-a-shell-sentence)
- [Execute um processo filho limitado](#run-one-bounded-child-process)
- [Controle explicitamente ambiente, diretório e saída](#control-environment-directory-and-output-explicitly)
- [Separe falhas de início, saída, timeout e captura](#separate-start-exit-timeout-and-output-failures)
- [Use binários de fixture, não comandos da máquina](#use-fixture-binaries-instead-of-host-commands)
- [Fronteira do componente de processos publicado](#released-process-component-boundary)
- [Prefira APIs Go e reserve processos para ferramentas reais](#prefer-go-apis-and-reserve-processes-for-real-tools)

> [!WARNING]
> **Limitações atuais.** A API de processos está publicada e testada, mas não tem consumidor de primeira parte. Nunca interpole entrada não confiável no shell; passe executável e argumentos explicitamente, defina prazo no contexto e limite a saída capturada.

<a name="execute-an-argument-vector-never-a-shell-sentence"></a>
## Execute um vetor de argumentos, nunca uma frase de shell

O componente de processos publicado envolve `os/exec` com contexto, vetores de argumentos, diretório de trabalho, ambiente herdado com substituições, saída capturada e informações de encerramento. Nenhuma aplicação atual com tag o consome; este guia limitado também explicita que stdout e stderr são retidos sem teto de bytes.

<a name="run-one-bounded-child-process"></a>
## Execute um processo filho limitado

Construa um vetor de argumentos, defina diretório, prazos total e de inatividade e apenas as substituições de ambiente necessárias. Env herda `os.Environ` e troca chaves correspondentes; não é uma allowlist. Nunca concatene entrada de requisição num comando shell.

```go example
func RunTool(ctx context.Context, workDir, binary, inputPath string, environmentOverrides map[string]string) error {
	result, err := process.NewFactory().NewPendingProcess().
		Path(workDir).
		Env(environmentOverrides).
		Timeout(2 * time.Minute).
		IdleTimeout(30 * time.Second).
		Run(ctx, []string{binary, "--input", inputPath}, nil)
	if err != nil {
		return err
	}
	if _, err := result.Throw(nil); err != nil {
		return err
	}
	return nil
}
```

<a name="control-environment-directory-and-output-explicitly"></a>
## Controle explicitamente ambiente, diretório e saída

Resolva caminhos de executáveis no boot, passe cada argumento como uma string, escolha diretório validado, audite o ambiente herdado e use Quietly ou `OutputHandler` quando reter saída ilimitada for inaceitável. Registre código de saída e duração sem segredos da saída.

- Passe argv sem sh `-c`.
- Defina prazos total e de inatividade.
- Audite variáveis de ambiente herdadas.
- Evite reter saída ilimitada de ferramentas verbosas.

<a name="separate-start-exit-timeout-and-output-failures"></a>
## Separe falhas de início, saída, timeout e captura

Executável não encontrado, permissão, falha ao iniciar, saída não zero, sinal, cancelamento, timeout total e de inatividade são distintos. Um filho iniciado pode ter efeitos externos parciais, e a saída capturada pode crescer sem teto da biblioteca.

> [!WARNING]
> Cancelar o contexto encerra o processo filho, mas não desfaz seus efeitos no filesystem ou em serviços externos. Timeout não prova que nada aconteceu.

<a name="use-fixture-binaries-instead-of-host-commands"></a>
## Use binários de fixture, não comandos da máquina

Construa pequenos binários de teste para sucesso, stderr, código de saída escolhido, espera, sinal, saída excessiva, eco de ambiente e limpeza de filhos. Evite asserções sobre utilitários de shell específicos da plataforma ou PATH do usuário.

- Argumentos com espaços permanecem como um único argumento.
- Timeout encerra o processo auxiliar.
- Substituições de ambiente prevalecem e os demais valores do pai são herdados.
- Saída verbosa segue a política de memória da aplicação.

<a name="released-process-component-boundary"></a>
## Fronteira do componente de processos publicado

`NewFactory`, `NewPendingProcess`, Path, Env, Timeout, `IdleTimeout`, Run, Quietly e `ProcessResult` compõem a fronteira publicada. Autorização, redução de ambiente, política de memória da saída e disponibilidade da ferramenta permanecem explícitas na aplicação.

| Superfície | Contrato publicado |
| --- | --- |
| `NewPendingProcess` | construtor de processos com vetor de argumentos |
| Path / Env | substituições de diretório e ambiente herdado |
| Timeout / `IdleTimeout` / Run | execução vinculada ao contexto |
| `ProcessResult` | saída integral capturada e metadados de encerramento |

<a name="prefer-go-apis-and-reserve-processes-for-real-tools"></a>
## Prefira APIs Go e reserve processos para ferramentas reais

Contexto possui cancelamento, filesystem possui arquivos, implantação possui ferramentas instaladas e contratos mantêm processos filhos atrás de um adapter específico da aplicação.

- [Contexto](/docs/{{version}}/architecture/context)
- [Sistema de arquivos](/docs/{{version}}/services/filesystem)
- [Implantação](/docs/{{version}}/operations/deployment)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
