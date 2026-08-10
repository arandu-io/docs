# ADR 0043 — A coleção de componentes se chama `hesape`

**Status:** Aceita quanto ao nome. **Recorte e endereço pendentes** da auditoria
contra os 38 componentes do Illuminate. Data: 10/08/2026.

## Contexto

No Laravel, `Illuminate` é o nome do conjunto de componentes: um repositório,
`laravel/framework`, com 37 `composer.json` sob `src/Illuminate/`, publicados por
subtree split para quem quer só `illuminate/http`. O desenvolvedor escreve
`use Illuminate\Http\Request;` todo dia.

Aqui o conjunto não tinha nome. Os pacotes moram sob `arandu-io/framework` e se
fala deles como "o framework", que é o mesmo nome do repositório e não distingue
o continente do conteúdo.

O proprietário pediu o nome, com o objetivo declarado: **quem vem do Laravel não
deve sofrer.**

## Decisão

O conjunto de componentes do Arandu se chama **`hesape`** — guarani para luz,
resplendor. É o que `Illuminate` é lá: o nome sob o qual os componentes vivem.

## O que ainda não está decidido, e por quê

**Onde o nome entra no caminho de import.** Três formas foram postas, com custos
diferentes — renomear o módulo para `arandu-io/hesape`, um diretório
`framework/hesape/`, ou o nome viver nos documentos e na marca com o código
seguindo em `framework/http`. A escolha depende de onde as fronteiras dos pacotes
caem, e isso é o que a auditoria contra os 38 componentes está apurando.

**Os nove pacotes cujo nome diverge do Illuminate** — `httpx`, `security`,
`data`, `jobs`, `kernel`, `observability`, `storage` e os demais — serão
renomeados, porque a REGRA 7 já manda copiar *estrutura e nomenclatura*. Não
agora: a auditoria responde se `httpx` são um pacote ou dois (`Http` e
`Routing`), se `security` são um ou três (`Auth`, `Hashing`, `Encryption`), e o
que ainda nem existe como pacote. Renomear antes é decidir o recorte no escuro e
refazer depois.

## Os filtros da REGRA 4, honestamente

Esta ADR existe em parte para que isto não se perca. O nome foi **aceito pelo
proprietário com três filtros em aberto**:

- **Filtro 0 — voz alta em português.** Não avaliado por quem tem o ouvido. Foi
  este filtro, e não a etimologia, que reprovou `karai`, `kuatia` e `heka`.
- **Filtro 2 — categoria vetada.** Luz e resplendor não são palavra neutra em
  guarani: estão no centro do Ayvu Rapyta, onde Ñamandu se cria em meio à
  escuridão pelo resplendor do próprio coração. A REGRA 4 veta divindade e ofício
  religioso, e a pendência *"revisar a página de origem do nome com quem estude
  tupi antigo"* já está no `checklist.sh`. Esta é a mesma checagem.
- **Filtro 3 — leitura em inglês.** `hesape` lido por um falante de inglês soa
  como *"he's ape"*. O produto é global e em inglês (REGRA 6, doc 04).

**Filtro 5:** `hesape.com` está ocupado; `.dev` e `.io` livres.

Se a checagem do filtro 2 reprovar, ou se o filtro 0 reprovar no ouvido de quem
fala, esta ADR é revogada e o nome muda antes de entrar em qualquer import — que
é exatamente por que o endereço ficou pendente em vez de ser decidido junto.

## Consequência de escopo

`hesape` **não é o décimo terceiro repositório** (REGRA 3). É o nome do conjunto
que já existe dentro de `arandu-io/framework`, do mesmo jeito que `Illuminate`
não é um repositório separado de `laravel/framework`.
