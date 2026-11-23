# Factories de models

- [Factories constroem registros tipados para testes e seeds](#factories-build-typed-rows-for-tests-and-seeds)
- [Configure o model antes de definir sua factory](#configure-the-model-before-defining-its-factory)
- [Gere dados determinísticos ou crie registros autorizados](#make-deterministic-data-or-create-authorized-rows)
- [Componha contagem, estado, sequência e relações](#compose-count-state-sequence-and-relationships)
- [Considere persistência parcial e erros de callbacks](#account-for-partial-persistence-and-callback-errors)
- [A criação ainda carrega o Grant do chamador](#creation-still-carries-the-caller-s-grant)
- [Comprove repetibilidade, isolamento e propriedade por tenant](#prove-repeatability-isolation-and-tenant-ownership)
- [Use factories em testes focados e seeders idempotentes](#use-factories-from-focused-tests-and-idempotent-seeders)

<a name="factories-build-typed-rows-for-tests-and-seeds"></a>
## Factories constroem registros tipados para testes e seeds

`aru make:factory` escreve uma factory tipada para uma entidade, como `factories.Users(db)` em `database/factories/UserFactory.go`, e `aru model:build` mantém atualizado tudo o que fica fora do bloco custom. Ela guarda um padrão válido e compõe estados determinísticos sobre a `factories.Factory` não genérica. A definição retorna User, estados recebem `*User` e os resultados são os mesmos registros dos terminais de consulta. Campos com nome errado falham na compilação em vez de desaparecer num mapa.

<a name="configure-the-model-before-defining-its-factory"></a>
## Configure o model antes de definir sua factory

Declare a entidade e sua tabela em app/Models e execute `aru model:build`. A factory tipada é construída sobre a consulta gerada da entidade: a tabela fornece as configurações de chave e a convenção de tenant, e o argumento db fornece a conexão; não há valor de model para construir à mão.

- Declare a entidade e sua tabela.
- Mantenha uma definição padrão válida.
- Use valores Faker determinísticos.
- Forneça Grant apenas às operações de persistência.

<a name="make-deterministic-data-or-create-authorized-rows"></a>
## Gere dados determinísticos ou crie registros autorizados

`MakeOne` não faz I/O nem exige contexto ou Grant; `CreateOne` grava pelo model com ambos. A semente padrão é 1 e `Seed(n)` torna falhas reproduzíveis.

```go
func defineUser(f faker.Faker) models.User {
	return models.User{Name: f.Name(), Email: f.Unique().Email(), Active: true}
}

made, err := factories.Users(db).MakeOne()
if err != nil { return err }
stored, err := factories.Users(db).Count(3).Create(ctx, grant)
if err != nil { return err }
```

<a name="compose-count-state-sequence-and-relationships"></a>
## Componha contagem, estado, sequência e relações

Count e Seed definem volume e repetibilidade; State executa em ordem; Sequence se repete em ciclos; `AfterMaking` e `AfterCreating` adicionam callbacks. A factory tipada repassa esses métodos; Has, que cria filhos depois dos pais, e `ForParent`, que cria primeiro um pai, ficam na factory do núcleo e são alcançados por um método no bloco custom. Cada modificador retorna uma nova factory.

- `Count(n)` e `Seed(n)` controlam um lote reproduzível.
- State e Sequence compõem alterações tipadas.
- `AfterCreating` recebe contexto e Grant.
- Has e `ForParent` exigem um callback explícito de associação.

<a name="account-for-partial-persistence-and-callback-errors"></a>
## Considere persistência parcial e erros de callbacks

Make pode falhar na atribuição ao model. Create salva em sequência: um erro posterior em Save ou em um callback pode deixar registros anteriores persistidos, pois Factory não envolve o lote inteiro em uma transação. Count negativo vira zero, e um callback de relacionamento válido, mas incorreto, gera dados enganosos.

> [!WARNING]
> Envolva a criação pela factory em uma transação controlada pelo chamador quando a atomicidade for importante. Verifique as chaves dos relacionamentos, sem presumir que um callback que compila conectou os registros pretendidos.

<a name="creation-still-carries-the-caller-s-grant"></a>
## A criação ainda carrega o Grant do chamador

> [!WARNING]
> Factory não contorna autorização. Create e callbacks de relação usam o Grant fornecido; nunca invente tenant nem use autoridade de sistema em requisições. Seeds de produção não contêm senhas fixas, tokens ou dados pessoais reais.

<a name="prove-repeatability-isolation-and-tenant-ownership"></a>
## Comprove repetibilidade, isolamento e propriedade por tenant

Comprove que a mesma semente repete e outra diverge, factories base não mudam, Make não chama banco, Create grava tenant e contagem esperados, Grants negados recusam, chaves de relação correspondem e falha parcial de callback é tratada deliberadamente.

- Valores com semente se repetem.
- Make não realiza E/S.
- Create preserva a autoridade do tenant.
- Modificadores não vazam entre testes.

<a name="use-factories-from-focused-tests-and-idempotent-seeders"></a>
## Use factories em testes focados e seeders idempotentes

Models definem a persistência, relacionamentos definem os vínculos, seeders controlam operações nomeadas de inicialização e testes de banco isolam cada conjunto de dados criado.

- [Models](/docs/{{version}}/database/models)
- [Relacionamentos](/docs/{{version}}/database/relationships)
- [Seeders](/docs/{{version}}/database/seeders)
- [Testes de banco de dados](/docs/{{version}}/testing/database)
