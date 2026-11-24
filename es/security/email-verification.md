# Verificación de correo

- [Un código vinculado a un propósito confirma una transición de correo](#a-purpose-bound-code-confirms-one-email-transition)
- [Comparta almacenamiento de un solo uso y configure correo](#share-the-one-time-store-and-configure-mail-delivery)
- [Emita, entregue y consuma mediante el servicio](#issue-deliver-and-consume-through-the-service)
- [Renueve códigos cuando cambie la dirección o solicitud](#rotate-codes-when-the-address-or-request-changes)
- [Los códigos inválidos y expirados revelan el mismo resultado acotado](#invalid-and-expired-codes-reveal-the-same-bounded-result)
- [Verificación no es vinculación de cuenta](#verification-is-not-account-linking)
- [Ejercite emisión, replay, expiración y aislamiento de tenant](#exercise-issuance-replay-expiry-and-tenant-isolation)
- [Complete la verificación antes de conceder la sesión final](#finish-verification-before-granting-the-final-session)

<a name="a-purpose-bound-code-confirms-one-email-transition"></a>
## Un código vinculado a un propósito confirma una transición de correo

Los códigos onetime de Hesape implementan la primitiva publicada y UI :uiVersion proporciona el flujo. El servicio emite un código de seis dígitos vinculado al propósito, el correo entrega su única copia legible y Consume lo acepta atómicamente una vez antes de escribir el estado verificado.

<a name="share-the-one-time-store-and-configure-mail-delivery"></a>
## Comparta almacenamiento de un solo uso y configure correo

Configure `CodeStore` compartido, clave criptográfica, transporte de correo y servicio de usuario consciente del tenant. Separe el propósito de verificación de los de recuperación de contraseña y dos factores, aunque compartan almacenamiento.

- Use propósito y tenant exactos.
- Renderice el código en el mailable de verificación.
- Mantenga códigos sin procesar fuera de logs duraderos.
- Envíe mediante el puente de correo de aplicación.

<a name="issue-deliver-and-consume-through-the-service"></a>
## Emita, entregue y consuma mediante el servicio

Construya el sujeto con tenant fiable, ID de usuario y dirección capturada normalizada. Emita y envíe el código; luego ejecute Consume para el mismo propósito y sujeto antes de que `MarkVerified` escriba condicionalmente la cuenta. Una repetición no reutiliza el código consumido.

```go
const verifyPurpose = "verify-email"
subject := user.AccountTenant() + "\x00" + user.AccountID() + "\x00" + account.NormalizeEmail(user.AccountEmail())

code, err := codes.Issue(ctx, verifyPurpose, subject)
if err != nil { return err }
if err := mailer.Send(ctx, VerifyEmail{User: user, Code: code}); err != nil { return err }

if err := codes.Consume(ctx, verifyPurpose, subject, submitted); err != nil { return err }
_, _, err = users.MarkVerified(ctx, user.AccountTenant(), user.AccountID(), user.AccountEmail())
return err
```

<a name="rotate-codes-when-the-address-or-request-changes"></a>
## Renueve códigos cuando cambie la dirección o solicitud

Reemita tras cambiar la dirección, sustituya códigos anteriores según el contrato del store, limite envíos repetidos y exija nuevo login tras verificar, sin convertir el consumo de código en sesión.

- Vincule propósito, tenant, cuenta y expiración.
- Consuma antes de registrar éxito.
- Envíe una sustitución sin revelar existencia de la cuenta.
- Mantenga los fallos de entrega de correo observables y reintentables.

<a name="invalid-and-expired-codes-reveal-the-same-bounded-result"></a>
## Los códigos inválidos y expirados revelan el mismo resultado acotado

Los códigos desconocidos, alterados, expirados, ya consumidos, de otro propósito y de otro tenant fallan sin revelar qué cuenta existe. Un propósito y un subject pueden presentar veinte códigos erróneos en veinticuatro horas sumando todos los códigos emitidos; después Consume rechaza todo código con `ErrTooManyFailures` hasta que se cierre la ventana, y emitir un código nuevo no reinicia la cuenta. Las pantallas de ui :uiVersion responden a ese rechazo con "too many attempts, try again later" en lugar de llamar inválido al código. Una caída del correo es distinta de un POST inválido; un fallo de base de datos después de Consume necesita un camino de reparación idempotente.

No marque como verificado a partir solo del id de la cuenta. El flujo POST publicado depende de la frontera del almacén de un solo uso, de la finalidad y de la dirección capturada.

<a name="verification-is-not-account-linking"></a>
## Verificación no es vinculación de cuenta

> [!WARNING]
> Nunca vincule cuentas solo por el mismo correo. No registre códigos sin procesar ni los envíe a analytics, URLs o redirecciones. Acepte el código solo en POST protegido por CSRF. La verificación confirma control de una dirección; la policy aún decide pertenencia al tenant y cambios de cuenta.

Use TLS y expiración corta. La sesión de navegador puede llevar estado del flujo, pero el registro duradero de un solo uso impide replay entre procesos.

<a name="exercise-issuance-replay-expiry-and-tenant-isolation"></a>
## Ejercite emisión, replay, expiración y aislamiento de tenant

Pruebe el formulario UI contra Codes de Hesape con códigos válidos, alterados, expirados, consumidos, de propósito/tenant incorrecto, sustitución, fallo de correo y envío entre réplicas. Confirme un único cambio en la cuenta y ninguna exposición de otra cuenta.

- Issue y Consume tienen éxito una vez.
- El replay y la expiración fallan.
- Otro propósito o tenant no puede consumir.
- El consumidor UI guarda estado verificado solo tras Consume.

<a name="finish-verification-before-granting-the-final-session"></a>
## Complete la verificación antes de conceder la sesión final

La autenticación determina cuándo exigir verificación, la recuperación de contraseña usa otro propósito, las sesiones rotan tras transiciones de identidad y el correo transporta el código.

- [Autenticación](/docs/{{version}}/security/authentication)
- [Recuperación de contraseña](/docs/{{version}}/security/password-reset)
- [Sesión](/docs/{{version}}/http/session)
- [Correo](/docs/{{version}}/services/mail)
