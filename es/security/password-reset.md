# Recuperación de contraseña

- [La recuperación consume una credencial vinculada al propósito](#a-reset-consumes-one-purpose-bound-credential)
- [Prepare almacenamiento de códigos, correo, hashing y sesiones](#prepare-code-storage-mail-hashing-and-sessions)
- [Valide, consuma, guarde y destruya otras sesiones](#validate-consume-write-then-destroy-other-sessions)
- [Evite enumeración en respuestas de solicitud y finalización](#keep-request-and-completion-responses-non-enumerating)
- [No consuma un código antes de validar el candidato](#do-not-consume-a-code-before-the-candidate-is-valid)
- [Recuperar contraseña cambia todas las fronteras activas de identidad](#a-reset-changes-every-active-identity-boundary)
- [Ejecute el recorrido completo publicado de UI](#run-the-complete-released-ui-journey)
- [La recuperación usa códigos únicos, pero no es verificación de correo](#reset-uses-one-time-codes-but-is-not-email-verification)

<a name="a-reset-consumes-one-purpose-bound-credential"></a>
## La recuperación consume una credencial vinculada al propósito

UI :uiVersion implementa la recuperación actual mediante códigos onetime de Hesape. El flujo acepta una solicitud sin enumeración, entrega un código de propósito definido, valida la contraseña nueva, consume el código, guarda el hash y destruye otras sesiones de la cuenta.

<a name="prepare-code-storage-mail-hashing-and-sessions"></a>
## Prepare almacenamiento de códigos, correo, hashing y sesiones

Conecte Codes, las reglas de cuenta con alcance de tenant (`*account.Users`), el hasher, el mailer y el mismo `SessionStore` usado en el inicio de sesión. Construya el origen del restablecimiento a partir de la configuración y mantenga el propósito de restablecimiento separado del de verificación y de la recuperación del segundo factor.

- Use almacenamiento de un solo uso con duración corta.
- Envíe mediante el transporte de correo configurado.
- Valide contraseña y confirmación antes de Consume.
- Comparta invalidación de sesión entre réplicas.

<a name="validate-consume-write-then-destroy-other-sessions"></a>
## Valide, consuma, guarde y destruya otras sesiones

Tras validar, busque la cuenta en el tenant fiable, consuma un código vinculado a la dirección y huella de contraseña actuales y ejecute el reset compare-and-swap de aplicación. El consumidor publicado intenta `DestroyOthers` tras guardar la contraseña y registra el fallo sin deshacer la nueva contraseña.

```go
if err := ctx.Bind(&request); err != nil { return err }
user, err := users.Lookup(ctx, tenant, request.Email)
if err != nil { return err }
fingerprint := account.Fingerprint(user.PasswordHash())
subject := user.AccountTenant() + "\x00" + user.AccountID() + "\x00" +
	account.NormalizeEmail(user.AccountEmail()) + "\x00" + fingerprint
if err := codes.Consume(ctx, "reset-password", subject, request.Code); err != nil { return err }
updated, err := users.ResetPassword(ctx, user.AccountTenant(), user.AccountID(), user.AccountEmail(), fingerprint, request.Password)
if err != nil { return err }
if err := sessions.DestroyOthers(ctx, updated.(*models.User).Subject(), ""); err != nil {
	observability.Log(ctx).Error("signing the account's other sessions out", "error", err)
}
return nil
```

<a name="keep-request-and-completion-responses-non-enumerating"></a>
## Evite enumeración en respuestas de solicitud y finalización

Responda a solicitudes de recuperación igual exista o no la cuenta; limite peticiones mediante política explícita de seguridad, permita que una solicitud nueva sustituya correo antiguo y preserve la sesión recién autenticada solo cuando el flujo publicado la haya elegido deliberadamente.

- Use una respuesta pública para correos conocidos y desconocidos.
- Limite solicitud y consumo por separado.
- Mantenga código y contraseña fuera de datos antiguos del formulario y logs.
- Notifique a la cuenta mediante un mensaje independiente tras el éxito.

<a name="do-not-consume-a-code-before-the-candidate-is-valid"></a>
## No consuma un código antes de validar el candidato

Una contraseña inválida no debe consumir un código utilizable. Los códigos desconocidos, expirados, repetidos, de otro propósito y de otro tenant fallan de la misma forma, y después de veinte códigos erróneos para un propósito y un subject en veinticuatro horas Consume rechaza todo código con `ErrTooManyFailures` hasta que se cierre la ventana, por muchos códigos nuevos que se envíen; las pantallas de ui :uiVersion dicen entonces "too many attempts, try again later". Un fallo de escritura después de Consume exige un camino de recuperación explícito; nunca marque éxito porque el código era válido.

> [!WARNING]
> El flujo publicado valida, ejecuta Consume, compare-and-swap de contraseña y luego `DestroyOthers` como mejor esfuerzo. Un fallo en almacenamiento de sesiones se registra sin deshacer la contraseña; el consumidor actual no puede garantizar una revocación correcta.

<a name="a-reset-changes-every-active-identity-boundary"></a>
## Recuperar contraseña cambia todas las fronteras activas de identidad

> [!WARNING]
> Nunca acepte autoridad de cuenta o tenant del formulario de recuperación, revele existencia de cuenta, registre código o contraseña, preserve otras sesiones por defecto ni use un token permanente de API como credencial de recuperación. Use TLS y orígenes configurados.

La recuperación demuestra control de la credencial de un solo uso entregada; la policy aún controla la cuenta del tenant modificada.

<a name="run-the-complete-released-ui-journey"></a>
## Ejecute el recorrido completo publicado de UI

Pruebe solicitudes conocidas y desconocidas, finalización válida, confirmación inválida, códigos alterados, expirados, repetidos, de propósito incorrecto y otro tenant, consumo concurrente, fallo de correo, fallo de escritura de contraseña y dos réplicas con sesiones compartidas.

- Las respuestas públicas de solicitud no enumeran cuentas.
- La entrada inválida mantiene utilizable el código.
- Consume concurrente tiene éxito una vez.
- Las demás sesiones fallan tras el éxito.

<a name="reset-uses-one-time-codes-but-is-not-email-verification"></a>
## La recuperación usa códigos únicos, pero no es verificación de correo

La verificación de correo usa otro propósito, el hashing protege la nueva contraseña, las sesiones aplican invalidación y el rate limit restringe intentos online.

- [Verificación de correo](/docs/{{version}}/security/email-verification)
- [Hashing](/docs/{{version}}/security/hashing)
- [Sesión](/docs/{{version}}/http/session)
- [Limitación de peticiones](/docs/{{version}}/services/rate-limiting)
