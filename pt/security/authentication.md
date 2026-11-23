# Autenticação

- [Autenticação transforma credenciais em Subject](#authentication-turns-credentials-into-a-subject)
- [Conecte armazenamento de usuários, hashing, throttle e sessões](#wire-the-user-store-password-hashing-throttle-and-sessions)
- [Verifique credenciais e rotacione só após o último fator](#verify-credentials-and-rotate-only-after-the-final-factor)
- [Carregue, lembre, confirme e encerre a sessão](#load-remember-confirm-and-sign-out)
- [Retorne uma única falha de credencial](#return-one-credential-failure)
- [Impeça enumeração e fixação de sessão](#keep-enumeration-and-fixation-closed)
- [Teste o caminho completo da requisição](#test-the-complete-request-path)
- [Autorização começa após identidade](#authorization-begins-after-identity)

<a name="authentication-turns-credentials-into-a-subject"></a>
## Autenticação transforma credenciais em Subject

Autenticação identifica quem faz a requisição. Entrada validada fornece credenciais normalizadas ao serviço; o provedor busca conta e verifica hash; a conta válida vira `security.Subject`; o store assina e persiste quem é a pessoa, e cada requisição seguinte relê os papéis pelo resolvedor de sessão. Permissão é decisão separada de policy.

<a name="wire-the-user-store-password-hashing-throttle-and-sessions"></a>
## Conecte armazenamento de usuários, hashing, throttle e sessões

A aplicação precisa de repositório que busque contas no tenant resolvido, chave, serviço de fatores, backend de sessão e throttle de login. Réplicas de produção exigem sessões compartilhadas. O throttle publicado usa memória por processo: cada réplica tem orçamento próprio; adicione controle compartilhado efetivo antes de prometer bloqueio global.

<a name="verify-credentials-and-rotate-only-after-the-final-factor"></a>
## Verifique credenciais e rotacione só após o último fator

Valide o formulário e chame a fronteira de credenciais com tenant confiável, endereço normalizado e senha em texto puro. Se houver segundo fator, grave tentativa pendente assinada de curta duração sem criar Subject autenticado. Só conta apenas com senha ou desafio concluído de autenticador/recuperação rotacionam para a sessão final.

```go
type CredentialVerifier interface {
	Verify(context.Context, string, string, string) (security.Subject, error)
}

subject, err := credentials.Verify(r.Context(), tenant, input.Email, input.Password)
if err != nil { return err }

required, err := factors.Required(r.Context(), subject.Tenant, subject.ID)
if err != nil { return err }
if required {
	return pending.Write(w, subject, input.Remember)
}

oldID := sessions.IDFromRequest(r)
_, err = sessions.Rotate(r.Context(), w, oldID, subject, security.Remember(input.Remember))
return err
```

<a name="load-remember-confirm-and-sign-out"></a>
## Carregue, lembre, confirme e encerre a sessão

`RequireAuth`, `RequireRole` e `RequireConfirmedPassword` carregam a sessão assinada, releem seus papéis pelo resolvedor de sessão e colocam seu Subject no contexto da requisição, onde handlers o leem com `ctx.User`; `middleware.LoadSubject` faz o mesmo numa rota pública e deixa passar um visitante sem Subject. Remember estende a sessão por opção nomeada. Confirmação de senha registra prova recente para rotas sensíveis. Logout invalida a sessão atual; sair de outros dispositivos destrói os outros registros do mesmo sujeito no mesmo tenant.

<a name="return-one-credential-failure"></a>
## Retorne uma única falha de credencial

Endereço desconhecido e senha errada retornam o mesmo erro público. Falhas de busca permanecem operacionais e não gastam orçamento de credenciais. Tentativa bloqueada retorna espera limitada antes da consulta. Cookies inválidos ou expirados viram requisições não autenticadas, não sujeitos parcialmente preenchidos.

<a name="keep-enumeration-and-fixation-closed"></a>
## Impeça enumeração e fixação de sessão

Limite por identidade normalizada e cliente confiável antes de verificar a senha, respeitando a política de tempo das comparações falhas. O throttle oficial em memória divide limites entre réplicas. Nunca registre senha, ID de sessão, payload de fator pendente, código de recuperação ou chave; rotacione só após o último fator e mudança de privilégio.

As telas publicadas contam em `drivers.Limiter` cada envio de um formulário que manda ou verifica um código de uso único, dez por hora por endereço digitado e trinta a cada quinze minutos por cliente, antes de qualquer consulta. O pedido de redefinição e o reenvio de verificação respondem depois do mesmo segundo tenha o endereço uma conta ou não, e a tela de configuração do segundo fator e os códigos de recuperação saem com Cache-Control: no-store, private.

<a name="test-the-complete-request-path"></a>
## Teste o caminho completo da requisição

Exercite login válido, senha errada, conta desconhecida, bloqueio, backend indisponível, remember, rotação, logout e outro tenant com mesmo endereço normalizado. Confirme falhas públicas indistinguíveis enquanto logs e status ainda distinguem indisponibilidade de credencial rejeitada para operadores.

<a name="authorization-begins-after-identity"></a>
## Autorização começa após identidade

Subject não é Grant. Após identificar o ator, o handler chama `security.Authorize` para ação e recurso e passa `security.Grant` à fronteira de dados. Repositórios derivam tenant por `data.Tenant(g)`. CSRF comprova separadamente que a mutação do navegador veio pela sessão vinculada.
