# Octane

**Veredito:** NÃO EXISTE EM GO
**Nome Arandu:** —

## O que o Laravel faz

Mantém a aplicação Laravel viva em memória entre requests, usando Swoole, RoadRunner ou FrankenPHP.

## Decisão

Não construir. O problema não existe.

## Por quê

**Este é o caso mais claro de todo o mapa.** O Octane existe porque o PHP reinicializa a aplicação inteira a cada request: carrega o autoloader, monta o container, registra os providers, tudo de novo, milhares de vezes por segundo.

Um binário Go boota uma vez e fica em memória. É o comportamento padrão, não um modo especial. Construir um "Octane do Arandu" seria escrever código para ativar algo que já está ativo.

Se aparecer alguém pedindo, a resposta é uma frase: já é assim.

## O que conteria, se existir

Nada.

## Gatilho para construir

Nenhum, em nenhuma circunstância.

## Alternativa enquanto não existe

O comportamento padrão do Go.
