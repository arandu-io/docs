# Marca

## Nome

**Arandu** — tupi-guarani: sabedoria, inteligência, conhecimento profundo.
Entre os Kaiowá e Guarani, a leitura literal registrada é *"ouvir o tempo"*,
*"vivenciar"*, *"conhecer com a experiência de vida, na relação intrínseca com o
ambiente"*.

A história da marca para no significado linguístico. **Não usamos a figura
mitológica** (Tumé Arandú / Pay Sumé) como narrativa de produto: apoiar a marca
em ancestral sagrado é apropriação, mesmo quando a palavra em si é substantivo
comum. `arandu` entra como palavra do dicionário, não como nome de entidade.

A mesma regra vale para todo o ecossistema — ver a regra de admissibilidade em
`07-pacotes.md`.

## Por que o nome não descreve tecnologia

Não há "Go" nem "HTMX" no nome, de propósito. A plataforma pode mudar de runtime
— WebAssembly, outra linguagem — sem que a marca fique datada. Mesma escolha de
Laravel, Django e Rails.

## Posicionamento

> Não é o desenvolvedor que garante a arquitetura. É o compilador.

## Atitudes, e como se provam no produto

| Atitude | Prova verificável |
|---|---|
| Verificável | o que a doc decide tem quem cobre — o compilador, o `aru doctor` ou o `plans/checklist.sh` |
| Explícito | nenhum comportamento acontece sem estar em arquivo commitado |
| Diagnóstico | erro nunca é só erro: vem com a causa provável |
| Sóbrio | sem gamificação, sem emoji na CLI, sem "🚀 pronto!" |

*Corrigido em 07/08/2026:* a primeira linha prometia que "toda regra da doc tem
checagem correspondente no `aru doctor`". Não tem, e a promessa era maior que o
produto. O `doctor` lê AST, e as vinte e duas identificações de achado que ele
emite ficam no que se enxerga dentro de um arquivo Go: autorização, tenant, SQL,
sessão, view e manifesto, mais camada (`handler-reaches-data`,
`controller-reaches-repository`), dado sensível em log
(`sensitive-field-not-redacted`) e arquivo que não parseia
(`file-does-not-parse`). A REGRA 16 (migration não roda no boot) e a
REGRA 13 (nenhum arquivo de Node na árvore) são cobradas pelo
`plans/checklist.sh`, que olha o repositório inteiro e não um AST; config ser
struct tipada é cobrada pelo compilador, que é onde a tese do produto prefere
que fique. O que está verificado do `doctor` é a direção inversa, e essa vale:
toda regra que ele tem dispara em algum fixture, e todo nome de regra citado num
documento existe no código.

## Tom de voz

O tom está na mensagem de erro, não no logo.

```
✗  ✨ Módulo criado com sucesso! Bom código! 🎉

✓  cliente created, 12 files, in the Laravel tree.

   The policy denies every action. Open what this module needs in
   app/Policies/ClientePolicy.go, inside the custom block, and nothing else
   -- that is what makes the default safe.
```

*Corrigido em 07/08/2026:* o lado ✓ era um exemplo inventado, em português,
dizendo sete arquivos e mandando editar `cliente.policy.go`. O que está aí agora
é o que o comando imprime de verdade, copiado de uma execução de
`aru make:module cliente --fields "name:string!,email:email!u"` — sem a flag ele
para com `--fields is required` e sai com código 1: a saída da CLI é em inglês (REGRA 6), o gerador
escreve doze arquivos desde a ADR 0023, e a policy fica em
`app/Policies/ClientePolicy.go`. Um exemplo de tom de voz que ninguém pode
comparar com a saída real não prova tom nenhum.

## Nomes: marca vs. CLI

- Marca: **Arandu**
- CLI / binário: **`aru`**

O binário `arandu` já é distribuído por outro projeto (um visualizador de
Markdown, via Homebrew, snap e winget). Marca e comando divergentes são padrão:
Laravel/`artisan`, Kubernetes/`kubectl`, Angular/`ng`, .NET/`dotnet`.

`paje` foi considerado e vetado: pajé é ofício religioso. Se `aru` soar seco,
a alternativa tem que sair da mesma regra de admissibilidade — substantivo comum
da vida material, nunca cargo espiritual.

## Atribuição

O site deve ter uma página de origem do nome, creditando a língua tupi-guarani e
citando as fontes das traduções. Antes do lançamento público, a redação dessa
página deve ser revisada por alguém que estude tupi antigo ou nheengatu — a
referência lexicográfica é o *Dicionário de tupi antigo*, de Eduardo de Almeida
Navarro (2013).

Isso não é conformidade. É o que sustenta a história quando alguém perguntar, em
dez anos, por que o framework se chama assim.

## Checklist de namespace — estado em 29/07/2026

Verificado:

| Onde | `arandu` |
|---|---|
| npm | livre |
| PyPI | **ocupado** (pacote de memória para agentes de IA) |
| crates.io | livre |
| arandu.io | **nosso** desde 08/08/2026 — ver nota abaixo |
| arandu.app | **ocupado** (visualizador de Markdown, CLI `arandu`) |
| arandu.studio | **ocupado** (estúdio de UX/UI) |
| GitHub `@Arandu` | **ocupado** |
| GitHub `Arandusoft` | **ocupado** |

A fazer, e nesta ordem:

1. `arandu.dev` — consta **registrado** desde 09/02/2026, pela Namecheap, e os
   nameservers são `failed-whois-verification.namecheap.com` e
   `verify-contact-details.namecheap.com`. O RDAP não devolve `clientHold` nem
   `serverHold` — só `client transfer prohibited` — e o domínio resolve, então
   "suspenso" é inferência a partir do nome do nameserver, não medição. O que
   ela indica é verificação de contato pendente no registrador; responder a ela
   é ação de quem for o titular, e o RDAP não expõe quem é. `arandu.com` **não é
   registrável**: existe desde 05/01/2004 e está estacionado na Afternic — é
   compra de terceiro, não registro
2. npm `arandu` e crates.io `arandu` — squatting defensivo, hoje
3. ~~Org no GitHub: **`arandu-io`**~~ — **feito.** `@Arandu` e `Arandusoft`
   estavam tomados; a org é `arandu-io`, e os catorze repositórios vivem nela —
   treze públicos e o `docs`, que é privado. Módulo:
   `github.com/arandu-io/framework`
4. ~~Fórmula do Homebrew: registrar como `aru`, não `arandu`~~ — **feito.** É
   `Formula/aru.rb`, em `arandu-io/homebrew-tap`, e ela compila do fonte em vez
   de baixar binário
5. INPI classes 9 e 42, **em nome de HYZIS - Serviços Digitais LTDA**
   (CNPJ 57.283.262/0001-28), que é a titular do copyright — marca e código no
   mesmo titular, ou a separação vira problema no dia de uma aquisição. Buscar
   anterioridade antes: existem `arandu.studio` e `Arandusoft` em classes
   adjacentes, e vale opinião de advogado antes de investir em identidade visual.
6. **Não** usar ccTLD como vanity import path. `.io` é o código do Território
   Britânico do Oceano Índico, que deixa de existir pelo tratado de 2025 —
   se o código sair do ISO 3166-1, a ICANN inicia aposentadoria em cinco anos.
   Perder um domínio de marketing é chato; perder um import path quebra o
   `go build` de todos os clientes, para sempre.

*Corrigido em 07/08/2026:* os itens 3 e 4 estavam escritos como pendências e já
estavam cumpridos — a org existe e a fórmula está publicada. Os itens 2 e 5
continuam abertos: `arandu` segue sem dono no npm e no crates.io — o registro
responde 404 nos dois —, o que é o oposto do squatting defensivo que o item 2
pede.

O item 1 dizia "registrar `arandu.dev` e `arandu.com` no mesmo dia", e registrar
os dois já não é uma ação disponível. O RDAP em 07/08/2026 dá `arandu.dev`
criado em 09/02/2026 e `arandu.com` em 05/01/2004 — o RDAP não expõe o titular,
então de quem é o `.dev` não está verificado aqui, e o `.com` é anterior ao
projeto por vinte e dois anos. Sobre o `.com` o que resta é comprar ou não
comprar, que é outra pergunta e outro preço.

*Corrigido em 14/08/2026:* o item 3 dizia "os dez repositórios estão publicados
nela", e as duas metades estavam erradas. São catorze — `ui` (ADR 0026), `kyse`
(ADR 0027), `mcp` (ADR 0030) e `hesape` (ADR 0049) entraram depois — e nem todos
são públicos: `gh repo list arandu-io` devolve o `docs` como privado, além do
`plans`, que não é produto e não entra na conta.
