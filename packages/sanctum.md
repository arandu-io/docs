# Sanctum

**Veredito:** CORE — sem repositório próprio
**Nome Arandu:** —

## O que o Laravel faz

Tokens de API leves e autenticação de SPA por cookie, sem a complexidade de OAuth2.

## Decisão

Vive no core, em `security/token.go` — **decidido, e não construído.**

*Corrigido em 08/08/2026:* esta linha estava no presente, como se o arquivo
existisse. Não existe: `framework/security/` tem `csrf.go`, `password.go`,
`policy.go` e `session.go`, e `find` por `*token*` no projeto inteiro não devolve
nada. O único caminho autenticado hoje é o cookie de sessão.

Foi assim que este pacote atravessou duas fases: veredito CORE, gatilho "não se
aplica" — que é a forma de dizer "é para construir" — e zero linhas. Um veredito
de construir sem código é uma promessa, e o placar de `00-indice.md` a contava
como entregue.

## Por quê

Segurança é a tese do produto. Um mecanismo de autenticação que o usuário precisa lembrar de instalar é um mecanismo que metade dos projetos não vai ter.

Tecnicamente, token de API é uma variação de sessão: emitir, assinar, verificar, revogar. O core já tem o `SessionStore` e o `Grant`. O token reusa os dois.

## O que conteria, se existir

- emissão de token com escopo (lista de Actions permitidas)
- verificação em middleware, produzindo o mesmo `security.Subject` da sessão
- revogação e rotação
- o token **não** amplia poder: ele continua passando pela Policy

## Gatilho para construir

Não se aplica: a decisão é construir, e o que falta é a ordem em que isso entra.

**O gatilho verdadeiro é o `aru new --api`**, que também não existe e é anunciado
no `docs/14`. Os dois são a mesma entrega vista de dois lados: um esqueleto sem
`resources/views/` não tem sessão por cookie para oferecer, e um token sem
esqueleto de API é um mecanismo sem quem o use.

Enquanto nenhum dos dois existir, o registro diz isso — e é decisão do
proprietário construir os dois ou rebaixar ambos para DEFERIDO com gatilho
escrito.

## Alternativa enquanto não existe

—
