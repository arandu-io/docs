# Servidor `WebSocket` Joaju

- [Un servidor `WebSocket` protegido por la Policy de la aplicación](#a-websocket-server-behind-application-policy)
- [Elige un proceso o varias instancias conectadas mediante relay](#choose-one-process-or-a-relayed-fleet)
- [Ejecuta el proceso standalone publicado](#run-the-released-standalone-process)
- [Conecte, suscriba, publique y observe](#connect-subscribe-publish-and-observe)
- [Limite clientes malformados y tráfico excesivo](#bound-malformed-and-noisy-clients)
- [Todo nombre de canal proviene de un Grant](#every-channel-name-comes-from-a-grant)
- [Verifique las llamadas a las policies antes de afirmar entrega entre instancias](#verify-policy-calls-before-claiming-fleet-delivery)
- [Define explícitamente la degradación cuando el bus no esté disponible](#degrade-explicitly-when-the-bus-is-unavailable)
- [Protocolo, transporte y aplicación siguen separados](#protocol-transport-and-application-remain-separate)

<a name="a-websocket-server-behind-application-policy"></a>
## Un servidor `WebSocket` protegido por la Policy de la aplicación

Joaju implementa el protocolo de cable de Pusher como un `http.Handler`. La aplicación lo monta detrás de la autenticación de subject y proporciona tres policies: Connect decide si un subject puede mantener un socket, Subscribe decide si ese subject puede llegar a cada canal, y Operate decide quién puede llamar a la API HTTP. El propio servidor no inventa una identidad. El esqueleto no lo monta; monte el handler en la aplicación que lo sirve y pruebe allí las tres policies.

<a name="choose-one-process-or-a-relayed-fleet"></a>
## Elige un proceso o varias instancias conectadas mediante relay

Un `ServerConfig` exige id y clave de la aplicación, Broker, `ConnectPolicy`, `SubscriptionPolicy` y Protocol. `ServerConfig.Operate`, una `OperationPolicy` consultada sobre `joaju.Operate`, es opcional, y nil rechaza toda ruta de la API HTTP con 403; una aplicación que publica a través de la API escribe una que permite su backend y nunca a un usuario común. El recorrido y los tests de extremo a extremo del propio Joaju cubren un proceso con el broker en memoria. Una flota añade un `Relay`: constrúyalo con `joaju.NewRelay(ctx, instanceID, bus, logger)` sobre la conexión RESP del módulo `github.com/arandu-io/joaju/redis` (`redis.NewBus`), envuelva el Broker con `joaju.RelayedBroker(base, relay)`, entregue ese mismo valor al Protocol y defina `ServerConfig.Relay`; `NewServer` rechaza un Relay cuyo Broker no fue envuelto. Cada instancia necesita su propio `InstanceID`, como el hostname o el nombre del pod, porque una instancia descarta todo mensaje retransmitido que lleve su id. El comando standalone construye el servidor sin Relay, así que una flota monta el handler en la aplicación. Los tests de Joaju cubren la entrega a las otras instancias exactamente una vez y dentro del tenant, las métricas respondidas entre instancias y la recuperación tras una caída, con un bus en memoria, y el bus RESP contra un servidor real bajo el build tag `integration`.

<a name="run-the-released-standalone-process"></a>
## Ejecuta el proceso standalone publicado

El comando independiente sirve como endpoint del protocolo en un solo proceso. Proporcione los cuatro valores sin un valor predeterminado seguro. El comando autentica el protocolo de la aplicación, no a las personas; monte un handler cuando las decisiones sobre canales privados deban usar el Subject autenticado por la propia aplicación.

```bash
JOAJU_APP_ID=app JOAJU_APP_KEY=key JOAJU_APP_SECRET=secret JOAJU_TENANT=acme \
  go run github.com/arandu-io/joaju/cmd/joaju@:joajuVersion
```

<a name="connect-subscribe-publish-and-observe"></a>
## Conecte, suscriba, publique y observe

El cliente del navegador se conecta, autoriza suscripciones privadas o de presencia y se vuelve a suscribir después de un backoff de reconexión acotado. Las rutas de la API publican eventos y consultan cantidades de canales o usuarios una vez que la policy Operate permite a quien llama; sus rutas POST responden 415 a un cuerpo cuyo Content-Type no es application/json, antes de leerlo. Observer recibe eventos de ciclo de vida completados sin obtener el poder de rechazarlos, manteniendo métricas y logging separados de la política de acceso.

<a name="bound-malformed-and-noisy-clients"></a>
## Limite clientes malformados y tráfico excesivo

`NewServer` rechaza configuración obligatoria ausente y ajustes acotados inválidos. `MaxMessageSize` y `MaxBodySize` limitan la entrada, `OutboundQueue` y `WriteTimeout` limitan consumidores lentos y `MaxMessagesPerSecond` puede medir cada socket después de medir el tráfico. `MaxConnections` menor que cero selecciona explícitamente una cantidad ilimitada de conexiones; no es un fallo de validación. Una suscripción rechazada nunca debe crear ni exponer el canal nombrado.

El comando independiente limita cada socket a 100 frames por segundo, salvo que `JOAJU_MAX_MESSAGES_PER_SECOND` diga otra cosa, con 0 para ningún límite, y rechaza un socket con 429 cuando su dirección de cliente ya mantiene `JOAJU_MAX_CONNECTIONS_PER_ADDRESS` sockets, 256 por defecto. La dirección es el par TCP, salvo que ese par sea un proxy que `JOAJU_TRUSTED_PROXIES` lista, que es loopback y los rangos privados cuando no está definido y nadie cuando vale none; entonces X-Forwarded-For se lee desde la derecha. Un proxy en una dirección pública debe listarse, o cada visitante cuenta como ese proxy.

<a name="every-channel-name-comes-from-a-grant"></a>
## Todo nombre de canal proviene de un Grant

`NewConnection` y `NewChannelName` exigen un Grant, del que deriva el prefijo de tenant. Un nombre enviado por el cliente con el separador se rechaza en vez de aceptarse como un canal ya prefijado. TLS corresponde al proxy inverso y la política de origen del handshake debe coincidir con los orígenes confiables del despliegue.

<a name="verify-policy-calls-before-claiming-fleet-delivery"></a>
## Verifique las llamadas a las policies antes de afirmar entrega entre instancias

El consumidor en proceso debe probar que Connect se ejecuta antes del upgrade, que Subscribe se ejecuta para cada ruta y suscripción que toca un canal, que una llamada a la API HTTP sin policy Operate responde 403 y que el tenant A no puede nombrar el canal del tenant B. Para una flota, añada la verificación entre instancias a los tests de la propia aplicación: publique en una instancia, verifique que un socket en otra recibe el evento y consulte métricas de ambos lados.

<a name="degrade-explicitly-when-the-bus-is-unavailable"></a>
## Define explícitamente la degradación cuando el bus no esté disponible

Cuando Redis no responde, el Relay sigue atendiendo los sockets de su propia instancia: `Join` y `Publish` siguen funcionando, y un evento llega solo a esa instancia hasta que el bus vuelva, porque nada queda en cola para las demás. Escribe una línea de error cuando empieza la caída y una de info cuando termina, reabre las suscripciones cada segundo, y `Relay.Degraded()` responde true hasta que un publish tenga éxito o llegue un mensaje; conéctelo a la comprobación de disponibilidad y alerte sobre la línea de error. Las rutas de métricas esperan `ServerConfig.MetricsTimeout` a las otras instancias, 500 ms por defecto, y responden con los números que llegaron. Tras restaurar Redis, repita la verificación de entrega entre instancias.

<a name="protocol-transport-and-application-remain-separate"></a>
## Protocolo, transporte y aplicación siguen separados

Protocol controla la semántica de frames y rutas HTTP, el transporte `WebSocket` controla upgrades y plazos, Broker controla el estado de canales y Relay la propagación entre instancias. `HesapeBroadcaster` adapta el contrato de broadcasting de Hesape al mismo Broker o Carrier de Joaju; no elude las policies de conexión y suscripción de Joaju. El middleware de autenticación proporciona Subject y las policies de auth emiten Grants. Mantener visibles esos límites permite ejecutar el mismo servidor montado o independiente sin fingir que ambos despliegues tienen el mismo comportamiento de identidad.
