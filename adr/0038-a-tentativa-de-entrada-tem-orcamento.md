# ADR 0038 — A tentativa de entrada tem orçamento, e ele é cobrado antes da senha

**Status:** Aceita. Data: 10/08/2026.

## Contexto

Não existia contador de tentativas de login em lugar nenhum:
`grep -rniE 'throttle|lockout|maxAttempts'` sobre `framework/` e `ui/` não
achava nada em auth. O `laravel/ui` tem o trait `ThrottlesLogins` — 5 tentativas,
janela de 1 minuto, chave por identidade + endereço, zerada no acerto — e o
mapeamento contra ele registrou a ausência como lacuna de segurança, não de
paridade.

O limitador global que existe não substitui, e o próprio arquivo diz por quê:
`framework/httpx/middleware/ratelimit.go` já observa que limitar login por IP
sozinho não impede credential stuffing distribuído. O que estava montado nos dois
bootstraps é 300/min compartilhado entre **todas** as rotas, com janela fixa por
processo — então o orçamento multiplica pelo número de réplicas.

## Decisão

Quatro portas fecham juntas.

**O orçamento é cobrado no momento da decisão, não depois da falha.** A primeira
implementação era check-then-act — `Blocked()` perguntava, `Fail()` registrava, e
entre os dois havia um argon2 de 50 a 100 ms em que nada tinha sido anotado.
Medido com oito chamadas simultâneas contra um orçamento de cinco: as oito senhas
foram verificadas. "Cinco tentativas por minuto" valia contra um atacante educado
o bastante para esperar cada resposta. A API é `Attempt` / `Refund` / `Clear`: a
leitura dos dois orçamentos e a escrita dos dois contadores acontecem na mesma
seção crítica. Tentativa recusada não custa nada, então martelar uma conta já
travada não estende a trava.

**O controle mora sob o `Authenticate`, não no handler.** O contrato do kit diz
que o arquivo publicado pode ser apagado pelo projeto (ADR 0026); controle de
segurança em handler editável é controle que some no dia em que alguém reescreve
a tela. O `Authenticate` ganhou um parâmetro `client`, e é o único lugar por onde
a verificação de senha passa.

**`SignInThrottle` é interface própria, não um alargamento de
`middleware.Limiter`.** Não é preferência: `httpx` importa `security`, então
`security` nomear `middleware.Limiter` é ciclo de import.

**A janela é fixa, e a ADR assume o custo.** Numa virada de janela cabem até dez
tentativas em vez de cinco. Janela deslizante seria uma segunda noção de "janela"
no projeto (REGRA 9) para transformar dez em cinco num caso de borda.

Dois orçamentos, e o segundo existe por um motivo específico: por identidade
(5) impede adivinhar a senha de uma conta; por endereço (25) impede um endereço
percorrer uma lista de contas. O `Clear` devolve **uma** unidade ao endereço, não
a contagem inteira — senão um script que possui uma conta no sistema comprava
uma lista nova a cada acerto.

`Refund` devolve exatamente a unidade que o `Attempt` tomou, e é chamado no ramo
de erro interno. Ele nunca faz `Clear`: aquele ramo é alcançável por quem
desligar no meio da requisição, e um `Clear` ali faria de "quatro palpites e uma
requisição abandonada" um número ilimitado de palpites.

## Consequências

- O `KeyByIP` passou a mascarar IPv6 para o /64. Devolvia o endereço inteiro, e
  um /64 é o menor bloco que qualquer ponta recebe — 2^64 endereços que chegam
  todos ao servidor. Uma máquina com um /64 roteado tinha orçamento por cliente
  ilimitado, podia encher a tabela do throttle sozinha, e furava também o
  limitador global. IPv4 não mudou um byte.
- O `Attempt` é consultado **antes** da tabela de usuários, então uma entrada
  travada roda zero statements e a resposta é idêntica para endereço cadastrado e
  desconhecido — a propriedade anti-enumeração continua valendo.
- Tabela limitada, com despejo que perdoa em vez de recusar. Um milhão de
  endereços digitados no formulário não é um milhão de contadores.
- Nada é configurável. Um projeto que precise de outro número tem que dizer
  por quê, e o pedido é o gatilho para esta ADR mudar.

## Alternativa descartada

Cobrar a unidade **depois** da falha, que é o que o Laravel faz. É a forma
intuitiva e é a que abre a janela do parágrafo acima: entre perguntar e anotar
existe o custo da própria verificação de senha, que é grande de propósito.
