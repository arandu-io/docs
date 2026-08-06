# Política de documentação

## A regra

**Nada de guia, tutorial, referência em prosa ou site enquanto o código está se
movendo.** Documentação de uso é uma fase própria, depois do congelamento da API.

Motivo: guia escrito contra API instável é trabalho feito duas vezes, e a segunda
vez é pior — porque agora existe documentação errada publicada, que é mais
prejudicial que documentação ausente.

## O que NÃO é adiável

Três coisas não entram nessa regra, e é importante entender por quê:

### 1. Comentário de doc em Go

Não é documentação, é código. O `pkg.go.dev` gera a referência de API a partir dos
comentários no fonte, automaticamente. Escrever depois significa reabrir cada
arquivo para adivinhar o que cada função fazia.

Continua obrigatório em todo símbolo exportado, como já está no esqueleto.

### 2. ADRs

Dez linhas cada, e capturam o **porquê**. É a única categoria que não se
reconstrói depois: em três meses você lembra o que o código faz e esquece por que
rejeitou a alternativa.

Um ADR novo a cada decisão que fecha uma porta. Ver `adr/`.

### 3. Este conjunto de documentos (`docs/01` a `08`)

Não são documentação de produto, são **decisões**. Continuam vivos durante todo o
desenvolvimento. Quando a fase de documentação começar, eles são a fonte da qual
os guias serão escritos — não são os guias.

## O site é construído em cima do Arandu

Decisão tomada: `arandu-io/website` é uma aplicação Arandu. Não Next.js, não
Astro, não gerador estático de terceiro.

Isso não é vaidade, é a prova mais forte que existe. O `laravel.com` é uma
aplicação Laravel; o argumento "o framework funciona" fica vazio se o próprio site
do framework usar outra coisa.

E, mais importante que a prova: **é o primeiro cliente real.** Um site de
documentação precisa de roteamento, templates, busca, versionamento por branch,
cache e navegação — exatamente a superfície que um SaaS usa. Todo defeito de
ergonomia do framework aparece aqui antes de aparecer num cliente pagante.

### Consequência de cronograma

O site não pode existir antes de o framework servir HTML de verdade. Ou seja: ele
depende da camada de view (Fase 2). Então a fase de documentação vem
depois da Fase 3, não em paralelo.

### Divisão

| Repositório | Papel |
|---|---|
| `arandu-io/docs` | markdown puro, um branch por versão. Sem código. Espelho de `laravel/docs` |
| `arandu-io/website` | a aplicação Arandu que renderiza o `docs`. Consome o markdown, não o contém |

Manter separado importa: o markdown precisa ser editável por quem não sabe Go, e
o site precisa poder ser reescrito sem tocar no conteúdo.

## A fase de documentação

Entra no roadmap depois da Fase 3, como fase própria com critério de saída.

Escopo:

1. Congelar a API pública do core. Nenhuma mudança de assinatura durante a fase
2. Escrever os guias a partir de `docs/01` a `08`
3. Construir `arandu-io/website` sobre o Arandu
4. Revisar a página de origem do nome com quem estude tupi antigo
   (ver `04-marca.md`)
5. Revisar toda mensagem de erro do framework contra o tom de voz definido em
   `04-marca.md` — mensagem de erro é documentação, e é a mais lida de todas

Critério de saída: alguém de fora instala, cria um CRUD e sobe em produção
seguindo só o site, sem perguntar nada a você.

Esse critério é o único que prova que a documentação existe.
