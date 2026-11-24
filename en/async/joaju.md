# Joaju `WebSocket` server

- [A `WebSocket` server behind application policy](#a-websocket-server-behind-application-policy)
- [Choose one process or a relayed fleet](#choose-one-process-or-a-relayed-fleet)
- [Run the released standalone process](#run-the-released-standalone-process)
- [Connect, subscribe, publish, and observe](#connect-subscribe-publish-and-observe)
- [Bound malformed and noisy clients](#bound-malformed-and-noisy-clients)
- [Every channel name comes from a Grant](#every-channel-name-comes-from-a-grant)
- [Verify policy calls before claiming fleet delivery](#verify-policy-calls-before-claiming-fleet-delivery)
- [Degrade explicitly when the bus is unavailable](#degrade-explicitly-when-the-bus-is-unavailable)
- [Protocol, transport, and application remain separate](#protocol-transport-and-application-remain-separate)

<a name="a-websocket-server-behind-application-policy"></a>
## A `WebSocket` server behind application policy

Joaju implements the Pusher wire protocol as an `http.Handler`. The application mounts it behind subject authentication and supplies three policies: Connect decides whether a subject may hold a socket, Subscribe decides whether that subject may reach each channel, and Operate decides who may call the HTTP API. The server itself does not invent an identity. The starter does not mount it, so mount the handler in the application that serves it and test the three policies there.

<a name="choose-one-process-or-a-relayed-fleet"></a>
## Choose one process or a relayed fleet

A `ServerConfig` requires application id and key, Broker, `ConnectPolicy`, `SubscriptionPolicy`, and Protocol. `ServerConfig.Operate`, an `OperationPolicy` asked about `joaju.Operate`, is optional, and nil refuses every HTTP API route with 403; an application that publishes through the API writes one that allows its backend and never an ordinary user. Joaju's own journey and end-to-end tests cover one process with the memory broker. A fleet adds a `Relay`: build it with `joaju.NewRelay(ctx, instanceID, bus, logger)` over the RESP connection of the `github.com/arandu-io/joaju/redis` module (`redis.NewBus`), wrap the Broker with `joaju.RelayedBroker(base, relay)`, give the Protocol that same value, and set `ServerConfig.Relay`; `NewServer` refuses a Relay whose Broker was not wrapped. Each instance needs its own `InstanceID`, such as the hostname or pod name, because an instance drops every relayed message that carries its own id. The standalone command builds its server without a Relay, so a fleet mounts the handler in the application. Joaju's tests cover delivery to the other instances exactly once and inside the tenant, metrics answered across instances, and recovery after an outage, on an in-memory bus, and the RESP bus against a live server under the `integration` build tag.

<a name="run-the-released-standalone-process"></a>
## Run the released standalone process

The standalone command is useful for a single-process protocol endpoint. Supply the four values for which no safe default exists. The command authenticates the application protocol, not people; a mounted handler is the path when private channel decisions must use the application's own authenticated Subject.

```bash
JOAJU_APP_ID=app JOAJU_APP_KEY=key JOAJU_APP_SECRET=secret JOAJU_TENANT=acme \
  go run github.com/arandu-io/joaju/cmd/joaju@:joajuVersion
```

<a name="connect-subscribe-publish-and-observe"></a>
## Connect, subscribe, publish, and observe

The browser client connects, authorizes private or presence subscriptions, and resubscribes after bounded reconnect backoff. API routes publish events and query channel or user counts once the Operate policy allows the caller; their POST routes answer 415 to a body whose Content-Type is not application/json, before reading it. Observer receives completed lifecycle events without gaining the power to refuse them, keeping metrics and logging separate from access policy.

<a name="bound-malformed-and-noisy-clients"></a>
## Bound malformed and noisy clients

`NewServer` refuses missing required config and invalid bounded settings. `MaxMessageSize` and `MaxBodySize` bound input, `OutboundQueue` and `WriteTimeout` bound slow consumers, and `MaxMessagesPerSecond` can meter each socket after traffic is measured. `MaxConnections` below zero explicitly selects an unlimited connection count; it is not a validation failure. A refused subscription must never create or expose the channel it named.

The standalone command meters every socket at 100 frames a second unless `JOAJU_MAX_MESSAGES_PER_SECOND` says otherwise, with 0 for no limit, and refuses a socket with 429 when its client address already holds `JOAJU_MAX_CONNECTIONS_PER_ADDRESS` sockets, 256 by default. The address is the TCP peer unless that peer is a proxy `JOAJU_TRUSTED_PROXIES` lists, which is loopback and the private ranges when unset and nobody when set to none; then X-Forwarded-For is read from the right. A proxy on a public address must be listed, or every visitor counts as that proxy.

<a name="every-channel-name-comes-from-a-grant"></a>
## Every channel name comes from a Grant

`NewConnection` and `NewChannelName` require a Grant, and the tenant prefix is derived from it. A client-supplied name containing the separator is refused instead of accepted as a pre-prefixed channel. TLS belongs at the reverse proxy, and the handshake origin policy must match the origins the deployment intends to trust.

<a name="verify-policy-calls-before-claiming-fleet-delivery"></a>
## Verify policy calls before claiming fleet delivery

The in-process consumer must prove that Connect runs before upgrade, Subscribe runs for every route and subscription touching a channel, an HTTP API call without an Operate policy answers 403, and tenant A cannot name tenant B's channel. For a fleet, add the cross-instance check to the application's own tests: publish on one instance, assert a socket on another receives the event, and query metrics from both sides.

<a name="degrade-explicitly-when-the-bus-is-unavailable"></a>
## Degrade explicitly when the bus is unavailable

When Redis cannot be reached, the Relay keeps serving the sockets its own instance holds: `Join` and `Publish` still succeed, and an event reaches only this instance until the bus answers again, because nothing is queued for the others. It writes one error line when the outage begins and one info line when it ends, reopens its subscriptions every second, and `Relay.Degraded()` answers true until a publish succeeds or a message arrives; wire it to the readiness check and alert on the error line. The metrics routes wait `ServerConfig.MetricsTimeout` for the other instances, 500 ms by default, and answer with the numbers that arrived. After restoring Redis, repeat the cross-instance delivery check.

<a name="protocol-transport-and-application-remain-separate"></a>
## Protocol, transport, and application remain separate

The Protocol owns frame and HTTP route semantics, the `WebSocket` transport owns upgrades and deadlines, the Broker owns channel state, and Relay owns cross-instance propagation. `HesapeBroadcaster` adapts Hesape's broadcasting contract to the same Joaju Broker or Carrier; it does not bypass Joaju connection and subscription policies. Authentication middleware supplies the Subject and auth policies issue Grants. Keeping those boundaries visible is what lets the same server run mounted or standalone without pretending the two deployments have identical identity behavior.
