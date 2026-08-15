# Sail

**Veredito:** NÃO EXISTE EM GO
**Nome Arandu:** —

## O que o Laravel faz

Ambiente Docker de desenvolvimento: PHP, banco, Redis, mail catcher, tudo orquestrado por um CLI.

## Decisão

Não construir. Um `compose.yml` no esqueleto resolve.

## Por quê

O Sail existe porque montar PHP local dói: versão do php-fpm, extensões, nginx, permissões. É um problema de instalação da linguagem.

Em Go você roda `go run ./cmd/app`. O que ainda precisa de container é Postgres e Redis, e isso é um arquivo de vinte linhas dentro de `arandu-io/arandu`, não um projeto.

## O que conteria, se existir

Um `compose.yml` no esqueleto, com Postgres e Redis. Nada mais.

## Gatilho para construir

Nenhum.

## Alternativa enquanto não existe

`go run` mais `docker compose up -d postgres redis`.
