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

A `ServerConfig` requires application id and key, Broker, `ConnectPolicy`, `SubscriptionPolicy`, and Protocol. `ServerConfig.Operate`, an `OperationPolicy` asked about `joaju.Operate`, is optional, and nil refuses every HTTP API route with 403; an application that publishes through the API writes one that allows its backend and never an ordinary user. Joaju's own journey and end-to-end tests cover one process with the memory broker. Released relay and RESP pub/sub APIs exist for a fleet, but no tagged consumer proves their application wiring and cross-instance behavior yet; this guide does not present multi-replica delivery as an established deployment path.

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

The in-process consumer must prove that Connect runs before upgrade, Subscribe runs for every route and subscription touching a channel, an HTTP API call without an Operate policy answers 403, and tenant A cannot name tenant B's channel. Fleet support is not yet a supported deployment contract: before relying on it, publish on one instance, assert a socket on another receives the event, and query metrics from both sides.

<a name="degrade-explicitly-when-the-bus-is-unavailable"></a>
## Degrade explicitly when the bus is unavailable

The released relay API can answer metrics with responses that arrive before its timeout, but the first-party application has not established the fleet recovery contract. A deployment that adopts it must alert on bus errors, define whether local service continues, restore connectivity, and re-run cross-instance delivery before considering the fleet recovered.

<a name="protocol-transport-and-application-remain-separate"></a>
## Protocol, transport, and application remain separate

The Protocol owns frame and HTTP route semantics, the `WebSocket` transport owns upgrades and deadlines, the Broker owns channel state, and Relay owns cross-instance propagation. `HesapeBroadcaster` adapts Hesape's broadcasting contract to the same Joaju Broker or Carrier; it does not bypass Joaju connection and subscription policies. Authentication middleware supplies the Subject and auth policies issue Grants. Keeping those boundaries visible is what lets the same server run mounted or standalone without pretending the two deployments have identical identity behavior.
