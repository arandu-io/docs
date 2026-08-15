# Pennant

**Veredito:** CORE — sem repositório próprio
**Nome Arandu:** —

## O que o Laravel faz

Feature flags com escopo por usuário ou time, resolução em memória e persistência opcional.

## Decisão

Vive no core, em `features/`.

## Por quê

São aproximadamente duzentas linhas: um mapa, uma interface de resolução, um middleware e um helper de template.

Repositório separado para isso custa mais em manutenção — release, versionamento, matriz de compatibilidade com o core — do que o código que ele carrega.

## O que conteria, se existir

- `features.Enabled(ctx, "nome")`, resolvendo por Subject do Grant
- resolução em memória com override por env, e persistência opcional via adapter
- integração com a camada de view: componente que não renderiza se a flag está desligada

## Gatilho para construir

Não se aplica — entra no core quando alguém precisar da primeira flag.

## Alternativa enquanto não existe

—
