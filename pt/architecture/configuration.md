# Configuração

- [Configuração é uma entrada tipada da inicialização](#configuration-is-a-typed-boot-input)
- [Carregue uma vez e falhe antes de atender tráfego](#load-once-and-fail-before-serving-traffic)
- [Use um responsável para cada configuração](#use-one-owner-for-each-setting)
- [Rejeite de forma visível valores inválidos e sem efeito no runtime](#reject-invalid-and-unwired-values-visibly)
- [Rastreie um valor até o construtor que o utiliza](#trace-a-value-into-the-constructor-that-uses-it)
- [Loaders e configuração tipada](#loaders-and-typed-configuration)
- [A configuração escolhe adapters; o código mantém o contrato](#configuration-chooses-adapters-code-keeps-the-contract)

<a name="configuration-is-a-typed-boot-input"></a>
## Configuração é uma entrada tipada da inicialização

A aplicação lê a configuração uma vez, valida-a antes de abrir o servidor e passa valores tipados aos construtores. Handlers não consultam variáveis de ambiente, e componentes não buscam chaves de texto no momento do uso.

O framework é dono de toda configuração que constrói um driver. `bootstrap.LoadSettings` lê as variáveis de aplicação, HTTP, banco de dados, cache, sessão, fila, e-mail, sistema de arquivos e logging em `bootstrap.Settings`, e o `config.Config` do esqueleto incorpora essas configurações ao lado de suas próprias seções Auth, Geo e Services. Não leia uma variável do framework uma segunda vez em `config/`: dois leitores de um mesmo nome se afastam.

<a name="load-once-and-fail-before-serving-traffic"></a>
## Carregue uma vez e falhe antes de atender tráfego

`bootstrap/app.go` carrega a configuração, faz `foundation.NewDrivers` construir cada driver que as configurações nomeiam, constrói o provider da aplicação sobre esses drivers e entrega os mesmos drivers a `foundation.Configure`. Create devolve a aplicação composta ou a primeira contradição que encontrar.

```go
cfg, err := appconfig.Load()
if err != nil {
	return nil, err
}

drivers, err := foundation.NewDrivers(cfg.Settings)
if err != nil {
	return nil, err
}

app, err := providers.NewAppServiceProvider(cfg, drivers)
if err != nil {
	_ = drivers.Close()
	return nil, err
}

return foundation.Configure(drivers).WithProviders(app).Create()
```

> [!NOTE]
> Um driver muda mudando o ambiente, nunca no código. `CACHE_STORE`, `SESSION_DRIVER`, `QUEUE_CONNECTION`, `MAIL_URL` e `FILESYSTEM_DISK` nomeiam o que `NewDrivers` constrói, e um valor que nomeia um driver que o binário não vincula interrompe o boot com as duas linhas que o acrescentam.

<a name="use-one-owner-for-each-setting"></a>
## Use um responsável para cada configuração

| Família de configuração | Responsável no esqueleto | Unidade ou formato |
| --- | --- | --- |
| `APP_*`, `HTTP_*`, `TRUSTED_PROXIES` e observabilidade | `bootstrap.Settings` do framework | valores tipados |
| `DATABASE_URL` e `DB_*` | `bootstrap.Settings` do framework | uma URL; contagens do pool e segundos |
| `SESSION_DRIVER`, `SESSION_LIFETIME`, `SESSION_SECURE_COOKIE` e `CSRF_TTL` | `bootstrap.Settings` do framework | memory ou redis; minutos; booleano; segundos positivos |
| `CACHE_STORE`, `CACHE_PREFIX`, `CACHE_TTL` e `REDIS_*` | `bootstrap.Settings` do framework | memory ou redis; prefixo; segundos; uma URL e arquivos TLS |
| `QUEUE_CONNECTION`, `QUEUE_DEFAULT`, `QUEUE_WORKERS`, `QUEUE_RETRY_AFTER` e `QUEUE_MAX_ATTEMPTS` | `bootstrap.Settings` do framework | database ou redis; nome da fila; contagens positivas; retry em segundos |
| `MAIL_URL` e `MAIL_FROM_*` | `bootstrap.Settings` do framework | uma URL de transporte; remetente |
| `FILESYSTEM_*` | `bootstrap.Settings` do framework | local, r2 ou s3; raiz, bucket, endpoint, região e chaves |
| `LOG_LEVEL`, `LOG_FORMAT` e `LOG_SAMPLING` | `bootstrap.Settings` do framework | nível; json ou text; um a cada N |
| `ARANDU_TENANT_ID`, `AUTH_*`, `GEO_*` e credenciais de serviços | as seções da própria aplicação em `config/` | campos tipados |

O arquivo `.env` preenche valores ainda não fornecidos pelo ambiente do processo. Portanto, as variáveis do ambiente de publicação têm prioridade. Mantenha segredos fora do controle de versão e use o arquivo de exemplo apenas como catálogo de entradas possíveis, não como prova de que todas chegam a um consumidor. `LOG_SAMPLING` é lido e validado, mas o logger ainda não o aplica.

### `APP_ENV` decide o que a produção recusa

`APP_ENV` é dev, staging ou prod. O `.env` copiado de `.env.example` diz dev; um processo sem `APP_ENV` é prod, o que desliga o console de depuração, a pilha na página de erro e os comandos que esvaziam um banco de dados. Uma máquina de desenvolvimento declara `APP_ENV=dev`, e uma imagem de container define `APP_ENV=prod` em vez de deixar a resposta para uma variável não definida.

### Nomes de sessão não são aliases

A duração da sessão tem um único nome. O framework lê `SESSION_LIFETIME`, em minutos, e `foundation.NewDrivers` constrói o store de sessão a partir dele. `SESSION_TTL`, a contagem de segundos que o esqueleto lia antes, interrompe o boot com o valor de `SESSION_LIFETIME` a escrever no lugar.

<a name="reject-invalid-and-unwired-values-visibly"></a>
## Rejeite de forma visível valores inválidos e sem efeito no runtime

Booleanos malformados, contagens não positivas, URLs inválidas, chaves ausentes e nomes de store não suportados devem interromper o boot com a variável identificada no erro. Não substitua um valor explícito inválido pelo padrão, pois isso faz o processo executar uma configuração diferente da solicitada pelo operador.

> [!WARNING]
> A produção recusa dois valores no boot. `SESSION_SECURE_COOKIE=false` enviaria a sessão, a própria credencial, em texto claro; uma implantação servida por http é de staging e define `APP_ENV=staging`. `MAIL_URL=log://`, que também é o valor de um `MAIL_URL` não definido, escreveria cada mensagem, inclusive os códigos de redefinição de senha, no log; a produção define `smtp://`, `smtps://`, `resend://` ou `sendgrid://`.

No framework :frameworkVersion, o nome e o escopo do cookie de sessão não são configurações: ele é `arandu_session`, no caminho /, para o host que respondeu. `SESSION_PATH`, `SESSION_DOMAIN` ou `SESSION_SAME_SITE` pedindo outro escopo interrompem o boot, `SESSION_SECURE` o interrompe nomeando `SESSION_SECURE_COOKIE`, e `SESSION_COOKIE` não é lido por nada.

<a name="trace-a-value-into-the-constructor-that-uses-it"></a>
## Rastreie um valor até o construtor que o utiliza

Um teste de parser verifica apenas a sintaxe. Para cada configuração operacional, adicione ou encontre um teste que mude o ambiente, construa a aplicação e observe o componente consumidor. O teste de réplicas de sessão e os testes de configuração de workers de fila são o modelo: falham se a composição restaurar silenciosamente o valor padrão do construtor.

```sh
aru doctor --strict
go test ./tests/Unit ./tests/Feature
```

- Valores explícitos inválidos fazem o boot falhar e identificam a configuração.
- Um valor ausente produz o padrão documentado.
- Um valor alterado chega ao comportamento observável do componente.
- Segredos e credenciais de conexão são ocultados nos erros.

<a name="loaders-and-typed-configuration"></a>
## Loaders e configuração tipada

`bootstrap.LoadSettings` lê o ambiente uma vez e responde `bootstrap.Settings`, uma seção tipada por driver ao lado da configuração do próprio framework. No esqueleto, `config.Load` acrescenta as seções da aplicação e `config.From` faz o mesmo sobre configurações que um teste fornece. `foundation.NewDrivers` constrói o banco de dados, os stores de cache, os locks, o store de sessão, o emissor de CSRF, o rate limiter, a fila, o mailer e os discos a partir das configurações, e `foundation.Configure` compõe a aplicação sobre esses drivers.

- [Configuração de bootstrap do Framework](https://pkg.go.dev/github.com/arandu-io/framework/foundation/bootstrap)
- [Drivers e builder do framework](https://pkg.go.dev/github.com/arandu-io/framework/foundation)
- [Utilitários de configuração do Hesape](https://pkg.go.dev/github.com/arandu-io/hesape/config)

<a name="configuration-chooses-adapters-code-keeps-the-contract"></a>
## A configuração escolhe adapters; o código mantém o contrato

A escolha de store, driver ou transporte pertence à configuração, e a raiz de composição a constrói uma vez por `foundation.NewDrivers`. Os serviços da aplicação continuam dependendo da interface estável que recebem. Isso impede que a configuração vire um service locator e permite que um teste construa o mesmo grafo com substitutos explícitos.

- [Módulos e raiz de composição](/docs/{{version}}/architecture/modules-and-composition)
- [Sessão HTTP](/docs/{{version}}/http/session)
- [Contratos e adapters](/docs/{{version}}/architecture/contracts-adapters)
