# ADR 0024-sqlc-nao-entra — o SQL é escrito à mão, e o `sqlc` vira gatilho

**Status:** Aceita.
**Supersede parcialmente:** a promessa de `sqlc` na REGRA 9 e em doze pontos do
registro, listados abaixo.

## Contexto

Doze lugares do registro afirmam que o caminho de dados do Arandu é
`sqlc` + `Repository` — inclusive a REGRA 9 do `CLAUDE.md`, que o cita como
*a* forma única de chegar ao banco, `docs/25:35`, que promete ao comprador que
"as queries estão em `.sql`", e `docs/22:234`, que apresenta o `sqlc` como a
resposta validada ao N+1.

O `sqlc` nunca entrou. O SQL é escrito à mão nos templates do gerador, em
`aru/internal/gen/`, e é isso que todo módulo gerado carrega desde a Fase 2.

A divergência sobreviveu porque nada a mede: o `sqlc` não é um arquivo que falta
nem um pacote que não compila, e o checklist verifica exatamente essas duas
coisas. Ela apareceu num levantamento de decisão-registrada-contra-código.

## Decisão

**O SQL continua escrito à mão, parametrizado, nos templates do gerador.** O
`sqlc` sai de promessa e vira gatilho, na forma do `docs/06-escopo.md`.

O gatilho: **quando o perfil de performance (`docs/12`) exigir a segunda emissão
de queries.** Ali passam a existir dois conjuntos de queries para uma API — o
convencional e o de agregado — e um gerador de queries deixa de ser conveniência
e passa a ser o que impede as duas de divergirem. Antes disso não é.

Os doze pontos do registro passam a dizer "SQL parametrizado escrito à mão".

## Por que não agora

**Não há o problema que o `sqlc` resolve.** Ele existe para dar tipo ao
resultado de query escrita à mão e para pegar erro de SQL em tempo de build.
O primeiro já está resolvido: o gerador emite a struct e o `Scan` juntos, da
mesma spec, então eles não podem divergir — não há query escrita à mão sem tipo
correspondente. O segundo está coberto pela suíte de conformidade da ADR 0018,
que roda as queries geradas contra PostgreSQL e MySQL de verdade.

**Não há risco de injeção.** A concatenação que existe nos templates usa apenas
constante e nome de coluna de allowlist; todo valor vai em placeholder. O
`aru doctor` tem a regra `sql-built-with-sprintf`, e ela foi ampliada para
concatenação de variável livre.

**O custo cairia inteiro no usuário.** O `sqlc` é um binário externo com
`sqlc.yaml`, passo de geração próprio e um segundo conjunto de arquivos gerados.
Um projeto Arandu roda com `git clone && aru dev` (REGRA 13) — o compromisso é
sobre Node, mas a razão dele é o número de ferramentas entre clonar e rodar, e
o `sqlc` é mais uma.

**E seria a segunda virada grande em duas semanas.** A árvore do Laravel, o
`kyse` e a dissolução do `porang` acabaram de acontecer. Trocar o motor de dados
por cima disso é o tipo de mudança que impede saber qual das duas quebrou —
a mesma razão pela qual a Fase A manteve o `templ` enquanto a estrutura mudava.

## Consequências

- Query nova continua sendo Go com string SQL e placeholder, no template do
  gerador. Quem escreve query fora do gerador escreve no bloco
  `// arandu:begin custom` (REGRA 15).
- A promessa de `.sql` legível some do material de venda. É perda real: era um
  argumento bom, e sai porque não é verdade.
- O `aru doctor` passa a ser a única barreira contra SQL montado com variável.
  Isso concentra risco numa regra, e é por isso que ela ganhou cobertura de
  mutação — regra de segurança que ninguém prova que dispara é decoração.
- Se o gatilho disparar, a migração é do gerador, não do usuário: o SQL vive nos
  templates, e nenhum projeto existente tem `.sql` para converter.

## Alternativa considerada

**Entrar com o `sqlc` só no perfil de performance.** Rejeitada agora, mantida
como o gatilho: enquanto o perfil não existir, teríamos um caminho de dados no
convencional e outro no de performance, que é precisamente a segunda forma que a
REGRA 9 proíbe. Quando o perfil existir, os dois conjuntos já serão dois, e aí o
`sqlc` reduz formas em vez de somar.
