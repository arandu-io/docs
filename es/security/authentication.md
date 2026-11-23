# Autenticación

- [La autenticación transforma credenciales en Subject](#authentication-turns-credentials-into-a-subject)
- [Conecte almacenamiento de usuarios, hashing, throttle y sesiones](#wire-the-user-store-password-hashing-throttle-and-sessions)
- [Verifique credenciales y rote solo tras el último factor](#verify-credentials-and-rotate-only-after-the-final-factor)
- [Cargue, recuerde, confirme y cierre la sesión](#load-remember-confirm-and-sign-out)
- [Devuelva un único fallo de credencial](#return-one-credential-failure)
- [Impida enumeración y fijación de sesión](#keep-enumeration-and-fixation-closed)
- [Pruebe el recorrido completo de petición](#test-the-complete-request-path)
- [La autorización empieza después de identidad](#authorization-begins-after-identity)

<a name="authentication-turns-credentials-into-a-subject"></a>
## La autenticación transforma credenciales en Subject

La autenticación identifica quién hace la petición. La entrada validada proporciona credenciales normalizadas al servicio; el proveedor busca la cuenta y verifica el hash; la cuenta válida se vuelve `security.Subject`; el store firma y persiste quién es la persona, y cada petición posterior vuelve a leer los roles mediante el resolvedor de sesión. El permiso es una decisión separada de policy.

<a name="wire-the-user-store-password-hashing-throttle-and-sessions"></a>
## Conecte almacenamiento de usuarios, hashing, throttle y sesiones

La aplicación necesita un repositorio que busque cuentas en el tenant resuelto, clave, servicio de factores, backend de sesión y throttle de login. Las réplicas de producción exigen sesiones compartidas. El throttle publicado usa memoria por proceso: cada réplica tiene presupuesto propio; añada control compartido efectivo antes de prometer bloqueo global.

<a name="verify-credentials-and-rotate-only-after-the-final-factor"></a>
## Verifique credenciales y rote solo tras el último factor

Valide el formulario y llame a la frontera de credenciales con tenant fiable, dirección normalizada y contraseña en texto plano. Si hay segundo factor, guarde un intento pendiente firmado de corta duración sin crear Subject autenticado. Solo una cuenta de contraseña única o un desafío completado de autenticador/recuperación rotan hacia la sesión final.

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
## Cargue, recuerde, confirme y cierre la sesión

`RequireAuth`, `RequireRole` y `RequireConfirmedPassword` cargan la sesión firmada, vuelven a leer sus roles mediante el resolvedor de sesión y colocan su Subject en el contexto de la solicitud, donde los handlers lo leen con `ctx.User`; `middleware.LoadSubject` hace lo mismo en una ruta pública y deja pasar a un visitante sin Subject. Remember amplía la sesión mediante opción con nombre. La confirmación de contraseña registra prueba reciente para rutas sensibles. Logout invalida la sesión actual; cerrar otros dispositivos destruye los demás registros del mismo sujeto en el mismo tenant.

<a name="return-one-credential-failure"></a>
## Devuelva un único fallo de credencial

Una dirección desconocida y contraseña incorrecta devuelven el mismo error público. Los fallos de búsqueda siguen siendo operativos y no consumen presupuesto de credenciales. Un intento bloqueado devuelve espera acotada antes de consultar. Las cookies inválidas o expiradas se vuelven peticiones no autenticadas, no sujetos parcialmente rellenados.

<a name="keep-enumeration-and-fixation-closed"></a>
## Impida enumeración y fijación de sesión

Limite por identidad normalizada y cliente fiable antes de verificar contraseña, respetando la política temporal de comparaciones fallidas. El throttle oficial en memoria divide límites entre réplicas. Nunca registre contraseña, ID de sesión, payload de factor pendiente, código de recuperación ni clave; rote solo tras el último factor y cambios de privilegio.

Las pantallas publicadas cuentan en `drivers.Limiter` cada envío de un formulario que manda o comprueba un código de un solo uso, diez por hora por dirección escrita y treinta cada quince minutos por cliente, antes de cualquier consulta. La solicitud de restablecimiento y el reenvío de verificación responden después del mismo segundo tenga o no cuenta la dirección, y la pantalla de configuración del segundo factor y los códigos de recuperación salen con Cache-Control: no-store, private.

<a name="test-the-complete-request-path"></a>
## Pruebe el recorrido completo de petición

Ejercite login válido, contraseña incorrecta, cuenta desconocida, bloqueo, backend indisponible, remember, rotación, logout y otro tenant con la misma dirección normalizada. Confirme fallos públicos indistinguibles mientras logs y estados aún distinguen indisponibilidad de credencial rechazada para operadores.

<a name="authorization-begins-after-identity"></a>
## La autorización empieza después de identidad

Subject no es Grant. Tras identificar al actor, el handler llama a `security.Authorize` para acción y recurso y pasa `security.Grant` a la frontera de datos. Los repositorios derivan tenant mediante `data.Tenant(g)`. CSRF demuestra por separado que la mutación del navegador llegó mediante la sesión vinculada.
