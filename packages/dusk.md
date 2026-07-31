# Dusk

**Veredito:** NÃO CONSTRUIR — biblioteca de terceiro
**Nome Arandu:** —

## O que o Laravel faz

Testes de navegador com API expressiva, gerenciando ChromeDriver por baixo.

## Decisão

Não existe equivalente Arandu. Documentar o uso de chromedp ou rod.

## Por quê

Em Go, chromedp e rod são maduros, mantidos, e falam CDP direto — sem precisar de driver externo, que é metade do problema que o Dusk resolve.

Construir um wrapper nosso adicionaria uma camada de tradução entre o teste e uma biblioteca que já é boa, e nos tornaria responsáveis por acompanhar mudanças do Chrome.

## O que conteria, se existir

Nada. Uma página de guia na Fase D mostrando como testar uma aplicação Arandu com chromedp.

## Gatilho para construir

Nenhum.

## Alternativa enquanto não existe

chromedp ou rod, direto.
