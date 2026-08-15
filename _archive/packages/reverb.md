# Reverb

**Veredito:** DEFERIDO
**Nome Arandu:** uirapuru — nome escolhido, repositório **deferido**: só existe
quando o gatilho disparar (o SSE do HTMX não bastar)

## O que o Laravel faz

Servidor WebSocket próprio, compatível com o protocolo Pusher, para broadcasting em tempo real.

## Decisão

Não construir antes de SSE não bastar.

## Por quê

HTMX faz Server-Sent Events por extensão, e SSE cobre a maioria absoluta do que um SaaS chama de tempo real: notificação, atualização de dashboard, progresso de job, badge de mensagem.

SSE é HTTP comum: passa por proxy, não precisa de servidor separado, reconecta sozinho, e funciona com o mesmo middleware de auth e CSRF que já existe.

WebSocket só é necessário com tráfego bidirecional de alta frequência — chat, colaboração ao vivo, jogo. É um caso, não o caso.

## O que conteria, se existir

- servidor WebSocket com autenticação pelo mesmo Subject da sessão
- canais privados autorizados por Policy — canal sem policy é vazamento
- presença

## Gatilho para construir

Cliente com requisito bidirecional real que SSE não atende.

## Alternativa enquanto não existe

SSE via extensão do HTMX, com a camada de view fornecendo o componente.
