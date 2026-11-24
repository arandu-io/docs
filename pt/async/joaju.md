# Servidor `WebSocket` Joaju

- [Um servidor `WebSocket` protegido pela Policy da aplicação](#a-websocket-server-behind-application-policy)
- [Escolha um processo ou várias instâncias conectadas por relay](#choose-one-process-or-a-relayed-fleet)
- [Execute o processo standalone publicado](#run-the-released-standalone-process)
- [Conecte, assine, publique e observe](#connect-subscribe-publish-and-observe)
- [Limite clientes malformados e tráfego excessivo](#bound-malformed-and-noisy-clients)
- [Todo nome de canal vem de um Grant](#every-channel-name-comes-from-a-grant)
- [Verifique as chamadas às policies antes de afirmar entrega entre instâncias](#verify-policy-calls-before-claiming-fleet-delivery)
- [Defina explicitamente a degradação quando o barramento estiver indisponível](#degrade-explicitly-when-the-bus-is-unavailable)
- [Protocolo, transporte e aplicação continuam separados](#protocol-transport-and-application-remain-separate)

<a name="a-websocket-server-behind-application-policy"></a>
## Um servidor `WebSocket` protegido pela Policy da aplicação

O Joaju implementa o protocolo de fio do Pusher como um `http.Handler`. A aplicação o monta atrás da autenticação de subject e fornece três policies: Connect decide se um subject pode manter um socket, Subscribe decide se esse subject pode alcançar cada canal, e Operate decide quem pode chamar a API HTTP. O próprio servidor não inventa uma identidade. O esqueleto não o monta; monte o handler na aplicação que o serve e teste as três policies ali.

<a name="choose-one-process-or-a-relayed-fleet"></a>
## Escolha um processo ou várias instâncias conectadas por relay

Um `ServerConfig` exige id e chave da aplicação, Broker, `ConnectPolicy`, `SubscriptionPolicy` e Protocol. `ServerConfig.Operate`, uma `OperationPolicy` consultada sobre `joaju.Operate`, é opcional, e nil recusa toda rota da API HTTP com 403; uma aplicação que publica pela API escreve uma que permite o seu backend e nunca um usuário comum. A jornada e os testes ponta a ponta do próprio Joaju cobrem um processo com o broker em memória. As APIs publicadas de relay e de pub/sub RESP existem para uma frota, mas nenhum consumidor publicado prova ainda sua ligação na aplicação e seu comportamento entre instâncias; este guia não apresenta a entrega com várias réplicas como um caminho de implantação estabelecido.

<a name="run-the-released-standalone-process"></a>
## Execute o processo standalone publicado

O comando independente é útil como endpoint do protocolo em um único processo. Forneça os quatro valores que não possuem padrão seguro. O comando autentica o protocolo da aplicação, não pessoas; monte um handler quando as decisões sobre canais privados precisarem usar o Subject autenticado pela própria aplicação.

```bash
JOAJU_APP_ID=app JOAJU_APP_KEY=key JOAJU_APP_SECRET=secret JOAJU_TENANT=acme \
  go run github.com/arandu-io/joaju/cmd/joaju@:joajuVersion
```

<a name="connect-subscribe-publish-and-observe"></a>
## Conecte, assine, publique e observe

O cliente do navegador conecta, autoriza assinaturas privadas ou de presença e assina de novo depois de um backoff de reconexão limitado. As rotas da API publicam eventos e consultam contagens de canais ou usuários depois que a policy Operate permite o chamador; suas rotas POST respondem 415 a um corpo cujo Content-Type não é application/json, antes de lê-lo. Observer recebe eventos de ciclo de vida concluídos sem ganhar o poder de recusá-los, mantendo métricas e logging separados da política de acesso.

<a name="bound-malformed-and-noisy-clients"></a>
## Limite clientes malformados e tráfego excessivo

`NewServer` recusa configuração obrigatória ausente e definições limitadas inválidas. `MaxMessageSize` e `MaxBodySize` limitam a entrada, `OutboundQueue` e `WriteTimeout` limitam consumidores lentos e `MaxMessagesPerSecond` pode medir cada socket depois que o tráfego for medido. `MaxConnections` abaixo de zero seleciona explicitamente uma quantidade ilimitada de conexões; não é uma falha de validação. Uma assinatura recusada nunca deve criar nem expor o canal nomeado.

O comando avulso limita cada socket a 100 frames por segundo, a menos que `JOAJU_MAX_MESSAGES_PER_SECOND` diga outra coisa, com 0 para nenhum limite, e recusa um socket com 429 quando seu endereço de cliente já mantém `JOAJU_MAX_CONNECTIONS_PER_ADDRESS` sockets, 256 por padrão. O endereço é o par TCP, a menos que esse par seja um proxy que `JOAJU_TRUSTED_PROXIES` lista, que é loopback e as faixas privadas quando não definido e ninguém quando definido como none; então X-Forwarded-For é lido da direita. Um proxy num endereço público precisa ser listado, ou todo visitante conta como esse proxy.

<a name="every-channel-name-comes-from-a-grant"></a>
## Todo nome de canal vem de um Grant

`NewConnection` e `NewChannelName` exigem um Grant, do qual deriva o prefixo de tenant. Um nome enviado pelo cliente com o separador é rejeitado, em vez de ser aceito como canal já prefixado. O TLS fica no proxy reverso, e a política de origem do handshake deve corresponder às origens confiáveis da implantação.

<a name="verify-policy-calls-before-claiming-fleet-delivery"></a>
## Verifique as chamadas às policies antes de afirmar entrega entre instâncias

O consumidor em processo precisa provar que Connect roda antes do upgrade, que Subscribe roda para toda rota e assinatura que toca um canal, que uma chamada à API HTTP sem policy Operate responde 403 e que o tenant A não consegue nomear o canal do tenant B. O suporte a frota ainda não é um contrato de implantação suportado: antes de depender dele, publique numa instância, verifique que um socket em outra recebe o evento e consulte métricas dos dois lados.

<a name="degrade-explicitly-when-the-bus-is-unavailable"></a>
## Defina explicitamente a degradação quando o barramento estiver indisponível

A API de relay publicada pode retornar métricas com as respostas recebidas antes do timeout, mas a aplicação de primeira parte ainda não comprovou o contrato de recuperação entre instâncias. Uma implantação que a adote deve alertar sobre erros no barramento, definir se o serviço local continua, restabelecer a conexão e testar novamente a entrega entre instâncias antes de considerar o conjunto recuperado.

<a name="protocol-transport-and-application-remain-separate"></a>
## Protocolo, transporte e aplicação continuam separados

Protocol cuida da semântica de frames e rotas HTTP, o transporte `WebSocket` cuida de upgrades e prazos, Broker cuida do estado de canais e Relay da propagação entre instâncias. `HesapeBroadcaster` adapta o contrato de broadcasting do Hesape ao mesmo Broker ou Carrier do Joaju; não contorna as policies de conexão e assinatura do Joaju. O middleware de autenticação fornece Subject e as policies de auth emitem Grants. Manter essas fronteiras visíveis permite executar o mesmo servidor montado ou isolado sem fingir que as duas implantações têm comportamento de identidade igual.
