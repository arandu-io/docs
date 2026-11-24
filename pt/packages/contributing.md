# Contribuição

- [As verificações do repositório definem uma alteração revisável](#repository-gates-define-a-reviewable-change)
- [Leia o contrato do repositório antes de editar](#read-the-repository-contract-before-editing)
- [Faça uma alteração focada com seu teste de regressão](#make-one-focused-change-and-its-regression)
- [Prove que cada nova proteção pode falhar](#prove-every-new-guard-can-fail)
- [Mantenha cada alteração associada a um motivo](#keep-one-change-attributable-to-one-reason)

<a name="repository-gates-define-a-reviewable-change"></a>
## As verificações do repositório definem uma alteração revisável

Formatação, vet e testes com detector de corrida são a base comum, complementada pelas verificações de geração, dependências, vulnerabilidades e conformidade de cada repositório. Os comandos rodam localmente e na CI; uma contribuição informa seus resultados reais, sem tratar os nomes dos comandos como aprovação.

```bash
unformatted=$(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go' -exec gofmt -l {} +)
test -z "$unformatted" || { echo "$unformatted"; exit 1; }
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 -p 1 ./...
```

<a name="read-the-repository-contract-before-editing"></a>
## Leia o contrato do repositório antes de editar

Use a versão do Go e a política de workspace declaradas pelo repositório, confira quem é responsável pelos arquivos gerados e identifique os testes específicos do comportamento alterado. Relatos de segurança seguem a política de comunicação privada; uma issue pública não é o primeiro lugar para divulgar uma vulnerabilidade.

<a name="make-one-focused-change-and-its-regression"></a>
## Faça uma alteração focada com seu teste de regressão

Reproduza a falha, adicione um teste que observe o comportamento incorreto, implemente a menor correção coerente e execute o teste novamente. Atualize a documentação pública quando a API observável mudar. Deixe formatação sem relação e alterações desnecessárias de arquivos gerados fora do patch para preservar a clareza da revisão.

<a name="prove-every-new-guard-can-fail"></a>
## Prove que cada nova proteção pode falhar

Introduza a violação, execute a proteção específica, observe a falha, restaure o código e observe o sucesso. Um teste de varredura também verifica se inspecionou uma quantidade plausível de alvos, para que um padrão desatualizado não passe sem encontrar nada. Exclua templates de views e fixtures intencionalmente inválidas das verificações genéricas de formatação pelos motivos documentados.

<a name="keep-one-change-attributable-to-one-reason"></a>
## Mantenha cada alteração associada a um motivo

Separe alterações mecânicas de mudanças de comportamento quando isso facilitar a revisão das evidências. A mensagem de commit explica o motivo que o diff não permite reconstruir. Instruções locais, responsabilidade pelos arquivos gerados, regras de release e procedimentos para relatar problemas de segurança continuam valendo além dos comandos comuns.
