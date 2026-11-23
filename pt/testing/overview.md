# Visão geral dos testes

- [Os testes inicializam a mesma aplicação sem servir uma porta](#tests-boot-the-same-application-graph-without-serving-a-port)
- [Gere as views e isole o estado externo](#generate-views-and-isolate-external-state)
- [Execute as verificações do projeto publicado](#run-the-released-project-gates)
- [Comprove o comportamento na menor fronteira completa](#prove-behavior-at-the-narrowest-complete-boundary)
- [Continue com HTTP, console, banco ou detalhes de contratos](#choose-http-console-database-or-contract-detail-next)

> [!WARNING]
> **Limitações atuais.** Os helpers de teste são publicados em vários módulos, sem uma fachada única. Componha testes Go comuns com o helper responsável pela fronteira e deixe explícita a preparação de cada pacote.

<a name="tests-boot-the-same-application-graph-without-serving-a-port"></a>
## Os testes inicializam a mesma aplicação sem servir uma porta

O Arandu usa testes Go comuns e helpers da aplicação. Um teste pode construir um serviço, inicializar a aplicação composta e despachar seu Handler ou executar um adapter de comando. A melhor fronteira é a menor que ainda comprova o comportamento, sem um framework de testes separado.

<a name="generate-views-and-isolate-external-state"></a>
## Gere as views e isole o estado externo

Execute `aru view:build` em um checkout limpo, pois os pacotes de views gerados podem estar no gitignore. Use fixtures descartáveis de banco e armazenamento compartilhado, transportes fake explícitos, relógios controlados onde disponíveis e restauração do ambiente por `t.Setenv`.

- Mantenha reproduzíveis as verificações de release com GOWORK=off.
- Isole banco e tenant por teste.
- Impeça requisições de rede inesperadas.
- Use `t.Cleanup` para todo recurso pertencente ao processo.

<a name="run-the-released-project-gates"></a>
## Execute as verificações do projeto publicado

Rode formatação, geração de views, vet, testes com race, build e o Doctor estrito a partir da raiz do projeto. Um parser ou teste unitário verde não basta quando o comportamento depende do bootstrap, de um driver selecionado ou de saída gerada. `aru test` roda `go test`, e `aru test --e2e` roda as jornadas de tests/E2E, que ficam atrás da build tag e2e, conduzem um navegador real contra a aplicação numa porta própria e são puladas, na máquina de quem desenvolve, onde não há navegador.

```sh
aru view:build
unformatted=$(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go' -exec gofmt -l {} +)
test -z "$unformatted" || { echo "$unformatted"; exit 1; }
GOWORK=off go vet ./...
GOWORK=off go test -race -count=1 ./...
GOWORK=off go build ./...
aru doctor --strict
```

<a name="prove-behavior-at-the-narrowest-complete-boundary"></a>
## Comprove o comportamento na menor fronteira completa

Use testes unitários para policies puras e transformações, contratos para cada adapter, features pela composição real da aplicação e fluxos semelhantes à produção para estado entre réplicas. Verifique resultados e recusas, não apenas parsing de valores ou retorno de construtores.

- Uma configuração alterada chega ao comportamento observável.
- Todo adapter de rede possui um teste contra requisições inesperadas.
- O isolamento do tenant é exercitado no armazenamento.
- Testes de race cobrem estado compartilhado e encerramento.

<a name="choose-http-console-database-or-contract-detail-next"></a>
## Continue com HTTP, console, banco ou detalhes de contratos

Testes HTTP despacham rotas, testes de console exercitam comandos delegados, testes de banco comprovam transações e Grants, e dublês capturam fronteiras de componentes sem estado global.

- [Testes HTTP](/docs/{{version}}/testing/http)
- [Testes de console](/docs/{{version}}/testing/console)
- [Testes de banco de dados](/docs/{{version}}/testing/database)
- [Dublês de teste](/docs/{{version}}/testing/test-doubles)
