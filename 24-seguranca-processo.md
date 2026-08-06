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

- endereço dedicado (`security@` do domínio), nunca issue pública
- GitHub Security Advisory como canal alternativo
- chave PGP publicada, opcional para quem quiser usar
- versões suportadas para correção de segurança (doc 23)

Issue pública sobre vulnerabilidade é fechada e movida para advisory privado,
com agradecimento — não com bronca.

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

## Supply chain

- Core com duas dependências (ADR 0004) — essa é a defesa principal
- `govulncheck` na CI de todo repositório, bloqueando merge
- Dependabot com revisão humana obrigatória; nunca merge automático
- `go.sum` verificado; nenhum `replace` apontando para fora em release
- Assets vendorizados (HTMX, Alpine, Bootstrap) com hash conferido no PR de
  atualização

## O `doctor` faz parte da resposta a incidente

Quando uma falha revela um padrão perigoso, a correção tem duas partes: o patch
e **uma regra nova no `doctor`** que impede o padrão de voltar.

Isso é o que transforma incidente em aprendizado permanente — e é a mesma lógica
dos `hints` da página de erro. Regra nova de doctor é breaking change (doc 23),
então entra como aviso na minor e vira erro na major seguinte.
