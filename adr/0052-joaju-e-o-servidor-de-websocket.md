# ADR 0052 — `joaju`, o servidor de WebSocket, sobre um fork do `gorilla/websocket`

**Status:** Aceita. Data: 14/08/2026.

## A porta que se fecha

O Arandu passa a ter servidor de WebSocket próprio, e deixa de poder responder
"use SSE" a quem precisa de bidirecional. Também deixa de valer a escolha de
driver que `docs/11-adapters.md` tinha feito.

## Contexto

`docs/06-escopo.md:150` punha o WebSocket em "com gatilho — quando o SSE do HTMX
não bastar". O gatilho foi puxado pelo proprietário: o recurso vai ser lançado.

O alvo é o `laravel/reverb`, e o que ele é, verificado na documentação em
14/08/2026:

- implementa o **protocolo Pusher** — os eventos internos dele recebem
  `Laravel\Reverb\Protocols\Pusher\Channel`;
- roda sobre um **event loop do ReactPHP**, que por padrão usa `stream_select` e
  trava em 1.024 conexões; passar disso exige a extensão `ext-uv`, instalada por
  PECL;
- **escala por Redis pub/sub**: cada servidor publica no Redis o que recebeu, e
  os outros entregam aos seus conectados;
- serve **várias aplicações** de uma instância, com `app_id`, `app_key` e
  `app_secret` por aplicação, e uma lista de origens permitidas;
- escuta WebSocket em `/app` e API em `/apps`.

## Decisão

**`arandu-io/joaju`** é o décimo quinto repositório, e é o servidor de WebSocket.
Esta ADR é a autorização que a REGRA 3 exige.

Ele nasce de um **fork do `gorilla/websocket`** (BSD-2-Clause, mantido, API
estável), e implementa o **protocolo Pusher**, para que a estrutura do Reverb
seja reconhecível.

## O que isto revoga

`docs/11-adapters.md:119-126` escolheu **`coder/websocket`** como padrão e
`gobwas/ws` para alta concorrência. **Fica revogado.** A razão de então era
"idiomático, próximo da stdlib", e ela continua verdadeira — o que mudou é que
não estamos escolhendo uma biblioteca para consumir, e sim uma base para
manter. Para isso o critério é outro: `gorilla/websocket` é a implementação mais
usada do ecossistema, tem a API estável há anos, e é a que qualquer pessoa que
chegue ao repositório já conhece.

O que **não** muda: Socket.IO e µWebSockets.js continuam recusados, e pela mesma
frase — são Node, e adotá-los transforma uma aplicação Arandu em dois processos e
dois runtimes.

## O que o fork obriga, e está aceito

Um fork não é uma dependência. O que se ganha é poder mudar o núcleo — e o que
se paga é que **toda correção de segurança do `gorilla/websocket` passa a ser
nossa para acompanhar**. Isso é trabalho permanente, e é a razão de a ADR 0004
manter o core com uma dependência só.

Fica dentro do `joaju`, que é repositório separado: o core continua sem.

O `THIRD_PARTY.md` do `joaju` carrega a licença BSD-2-Clause e a atribuição, do
mesmo jeito que a ADR 0034 fez com as fontes e a REGRA 7 com o `basecoat`.

## O cliente, que é a decisão de verdade

O protocolo Pusher só serve se algo o falar do lado do navegador. No Laravel isso
é `laravel-echo` mais `pusher-js`, e `docs/31:398` os lista entre os recusados.

**A recusa de lá é sobre npm no grafo do projeto, não sobre JavaScript no
navegador**, e essa diferença é o que resolve. A REGRA 13 é literal: *"HTMX,
Alpine e CSS entram por `go:embed`. Sem CDN — a CSP é `script-src 'self'`."*

Então o cliente entra do mesmo jeito que o HTMX já entra: **vendorizado, servido
pelo próprio binário, com a licença junto**. Não há `package.json`, não há
`node_modules`, não há CDN, e `git clone && aru dev` continua sendo tudo.

O que **não** entra é o `laravel-echo`: ele é a camada que resolve canal contra o
container do Laravel, e não há container (ADR 0001). O que a aplicação escreve é
um atributo HTMX ou uma função do cliente vendorizado.

`sabiá` era o nome reservado para o cliente e **fica livre** — não há pacote
nosso a nomear se o cliente é vendorizado.

## O que fica decidido junto, para não voltar como pergunta

**Escala por Redis pub/sub, como o Reverb, e o pré-requisito não existe.**

Esta ADR dizia, na primeira versão, que `hesape/redis` **não tem `Publish` nem
`Subscribe`** — repetindo `docs/31:527`, que afirma "verificado, zero
ocorrências". **As duas afirmações são falsas.** Os três estão escritos:

	redis/connections/commands.go:506       Publish(ctx, channel, message) (int64, error)
	redis/connections/subscriptions.go:27   Subscribe(ctx, channels, callback) error
	redis/connections/subscriptions.go:37   PSubscribe(ctx, channels, callback) error

A escala horizontal é possível hoje, sem trabalho anterior. A correção fica
escrita em vez de o parágrafo ser apagado, porque a afirmação errada foi copiada
de um documento para outro e a próxima cópia sairia do mesmo lugar.

**Não há event loop a escolher.** O `ext-uv` do Reverb existe porque PHP não tem
concorrência; o `stream_select` dele trava em 1.024 conexões. Em Go uma goroutine
por conexão é o modelo, e o limite é descritor de arquivo, não biblioteca.

**Multi-aplicação fica de fora da primeira versão.** O Reverb serve várias
aplicações de uma instância porque hospedar PHP é caro. Um binário Go é um
processo: subir um por aplicação é mais simples de operar e não tem estado
cruzado entre clientes — que é exatamente o que a REGRA 14 protege.

**O `Grant` vale aqui.** Assinar um canal é leitura, e a REGRA 17 não abre
exceção: nenhuma inscrição é aceita sem passar por uma Policy, e o tenant vem do
`Grant` e nunca do canal que o cliente mandou. O `hesape/broadcasting` já tem
esse buraco fechado — `LogBroadcaster.Auth` e `NullBroadcaster.Auth` devolviam
`nil` e autorizavam qualquer um.

## O que esta ADR não decide

**HTTP/1.1 versus HTTP/2**, que `docs/31:561` levanta e continua aberto: um
`EventSource` gasta uma das seis conexões que o HTTP/1.1 dá por origem, e
`aru dev` serve HTTP/1.1 em texto claro. Para WebSocket o problema é outro — o
upgrade sai do pool — mas a decisão de o `aru dev` servir HTTP/2 continua de pé e
é anterior a este repositório.

**O `potiron`**, dashboard de filas, e o `moitara` seguem sem data. Este ADR
autoriza um repositório, não a lista inteira.
