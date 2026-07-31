# 17 — Deploy

Você pediu "deploy em um comando" na primeira conversa. A resposta honesta é que
metade disso já está pronta e a outra metade não deve ser construída agora.

## O artefato

Um binário estático. `CGO_ENABLED=0`, assets embutidos por `go:embed`
(`14-view-e-assets.md`), templates compilados pelo templ. Sem runtime, sem
`vendor/`, sem `node_modules`, sem servidor web na frente obrigatório.

```
aru build            # binário + checksum + versão via ldflags
aru build --docker   # imagem distroless, ~25 MB
```

O `aru new` já gera `Dockerfile` e `compose.yml`. Isso cobre Dokploy, Coolify,
Kamal, Fly, Railway, ECS e Kubernetes sem nenhum código nosso — todos consomem
imagem OCI.

## O que já existe no kernel e conta como deploy

- `/_arandu/health` agregando o `Health()` de cada módulo
- shutdown gracioso com drenagem e fechamento na ordem inversa de registro
- config validada no boot: **falha rápido**, não sobe com config inválida
- `X-Request-ID` ecoado, para correlacionar log com incidente

## Migration em deploy: a regra

`aru migrate` **nunca** roda no boot do processo. Com N réplicas subindo juntas,
N migrations concorrem. É passo separado do pipeline, antes do rollout, e é
idempotente.

Corolário: toda migration precisa ser compatível com a versão anterior do binário
rodando ao lado durante o rollout. Coluna nova é nullable ou tem default; remoção
de coluna é duas releases.

## `aru deploy` não vai existir agora

Provisionar servidor é produto — no Laravel isso é o Forge, e o Forge tem um time.
Construir isso agora consome a Fase 2 inteira e compete com Dokploy, que já faz
bem e de graça.

**Gatilho:** depois da Fase 3, com cliente pagando pedindo. Antes disso, entregamos
documentação de deploy (Dokploy, systemd, Kubernetes) — que é documento, não código.
