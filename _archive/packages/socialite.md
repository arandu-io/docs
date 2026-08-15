# Socialite

**Veredito:** ADAPTER FINO ou nada
**Nome Arandu:** —

## O que o Laravel faz

Login por provedores OAuth: Google, GitHub, Facebook, e mais uma dezena.

## Decisão

Não construir pacote próprio. Se houver demanda, adapter fino sobre `markbates/goth`.

## Por quê

O goth já cobre os mesmos provedores em Go e é mantido. O valor que sobra para nós é pequeno e específico: traduzir o retorno do provedor para o nosso `security.Subject` e emitir a sessão pelo caminho correto.

Isso são cinquenta linhas dentro do starter kit que o `aru make:auth` gera, não um projeto.

## O que conteria, se existir

Uma função no starter kit que converte identidade do goth em Subject e rotaciona a sessão.

## Gatilho para construir

Primeiro cliente pedindo login social.

## Alternativa enquanto não existe

goth, direto.
