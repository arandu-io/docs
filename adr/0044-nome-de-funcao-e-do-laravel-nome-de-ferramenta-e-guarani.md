# ADR 0044 — Nome de função é do Laravel; nome de ferramenta é guarani

**Status:** Aceita. Data: 10/08/2026.

## Contexto

Duas conversas estavam sendo tratadas como uma, e por isso as duas andavam mal.

De um lado, a REGRA 4 manda todo nome novo passar por cinco filtros, e a REGRA 7
diz que do Laravel se copia "estrutura e nomenclatura". Do outro, a REGRA 6 diz
que o código é em inglês.

Na prática isso virou hesitação em cima de cada identificador: `Cache::remember`
vira `Remember` ou `GetOrSet`? `Str::slug` vira `Slug` ou `Slugify`? E a
hesitação produziu o que uma medição contra o clone mostrou depois — pacotes com
o nome certo e um terço dos métodos, e alguns métodos renomeados por preferência
de quem escreveu.

O proprietário cortou a confusão em uma frase:

> "as funções que são usadas dentro do Illuminate, a gente não tem como inventar
> outro nome, porque é padrão dentro da nomenclatura do ecossistema"

E, sobre o outro lado:

> "a kyse é o nosso Blade... as ferramentas que a gente vai desenvolver, igual o
> Telescope, o Reverb — a gente vai pegar o nome dentro da cultura guarani, mas
> vai fazer basicamente o mesmo que a ferramenta do Laravel faz, dentro do nosso
> ecossistema"

## Decisão

**São duas categorias, com regras opostas, e a fronteira é se a coisa tem nome
próprio de produto.**

### 1. Superfície do framework — o nome é o do Laravel, sem exceção

Método, função, campo de struct, constante, nome de rota, chave de config,
nome de diretório do projeto. Tudo o que uma pessoa **digita** ou **procura**
por já conhecer do Laravel.

    Cache::remember   ->  cache.Remember       nunca GetOrSet
    Str::slug         ->  str.Slug             nunca Slugify
    paginate()        ->  Paginate             nunca Page
    $errors->first()  ->  Errors.First         nunca FieldError

Inventar aqui não é estilo, é quebrar a única coisa que faz alguém vindo do
Laravel reconhecer o que está lendo. A REGRA 4 **não se aplica** a esta
categoria: não há nome novo a filtrar, há um nome existente a respeitar.

As únicas alterações permitidas são mecânicas, e cada uma tem de estar dita no
comentário de doc: maiúscula inicial para exportar, `(T, error)` onde o PHP
lança, genérico onde o PHP usa `mixed`, e sufixo quando o nome é palavra
reservada do Go — `Support/Defer` é `support/deferpkg`.

### 2. Ferramenta com nome próprio — o nome é guarani

Produto que se instala, se cita e ganha logo. O equivalente do que a Laravel
chama Blade, artisan, Illuminate, Telescope, Horizon, Reverb, Scout, Sanctum.

Já nomeados: **arandu** (o framework), **aru** (o artisan), **kyse** (o Blade),
**hesape** (o Illuminate, ADR 0043).

A ferramenta faz **o mesmo** que a do Laravel faz, dentro do nosso ecossistema.
O nome é que é nosso.

Aqui a REGRA 4 se aplica inteira, e os cinco filtros na ordem: gíria e palavrão
em português brasileiro, categoria vetada, leitura em inglês, etimologia
atestada, colisão de repositório e domínio. O filtro zero já reprovou `karai`,
`kuatia` e `heka` — nenhum dos três foi pego pelo dicionário, só pelo ouvido.

## Consequências

- A hesitação em cima de identificador acaba: se está dentro do framework, o
  nome já existe e é o do Laravel.
- A REGRA 4 deixa de ser invocada onde não cabia, que é onde ela estava
  atravancando.
- Fica uma lista de ferramentas por nomear, cada uma com o equivalente Laravel
  declarado, e nenhuma delas nomeada antes de passar pelos cinco filtros:
  o console de depuração (Telescope), o painel de filas (Horizon), o servidor
  de WebSocket (Reverb), a busca (Scout — `heka` foi vetado no filtro zero),
  os tokens de API (Sanctum).
- O que a REGRA 7 chama de "estrutura e nomenclatura" ganha limite claro: a
  nomenclatura da superfície é copiada; o nome do produto não é.
