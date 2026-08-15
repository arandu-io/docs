# ADR 0051 — A configuração mora no componente, e o `framework` só a carrega

**Status:** Aceita. Data: 14/08/2026.

## A porta que se fecha

O `framework` deixa de poder ter um `config/` com dezessete arquivos espelhando
`laravel/framework/config/`. Não haverá struct central com uma seção por
domínio, nem lugar único onde toda configuração do projeto é declarada.

## Contexto

O plano da reconstrução deu ao `framework` um `config/` com dezessete arquivos,
um por arquivo de `laravel/framework/config/`: `app`, `auth`, `broadcasting`,
`cache`, `concurrency`, `cors`, `database`, `filesystems`, `hashing`, `images`,
`logging`, `mail`, `queue`, `services`, `session` e `view`.

Ao ir escrever, sete deles já existiam — e não onde o plano os punha:

| `laravel/framework/config/` | onde está aqui |
|---|---|
| `session.php` | `hesape/session.Config` — `Driver`, `Cookie`, `Lifetime`, `ExpireOnClose`, `Encrypt`, `Files` |
| `cache.php` | `hesape/cache.Config` — `Default`, `Prefix`, `Stores` |
| `logging.php` | `hesape/log.Config` — `Default`, `Channels` |
| `database.php` | `hesape/database.Config` |
| `filesystems.php` | `hesape/filesystem.Config` |
| `hashing.php` | `hesape/hashing.Config` |
| `broadcasting.php` | `hesape/broadcasting.Config` |
| `app.php` | `hesape/config.App` |

Mais `auth.php`, que é `hesape/auth.ManagerConfig`, e `mail.php`, que é
`hesape/mail.MailerConfig` e `mail.ManagerConfig`.

Escrever os dezessete no `framework` teria criado uma segunda declaração de cada
um desses dez, e a segunda seria a que o desenvolvedor encontraria primeiro.

## Por que o Laravel precisa dos dezessete e nós não

Em Laravel, `config/queue.php` existe porque **qualquer código, em qualquer
lugar, chama `config('queue.default')`**. O valor é buscado por chave em tempo de
execução, do `Illuminate\Config\Repository` que vive no container. Para isso
funcionar, algum lugar tem de declarar a chave — e esse lugar é o arquivo.

A [ADR 0001](0001-sem-container-di.md) rejeitou o container. Sem ele, ninguém
busca configuração por chave: **o componente recebe a sua no construtor**, já
tipada, e o compilador verifica o campo. Um `config/queue.go` no `framework` não
teria quem o lesse — a fila lê `queue.Config`, que é dela.

Isto é a [ADR 0045](0045-o-que-se-escreve-no-lugar-do-container.md) aplicada à
configuração, e o mesmo raciocínio de `docs/13`: a interface vive no pacote que
consome. A configuração também.

## Decisão

**Cada componente declara a própria `Config`, no seu pacote.** O `framework` não
replica nenhuma delas.

O que o `framework` faz é o que o Laravel chama de
`Foundation\Bootstrap\LoadConfiguration`: **ler o ambiente uma vez, no boot, e
montar as structs que os componentes recebem.** É Foundation, não é componente,
e por isso mora aqui (ADR 0049).

```go
// framework/foundation/bootstrap — o que existe no lugar dos dezessete
cfg, err := bootstrap.LoadConfiguration()   // le o ambiente, valida, falha cedo
// cfg.Session, cfg.Cache, cfg.Database ... cada um o tipo do seu componente
sessions := session.NewRecordStore(cfg.Session, handler)
```

O que se ganha em relação ao `config.Config` de onze campos que existe hoje: um
campo errado é erro de compilação no componente que o usa, e não um `nil` que
aparece no primeiro request que precisava dele.

## O que continua existindo, e por quê

**`hesape/config.Repository`.** É o `Illuminate\Config\Repository`, com os
mesmos métodos, e existe para o caso que a struct não serve — chave cujo nome só
se conhece em tempo de execução. Ele é um **leitor**, não um segundo depósito:
`hesape/config/doc.go` diz que nada de que o framework depende é lido por ele, e
uma chave escrita ali e em nenhum outro lugar não configura coisa alguma.

**`config-stubs/`.** O Laravel publica um `app.php` reduzido dentro do projeto.
O nosso equivalente é o `arandu/config/`, que já existe, no esqueleto, e é
editado — como a view, e pela mesma razão da [ADR 0021](0021-porang-dissolvido.md).

## O que falta escrever, e onde

Seis dos dezessete não têm `Config` em lugar nenhum: `concurrency`, `cors`,
`images`, `queue`, `services` e `view`. Cada um é da frente do seu componente e
nasce **no pacote dele**, nunca no `framework`:

| falta | onde nasce |
|---|---|
| `queue.Config` | `hesape/queue` |
| `view.Config` | `hesape/view` |
| `cors.Config` | `hesape/http/middleware`, junto do middleware que o lê |
| `concurrency.Config` | `hesape/concurrency` |
| `images.Config` | `hesape/image` |
| `services.Config` | não nasce — ver abaixo |

**`services.php` não tem equivalente e não vai ter.** Ele é um saco de chaves de
terceiros — Postmark, Resend, SES, Slack — que existe porque em PHP o pacote de
cada serviço lê `config('services.ses.key')` do container. Aqui o transporte
recebe a chave no construtor, de quem o constrói, e um lugar central para
credencial de terceiro é um lugar onde ela vaza para todo mundo que leia
configuração.

## Consequência

- O `config/` de dezessete arquivos sai da árvore alvo de
  `plans/framework/ESTRUTURA.md`. A regra de dono por domínio que ele carregava
  continua valendo, mas aponta para o pacote do componente.
- `framework/config.Config`, de onze campos, é substituído por
  `bootstrap.LoadConfiguration`, e a casca de `framework/config` some com as
  outras na v1.0.0.
- `docs/31:173` fica cumprido pelo espírito e não pela letra: ele mandava
  repartir os onze campos entre os componentes, o que é isto — mas o destino é a
  `Config` do componente, e não um arquivo por domínio no `framework`.
