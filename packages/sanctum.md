# Sanctum

**Veredito:** CORE — sem repositório próprio
**Nome Arandu:** —

## O que o Laravel faz

Tokens de API leves e autenticação de SPA por cookie, sem a complexidade de OAuth2.

## Decisão

Vive no core, em `security/token.go`.

## Por quê

Segurança é a tese do produto. Um mecanismo de autenticação que o usuário precisa lembrar de instalar é um mecanismo que metade dos projetos não vai ter.

Tecnicamente, token de API é uma variação de sessão: emitir, assinar, verificar, revogar. O core já tem o `SessionStore` e o `Grant`. O token reusa os dois.

## O que conteria, se existir

- emissão de token com escopo (lista de Actions permitidas)
- verificação em middleware, produzindo o mesmo `security.Subject` da sessão
- revogação e rotação
- o token **não** amplia poder: ele continua passando pela Policy

## Gatilho para construir

Não se aplica.

## Alternativa enquanto não existe

—
