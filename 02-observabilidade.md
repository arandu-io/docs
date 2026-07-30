# Pilar: logs e debug

Esta é a razão de existir do framework. As outras partes são commodity — dá para
montar com Chi, sqlc e templ em uma semana. Isso não.

## Por que é aqui que dá para ganhar

No Laravel, a experiência de debug vem de dois pacotes separados: Ignition (a
página de erro) e Telescope (o inspetor). São bons justamente porque sabem coisas
específicas do framework — qual rota, qual middleware, quais queries, e mostram
por default só os frames da aplicação, não os do framework.

No Go, hoje, o equivalente disso não existe. Você tem `slog` na stdlib, que é
excelente e virou o default para projetos novos, e tem OpenTelemetry para
produção. Mas entre os dois falta a camada de dev: o momento em que o request
quebrou e você quer ver, na mesma tela, o stack, as queries com tempo, o que você
dumpou e o que o framework achou que está errado.

## O que o Collector captura

Um `Collector` por request, instalado só em dev ou com header autenticado:

- queries com SQL, args, duração, linhas e **arquivo:linha de quem chamou**
- dumps (`observability.Dump`) com origem e offset no tempo do request
- eventos emitidos
- chamadas HTTP externas com status e duração
- suspeita de N+1: mesmo SQL repetido acima de um limite

Em produção sem header, `FromContext` devolve `nil` e todos os `Record*` são
no-op com receiver nil — custo zero, não "custo baixo".

## A parte que importa: diagnóstico, não dados

Mostrar o stack é o mínimo. O que economiza hora de suporte é a página dizer a
causa provável. Cada padrão reconhecido em `hints()` é uma pergunta que ninguém
vai te fazer duas vezes:

- mesmo SQL repetido → "possível N+1, carregue em lote no repository"
- `grant ausente` → "você chamou repositório sem `security.Authorize`"
- token CSRF inválido → "falta `hx-headers` no template base"

Essa lista cresce com cada bug que você mesmo tomar. É o ativo que não dá para
copiar, porque vem de uso real.

## Regra absoluta

Nada em `observability/errorpage` pode ser alcançável com `Env != dev`. O
middleware `Recover` é o único chamador e checa o flag antes. Cookie e
`Authorization` são redigidos **mesmo em dev**: print de tela de erro vaza sessão
com facilidade absurda.
