# ADR 0036 — O recarregamento vem da conexão, não de uma mensagem

Data: 09/08/2026
Estado: Aceita

## Contexto

O `aru dev` reiniciava o servidor a cada mudança e o navegador nunca ficava
sabendo. Quem estava desenhando uma tela recarregava a mão, e a queixa registrada
foi: *"a pessoa vai modificando e vendo em tempo real a atualização no front. E
ele faz isso?"* — não fazia.

A saída óbvia é o processo antigo avisar o navegador antes de morrer. Ela não
existe: quando há o que avisar, o processo que avisaria já foi substituído.

Junto disso, o laço tinha sete defeitos confirmados por reprodução. Três
produziam exatamente o sintoma relatado — *"fico tendo que rodar `aru build` e
depois `aru dev`"*:

- `defer stopServer(server)` avaliava o argumento onde a linha estava escrita, e
  `server` era reatribuído a cada reinício. O `ctrl-C` matava o **primeiro**
  processo e deixava o vivo órfão, segurando a porta e servindo um build de horas
  atrás. É o único modo de falha em que nem o `aru build` resolve.
- Um `view:build` que falhava perdia a dívida: a mudança que o pediu já tinha
  sido consumida do snapshot, e o próximo `.go` salvo reiniciava contra views
  geradas antes da última edição. Só o `aru build` reparava.
- A saída do build era entrada do watcher. Uma edição em `resources/css/app.css`
  dava três reinícios, e nada no `aru` limitava isso — terminava porque o
  Tailwind não reescreve bytes idênticos, que é comportamento de binário de
  terceiro, não garantia.

E o primeiro comando de todo projeto novo falhava: `bootstrap/app.go` importa
`storage/framework/views`, que é saída de build e nasce vazia. `aru migrate`
respondia com um pacote que ninguém escreveu e não dá para obter.

## Decisão

**O sinal de recarregar é a identidade de quem responde a conexão, não uma
mensagem.**

Todo processo sorteia um `bootID` ao subir. O navegador ouve
`/_arandu/reload` por `EventSource`, que reconecta sozinho. Identidade diferente
da que a página carregou significa outro processo respondendo, e a página
recarrega. Rede que caiu e voltou no mesmo processo devolve a mesma identidade e
não recarrega nada.

Assim o mecanismo é consequência de como o `aru dev` já funciona, e não um
protocolo entre dois programas — uma coisa para manter em pé, não duas.

Decisões que vêm junto:

- **É middleware, não tag no layout.** Tag no layout é uma linha que todo projeto
  carrega, num arquivo que as pessoas editam, que não faz nada em produção e
  parece removível. O middleware monta com a rota, sob o mesmo `IsDev`, e um
  projeto que nunca ouviu falar disso ganha o comportamento.
- **Só injeta em documento com `</body>`.** Fragmento HTMX não tem, então uma
  troca nunca recebe um segundo `EventSource` — que acumularia uma conexão por
  interação até o limite de seis por origem travar a página.
- **O script é arquivo, não inline.** A CSP é `script-src 'self'` (REGRA 13), e
  tag inline é recusada por ela **em silêncio** — leria como a funcionalidade
  simplesmente não funcionar. Ele é registrado como asset comum, endereçado por
  conteúdo.
- **O script mora em `view/`, a rota no `kernel/`.** O `view` importa o `kernel`
  para ser um `Module`, então a seta só aponta para um lado: o kernel pede a tag
  por uma interface opcional e passa o endereço do próprio stream. Duas
  constantes para um endereço é como cliente e servidor passam a discordar dele.
- **Não existe hot-swap de CSS nem preservação de estado.** Isso exige runtime no
  cliente e grafo de módulos, que é o build de JavaScript que a REGRA 13 recusa.
  Recarga inteira contra um servidor Go que sobe em menos de um segundo é troca
  justa, e é um comportamento em vez de dois.
- **O `aru view:build --watch` foi removido.** Ele prometia a única coisa que não
  fazia: compilava as views uma vez e saía, porque entregava stdin fechado ao
  Tailwind. O `aru dev` é o laço, e duas implementações de watch são duas
  (REGRA 9).
- **`delegate` compila as views quando a saída não existe.** Não é um segundo
  caminho de build: é o mesmo `buildViews`, chamado quando o que o projeto
  importa não está lá. Nunca reconstrói o que está apenas desatualizado —
  desatualizado é o que o `aru dev` vigia; ausente é o que impede o projeto de
  compilar.

## Consequências

- `git clone && aru migrate` funciona. Era a primeira coisa que qualquer pessoa
  digitava, e falhava.
- Uma edição de CSS é um reinício, não três. Nada mais depende de o Tailwind não
  reescrever bytes iguais.
- `ctrl-C` não deixa órfão segurando a porta.
- Um `view:build` que falha continua devendo, e a próxima mudança paga.
- O `storage/` inteiro saiu do watcher: upload, cache e sessão são escrita do
  processo em execução, e reiniciavam o processo que acabara de servir a
  requisição.
- Asset embutido (`.svg`, `.png`, `.woff2`, `.toml`) passou a ser vigiado. Trocar
  um favicon não fazia nada, sem mensagem nenhuma, o que lia como bug de cache.
- Um `view:build` apaga o Go de uma view removida. O gerado se auto-registra, e o
  que ficava para trás mantinha a view apagada renderizável — e retomar o nome
  depois estourava no boot por registro duplicado. Era o caso em que o folclore
  mandava `rm -rf storage/framework/views`.

## O que não perseguir

`go:embed` mais o cache de build do Go. O **conteúdo** do arquivo embutido entra
no hash do action ID, e as entradas do `aru build` (`-trimpath`, `-ldflags`,
`CGO_ENABLED=0`) são disjuntas das do `go run` — então `aru build` não tem como
aquecer nada que o `aru dev` leia. É o primeiro suspeito óbvio e custa um dia.

## Alternativas descartadas

- **Watcher do sistema operacional** em vez de polling: dependência de terceiro
  numa CLI que não tem nenhuma, para ganhar meio segundo. Continua sendo a
  primeira coisa a trocar numa árvore muito grande, e só nela.
- **Checar a porta antes de subir**: o `aru` teria que aprender `HTTP_ADDR` do
  `.env`, virando segunda fonte de verdade de algo que a aplicação já é dona
  (REGRA 9). A linha `the application exited: …` resolve o mesmo caso e é
  autoexplicativa.
