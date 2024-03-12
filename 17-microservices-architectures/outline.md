# Outline

## Title

Rethinking Micro-services: Basics

## Thesis

Introduce the current state of micro-services, how much kit is currently required, and how NATS is well suited to replace many of these technologies in one place. Today I’m going to go over what’s generally required for a micro-services architecture, and what portions of NATS can be used to replace entire pieces of that infrastructure. In future videos in this series, I’m going to dive deep into each of these use cases, and go over practical implementations, pros and cons of each.

## Intro

[TBD: Talk about how we got to where we are, and why the monolith vs microservice debate is happening, and why nobody has truly verbalized the true problem]

## The micro-service stack

- API Gateway: NATS micro
- Load balancing: NATS Queue Groups, NATS Micro
- Discovery: NATS Micro
- Canary management: Subject Mapping
- Logging: Core NATS
- Monitoring & observability: NATS Micro, NATS Tracing
- Authentication/Authorization: NATS JWT Auth
- Configuration management: JetStream KV
- Data store: JetStream Streams, KV, ObjectStore
- Deployment: NEX

## Location Transparency

No longer about "Who do I talk to?", but more "What do I want to talk about?".
