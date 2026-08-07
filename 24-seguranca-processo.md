# 24 — Segurança: o processo

Este documento não é sobre arquitetura de segurança — isso é o doc 01, o `Grant`
e o pipeline. É sobre **processo**: como uma vulnerabilidade é reportada,
tratada e divulgada.

Para qualquer framework isso é higiene. Para um cujo posicionamento é "o caminho
errado é impossível", não ter processo de vulnerabilidade é contradição
performática — e a cadeia de CVEs do Next.js mostra o custo de improvisar sob
pressão.

## Canal

`SECURITY.md` na raiz de todo repositório `arandu-io`, com:

- GitHub Security Advisory, que é o canal que funciona, nunca issue pública
- endereço dedicado (`security@` do domínio) — **publicado e sem entrega**, ver
  abaixo
- versões suportadas para correção de segurança (doc 23)

O `SECURITY.md` existe nos dez repositórios: `framework`, `aru`, `arandu`,
`database`, `kv`, `queue`, `storage`, `examples`, `homebrew-tap` e `docs`.

Issue pública sobre vulnerabilidade é fechada e movida para advisory privado,
com agradecimento — não com bronca.

### O `security@arandu.dev` não recebe mensagem

*Corrigido em 07/08/2026:* este documento anunciava o endereço dedicado como
canal primário e o GitHub Security Advisory como alternativo. É o contrário do
que funciona hoje, e o endereço não funciona de forma alguma: **o domínio
`arandu.dev` não tem registro MX.** `dig MX arandu.dev` devolve zero respostas —
só o SOA do registrador (`verify-contact-details.namecheap.com`). Mail enviado
para `security@arandu.dev` não chega a lugar nenhum, e os dez `SECURITY.md`
publicam esse endereço.

Um canal de vulnerabilidade que engole a mensagem é pior que canal ausente:
quem reporta acredita ter reportado e o embargo começa a correr sozinho.

A decisão é do proprietário, e são duas saídas: apontar MX para uma caixa que
alguém leia, ou remover o endereço dos dez `SECURITY.md` e deixar o GitHub
Security Advisory como canal único. Até que uma das duas aconteça, o canal real
é o advisory.

*Corrigido em 07/08/2026:* a linha "chave PGP publicada, opcional para quem
quiser usar" foi removida. **Não existe chave PGP.** Não há `.asc`, `.gpg` nem
impressão digital em repositório nenhum, e nenhum `SECURITY.md` menciona PGP —
a única ocorrência da palavra no projeto inteiro era esta linha.

## Prazos, escritos

| Etapa | Prazo |
|---|---|
| Acuso de recebimento | 72 horas |
| Triagem e severidade (CVSS) | 7 dias |
| Correção ou plano com data | 30 dias |
| Divulgação coordenada | 90 dias, ou antes se já houver correção |

Se a falha estiver sendo explorada, o prazo cai: correção e aviso imediatos, sem
embargo.

## Divulgação

Via GitHub Security Advisory, que emite CVE. O advisory sai junto com a versão
corrigida — nunca antes.

O texto do advisory diz, sempre: versões afetadas, versões corrigidas, impacto,
mitigação temporária e crédito a quem reportou.

## O que é e o que não é vulnerabilidade nossa

**É nossa:** qualquer forma de burlar o `Grant`, vazar dado entre tenants,
falha em CSRF, sessão, hash, rate limit ou headers, e qualquer código gerado
pelo `aru` que nasça inseguro.

**A última é a mais importante.** Se o gerador emite módulo com furo, o furo é
nosso em todos os projetos que geraram — e a correção precisa vir com comando de
regeração, não só com nota.

**Não é nossa:** aplicação que desligou o `doctor`, que usou `SystemGrant` fora
do escopo, ou que escreveu SQL à mão fora do repositório. Isso vira item de
documentação, não advisory.

## Módulo de terceiro

Vulnerabilidade em módulo da registry é responsabilidade do autor. Nossa parte:

- revogar o selo `auditado` até a correção
- marcar a versão afetada no índice
- notificar quem instalou, se houver telemetria de instalação — e não haverá
  telemetria silenciosa, então o canal é o índice e o `aru doctor`

*Corrigido em 07/08/2026:* esta seção descrevia o procedimento no presente, como
se estivesse em vigor. A `moitara` (doc 18) **não existe**: não há registry,
índice nem selo `auditado` em código nenhum — `grep -rni "auditado\|registry"`
nos `.go` de `aru`, `framework` e `arandu` só devolve o registro de seeders do
esqueleto e comentários sobre registry de container. O que está acima é a
política que passa a valer no dia em que a registry existir.

## Supply chain

- Core com **uma** dependência direta, `golang.org/x/crypto` (ADR 0004), mais o
  `golang.org/x/sys` indireto que vem com ela — essa é a defesa principal
- `govulncheck` na CI dos oito repositórios que têm código Go, reprovando a
  execução — e só ela: nada bloqueia merge hoje, ver abaixo
- `go.sum` verificado; nenhum `replace` apontando para fora em release
- Assets vendorizados (HTMX, Alpine) em `framework/view/assets/`; o binário do
  Tailwind é **pinado por versão** em `aru/internal/toolchain`, e o download é
  conferido contra o `sha256sums.txt` da mesma release

*Corrigido em 07/08/2026:* quatro linhas desta lista não descreviam o que existe.

A primeira dizia "core com duas dependências". O `framework/go.mod` tem um único
`require` direto, `golang.org/x/crypto v0.31.0`, e um `require // indirect`,
`golang.org/x/sys`. A contagem que vale é a da ADR 0004 e do `CLAUDE.md`: **uma**.

A segunda dizia "Dependabot com revisão humana obrigatória; nunca merge
automático", e foi removida: **não há Dependabot em repositório nenhum.** Não
existe `.github/dependabot.yml` em nenhum dos dez — o `.github/` de cada um tem
só `workflows/`. A frase descrevia a política de revisão de uma ferramenta que
nunca foi ligada. Adotá-la ou não é decisão do proprietário; o que não cabe é
listá-la como defesa em vigor. Enquanto isso, quem observa dependência vulnerável
é o `govulncheck`, que existe e roda em oito dos nove `ci.yml`.

A terceira dizia "o binário do Tailwind com hash conferido no PR de atualização",
o que sugere um hash commitado e revisado por humano. Não é o desenho: o que está
commitado é a **versão** (`DefaultTailwindVersion`), e é ela que passa por PR ao
subir. O hash vem do `sha256sums.txt` publicado na mesma release, o que prova que
os bytes chegaram íntegros e não que a release seja confiável — a garantia contra
troca silenciosa é o pino de versão, e o comentário de doc do próprio
`toolchain.go` já diz isso.

A quarta dizia "`govulncheck` na CI de todo repositório, bloqueando merge", e são
duas afirmações erradas em uma. Ele roda em **oito** dos dez repositórios: o
`docs` não tem `.github/workflows/` nenhum, e a CI do `homebrew-tap` é de fórmula
(`brew style`, `brew audit`, `brew install`), sem Go para analisar. E ele não
bloqueia merge: `gh api repos/arandu-io/<repo>/branches/main/protection` devolve
404 "Branch not protected" nos nove repositórios públicos, e no `docs`, que é
privado, o recurso nem está disponível no plano. Sem proteção de branch nada
impede push direto em `main` com a CI vermelha — o que existe é a execução
reprovando, que é visível, não impeditiva. Ligar proteção de branch com o
`govulncheck` como check obrigatório é decisão do proprietário, e é o que faria a
linha original virar verdade.

## O `doctor` faz parte da resposta a incidente

Quando uma falha revela um padrão perigoso, a correção tem duas partes: o patch
e **uma regra nova no `doctor`** que impede o padrão de voltar.

Isso é o que transforma incidente em aprendizado permanente — e é a mesma lógica
dos `hints` da página de erro. Regra nova de doctor é breaking change (doc 23),
então entra como aviso na minor e vira erro na major seguinte.
