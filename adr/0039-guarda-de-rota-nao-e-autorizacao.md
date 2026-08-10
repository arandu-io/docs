# ADR 0039 — Guarda de rota não é autorização

**Status:** Aceita. Data: 10/08/2026.

## Contexto

`framework/httpx/middleware/` tinha `csrf.go`, `headers.go`, `observability.go`,
`ratelimit.go` e `recover.go`. Não havia equivalente do `auth` nem do `guest` do
Laravel, e `grep -rn 'RequireAuth|RequireGuest|RequireRole'` sobre os quatro
repositórios devolvia zero.

Pior que a ausência: `examples/routes/web.go` e `arandu/routes/web.go`
documentavam `r.Group("/admin", middleware.RequireRole("admin"))` em comentário,
nos **dois repositórios publicados**, exatamente no arquivo que se abre para
aprender a registrar rota. Middleware documentado que não existe é pior que
middleware ausente.

O que isso custava: `/dashboard` respondia 200 para qualquer um — o handler trata
falha de sessão como caso anônimo e renderiza. `GET /auth/login` e
`/auth/register` renderizavam para quem já estava dentro. A única checagem de
sessão do starter kit estava escrita à mão dentro de um handler.

A hesitação era legítima: pela REGRA 2 e pelo `docs/01`, **nenhum repositório é
alcançável sem passar por uma Policy**, e um middleware que decide acesso parece
exatamente a segunda porta que a tese do produto proíbe.

## Decisão

Entram `RequireAuth`, `RedirectIfAuthenticated` e `RequireRole` — e a linha fica
escrita no arquivo, no comentário que alguém lê quando estiver com a tentação de
movê-la:

> Uma guarda responde **"existe sessão"** e **"este sujeito carrega este papel"**.
> Ela não decide **nada** sobre um registro. Se este sujeito pode ler ou escrever
> **esta linha** é resposta da Policy, e a Policy continua rodando — em toda
> escrita e em toda leitura (REGRA 17), sempre por `security.Grant`.

A guarda é conveniência de roteamento: ela evita renderizar uma tela para quem
não tem o que ver nela. Não é camada de autorização, e nada que ela decida
substitui uma Policy.

Isso é o oposto do padrão que se copia sem pensar. Numa aplicação Laravel típica
o middleware `can:` é onde metade da autorização acontece; aqui ele não pode ser,
porque o `Grant` é inforjável por construção e a guarda não emite nenhum.

**O papel é lido da sessão, e a Policy não confia nele.** `RequireRole` olha o
`Subject` que a sessão carrega. Esse mesmo `Subject` chega na Policy, que decide
de novo, contra a linha. Duas leituras do mesmo fato, uma decisão só.

## Consequências

- As duas linhas de comentário mentirosas viraram chamada de verdade.
- O esqueleto `arandu` continua sem `RequireAuth`: ele registra uma rota, `GET /`,
  pública por desenho. `/dashboard` protegido vem do starter kit, que em termos do
  Laravel é o `laravel/ui` e não o `laravel/laravel`. Inventar rota em esqueleto
  publicado é REGRA 3.
- `RequireRole` responde 403, e o HTMX não troca resposta não-2xx — quem cair nela
  numa requisição HTMX não vê nada. O `CSRFProtect` tem a mesma forma e o mesmo
  problema. Ficou como está de propósito: mudar um dos dois cria a inconsistência
  em vez de remover. Se mudar, mudam os dois na mesma vez.
- As telas de recuperação de senha continuam alcançáveis por quem já entrou. O
  Laravel as põe no grupo `guest`; se o Arandu deve pôr é decisão de produto, e
  está registrada aqui como aberta em vez de decidida de passagem.
