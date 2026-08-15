# Cashier

**Veredito:** DEFERIDO — Fase 4
**Nome Arandu:** a definir

## O que o Laravel faz

Assinaturas, cobrança recorrente, trials, cupons e webhooks de Stripe e Paddle. Abstrai o ciclo de vida da assinatura sobre o modelo de usuário.

## Decisão

Não construir agora. Reavaliar na Fase 4.

## Por quê

Cobrança é o pacote com maior superfície regulatória e maior custo de manutenção de todo o ecossistema: muda quando o provedor muda a API, quando a legislação fiscal muda, e quebra em produção com dinheiro real envolvido.

Antes da Fase 4 não existe base instalada que justifique carregar isso. E um cliente que precise de cobrança consegue usar o SDK do Stripe direto em três dias.

## O que conteria, se existir

- entidades Subscription, Plan, Invoice como módulo Arandu
- adapters por provedor (Stripe, Polar), atrás de uma interface
- webhooks com verificação de assinatura e idempotência
- policies: quem pode ver e alterar assinatura de quem

## Gatilho para construir

Dois clientes pedindo cobrança, e a decisão de qual provedor sair de você, não deles.

## Alternativa enquanto não existe

SDK oficial do provedor, direto no módulo da aplicação.
