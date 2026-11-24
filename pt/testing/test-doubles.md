# Dublês de teste e tempo

- [Um dublê de teste implementa o contrato do consumidor](#a-test-double-implements-the-consumer-owned-contract)
- [Injete um fake local e inspecione sua chamada](#inject-one-local-fake-and-inspect-its-call)
- [Prefira APIs de captura dos componentes publicados](#prefer-capture-apis-for-released-components)
- [Não substitua por fake o comportamento que está sendo comprovado](#do-not-fake-away-the-behavior-being-proved)
- [Execute o mesmo contrato no fake e no adapter real](#run-the-same-contract-against-fake-and-real-adapter)
- [Superfícies publicadas de captura e fakes](#released-capture-and-fake-surfaces)
- [Combine dublês específicos com fluxos completos de feature](#combine-narrow-doubles-with-complete-feature-journeys)

<a name="a-test-double-implements-the-consumer-owned-contract"></a>
## Um dublê de teste implementa o contrato do consumidor

Um fake fica junto do teste da menor interface exigida pelo consumidor. Registra chamadas tipadas, devolve resultados configurados e respeita contexto e semântica do Grant. Registros globais de substituição interferem em testes paralelos e devem ser evitados quando existe captura por instância.

<a name="inject-one-local-fake-and-inspect-its-call"></a>
## Injete um fake local e inspecione sua chamada

Defina um fake local para o colaborador do serviço, injete-o pelo construtor normal, invoque o objeto testado e verifique o Subject ou Grant exato derivado do contexto e a entrada tipada, sem comparar um mapa sem estrutura.

```go example
type invoiceSenderFake struct {
	grant security.Grant
	id    string
	err   error
}

func (f *invoiceSenderFake) Send(_ context.Context, g security.Grant, id string) error {
	f.grant, f.id = g, id
	return f.err
}

func TestFakeRecordsTheCall(t *testing.T) {
	fake := &invoiceSenderFake{}
	var grant security.Grant

	if err := fake.Send(context.Background(), grant, "inv_123"); err != nil {
		t.Fatal(err)
	}
	if fake.id != "inv_123" {
		t.Fatalf("id = %q", fake.id)
	}
}
```

<a name="prefer-capture-apis-for-released-components"></a>
## Prefira APIs de captura dos componentes publicados

Use fakes do cliente HTTP bloqueando requisições inesperadas, Capture de notificações por teste, transportes de e-mail capturados, adapters em memória apenas onde o contrato coincide e relógios ou IDs explícitos quando publicados. Cada captura pertence a um único teste.

- Faça chamadas inesperadas falharem.
- Registre entradas tipadas e a ordem das chamadas.
- Respeite cancelamentos e resultados de erro.
- Evite fakes mutáveis globais no pacote.

<a name="do-not-fake-away-the-behavior-being-proved"></a>
## Não substitua por fake o comportamento que está sendo comprovado

Um fake que sempre tem sucesso pode esconder tentativas, rollback, cancelamento e falha parcial. Um adapter em memória pode divergir em transações, ordenação, locks, TTL e semântica entre processos. Nunca use o dublê para provar o comportamento que ele omite intencionalmente. Fakes de parser ou de construtor não provam que um valor configurado chega ao componente em execução; acrescente uma jornada de feature que altere a configuração e observe o comportamento a jusante.

<a name="run-the-same-contract-against-fake-and-real-adapter"></a>
## Execute o mesmo contrato no fake e no adapter real

Execute um contrato comportamental compartilhado no fake e no adapter real quando houver promessa de substituição; depois adicione um teste de feature pela inicialização. Execute testes em paralelo e com detecção de race para revelar capturas globais ou callbacks inseguros.

- Chamadas inesperadas falham.
- Contexto e Grant são preservados.
- Falhas configuradas chegam ao objeto testado.
- Testes paralelos não compartilham capturas.

<a name="released-capture-and-fake-surfaces"></a>
## Superfícies publicadas de captura e fakes

Factory do cliente HTTP fornece Fake, `FakeSequence`, `StubUrl`, `PreventStrayRequests` e asserções de envio; `notifications.Capture` retorna canais e entregas por teste; os demais dublês implementam diretamente interfaces dos consumidores.

| Superfície | Contrato publicado |
| --- | --- |
| fakes de Factory do cliente | captura de requisições de saída |
| `notifications.Capture` | entregas por canal isoladas por teste |
| fake de interface local | fronteira do serviço ou repositório da aplicação |
| contrato compartilhado | paridade comportamental entre adapters |

<a name="combine-narrow-doubles-with-complete-feature-journeys"></a>
## Combine dublês específicos com fluxos completos de feature

Testes HTTP usam dublês nas fronteiras externas, testes de banco mantêm armazenamento concreto, contratos definem substituição e a visão geral escolhe a menor fronteira completa.

- [Testes HTTP](/docs/{{version}}/testing/http)
- [Testes de banco de dados](/docs/{{version}}/testing/database)
- [Contratos](/docs/{{version}}/architecture/contracts-adapters)
- [Visão geral dos testes](/docs/{{version}}/testing/overview)
