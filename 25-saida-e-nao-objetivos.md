# 25 — Quando não usar, e como sair

O documento mais desconfortável de escrever, e o mais barato em credibilidade.

Os três projetos vivos do nicho — Goldr, Goilerplate e Pagoda — se posicionam
exatamente contra o que somos: eles vendem posse do código e ausência de lock-in.
E o risco de abandono num projeto de uma pessoa é máximo: o Gorilla parou de ser
mantido pelo time original em nov/2023 e 17% do ecossistema Go seguiu pendurado.

Dizer isso primeiro é mais forte do que deixar o cliente descobrir depois.

## Quando **não** usar o Arandu

| Caso | Por quê |
|---|---|
| **Aplicação offline-first** | HTMX depende de hipermídia vinda do servidor. Sem rede, não funciona |
| **Interatividade de alta frequência** | Editor gráfico, canvas, jogo, ferramenta de desenho. Latência de rede é o teto |
| **Colaboração em tempo real como núcleo** | Cursor compartilhado, edição simultânea. Times escolheram hipermídia para isso e passaram meses costurando WebSocket e estado no cliente |
| **App mobile nativo como cliente principal** | Dá para usar em modo API, mas aí você está pagando por metade do framework que não vai usar |
| **Equipe que não quer arquitetura imposta** | A cultura Go tende a rejeitar isso. Se o time briga com a Policy obrigatória, vai brigar todo dia |
| **Módulo com JOIN complexo no perfil de performance** | Não é elegível. O `doctor` reprova, e isso é por desenho (doc 12) |

Nada disso é defeito a corrigir. É fronteira.

## O que é seu e o que é nosso

Se você parar de usar o Arandu amanhã, isto continua sendo seu, e é Go comum:

- handlers, services, policies, entities, requests — código Go no seu repositório
- templates: você escreve `.kyse.go` e o `kyse` gera o `.go` ao lado, no mesmo
  pacote. O gerado é Go comum, está na árvore e é seu — mas é **gerado**, não
  escrito. Se você sair, mantém o `.go` como código final e abandona o fonte
  `.kyse.go`, ou reescreve as views em `html/template`. Estimativa: uma semana
  num projeto médio
- SQL: as queries estão em `.sql` e o `sqlc` gera Go a partir delas
- migrações: SQL versionado
- CSS: Tailwind v4 compilado por binário standalone. Sem npm e sem runtime

Isto some se você tirar o framework:

- `kernel` — boot, registro de módulos, pipeline
- `security.Grant` — a garantia de autorização
- `observability.Collector` e a página de erro
- `aru` — CLI, gerador e `doctor`
- `framework/view` — runtime de render, assets embutidos e a rota que os serve
- `kyse` — o motor de template, e ele **não existe fora do Arandu**

As views em si (`resources/views/`) ficam: elas moram no seu projeto, não em
repositório nosso.

## Roteiro de saída, honesto

1. **Trocar o kernel por `main.go` próprio.** O pipeline é
   `func(http.Handler) http.Handler` padrão e o router é `http.ServeMux`. Um dia.
2. **Substituir o `Grant`.** É a parte cara: as assinaturas de repositório
   exigem `Grant`, então ou você mantém o pacote `security` (MIT, pode
   copiar) ou remove o parâmetro de todas as assinaturas. Uma a duas semanas num
   projeto médio, e é mecânico.
3. **Substituir o `Collector`.** Trocar por OpenTelemetry. Alguns dias.
4. **Substituir o `kyse`.** Aqui está o custo novo, e ele é honesto: o `kyse`
   é nosso e não tem equivalente. Ou você mantém os `.go` gerados como código
   final e para de editar os `.kyse.go`, ou reescreve as views. Uma a duas
   semanas.
5. **Manter o resto.** Pacotes, sqlc e migrações não precisam mudar.

Estimativa honesta para um projeto de porte médio: **três a quatro semanas**, e
nenhuma reescrita de regra de negócio.

## O que fazemos para a saída ser barata

Não é gentileza; é decisão de arquitetura, e cada item já está travado:

- **Sem geração em runtime e sem reflection** — o que roda é o que está escrito
- **Wiring explícito e commitado** — nada a "descobrir" na saída
- **DSL não sobrevive no código final** — o que fica é Go legível (doc 19)
- **Contrato de módulo é uma função**, não uma hierarquia
- **`net/http` puro por baixo** — todo o ecossistema Go continua compatível
- **Sem CDN, sem serviço nosso em runtime** — a aplicação não liga para nós

Um framework que precisa te prender para sobreviver não deveria sobreviver.

## Se o projeto for abandonado

Compromisso escrito, porque é a pergunta que todo comprador sério faz:

- A licença é MIT (doc 20, ADR 0008). Fork é direito de qualquer um, sem pedir.
- Se a manutenção parar, publicamos aviso no README e no índice da registry,
  com tag final e estado conhecido dos bugs abertos.
- Não haverá "morte silenciosa": repositório sem commit e sem aviso por seis
  meses é considerado abandonado por esta política, e o aviso é obrigação nossa.
