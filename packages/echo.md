# Echo

**Veredito:** NÃO EXISTE SEPARADO — vive no porang
**Nome Arandu:** —

## O que o Laravel faz

Biblioteca JavaScript cliente que assina canais e escuta eventos do broadcasting.

## Decisão

Não é projeto. Se o `uirapuru` existir, o cliente vive no `porang`.

## Por quê

O Echo existe porque o front do Laravel é JavaScript e precisa de um cliente JS.

No Arandu a interatividade é HTMX, e o cliente de SSE e WebSocket é uma extensão de HTMX mais um atributo no template. Não há biblioteca JS nossa a distribuir — e não ter é justamente o pitch da stack.

Empacotar isso separado criaria uma dependência npm no ecossistema, que é exatamente o que a escolha de Go e HTMX evita.

## O que conteria, se existir

Um componente no `porang` que conecta e re-renderiza o fragmento.

## Gatilho para construir

Junto com o `uirapuru`, se ele existir.

## Alternativa enquanto não existe

Extensão SSE do HTMX.
