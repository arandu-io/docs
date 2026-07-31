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
| Verificável | toda regra da doc tem checagem correspondente no `aru doctor` |
| Explícito | nenhum comportamento acontece sem estar em arquivo commitado |
| Diagnóstico | erro nunca é só erro: vem com a causa provável |
| Sóbrio | sem gamificação, sem emoji na CLI, sem "🚀 pronto!" |

## Tom de voz

O tom está na mensagem de erro, não no logo.

```
✗  ✨ Módulo criado com sucesso! Bom código! 🎉

✓  modulo cliente criado
   7 arquivos, 1 migração pendente
   policy nega tudo por default — edite cliente.policy.go antes de subir
```

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
| arandu.io | **ocupado** (monitor de tendências, BR) |
| arandu.app | **ocupado** (visualizador de Markdown, CLI `arandu`) |
| arandu.studio | **ocupado** (estúdio de UX/UI) |
| GitHub `@Arandu` | **ocupado** |
| GitHub `Arandusoft` | **ocupado** |

A fazer, e nesta ordem:

1. `arandu.dev` e `arandu.com` — registrar no mesmo dia
2. npm `arandu` e crates.io `arandu` — squatting defensivo, hoje
3. Org no GitHub: **`arandu-io`** (`@Arandu` e `Arandusoft` estão tomados).
   Módulo: `github.com/arandu-io/framework`.
4. Fórmula do Homebrew: registrar como `aru`, não `arandu`
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
