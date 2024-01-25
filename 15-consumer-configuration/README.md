# Episode 15: This ONE feature makes NATS more powerful than Kafka, Pulsar, RabbitMQ and Redis

## Drawing

TODO

## Intro

There's one single feature of NATS JetStream that blows away all of the competition, and it often goes misunderstood or underutilized in most applications. The feature? It's how JetStream designed its **Consumers**.

Today I'm going to walk through a bunch of practical real-world examples of how to leverage JetStream Consumers, and why they are the NATS's secret superpower that allows you to have a flexible, adaptive software architecture without a lot of up-front design costs.

Let's get into it.

### The motivation

I'll be honest, at first I wanted to make a video that simply covered all of the options of JetStream consumers, similar to what I did with streams. Let's just say that would have easily been a multi-hour long workshop explaining all the the different options for consumers. I came to the conclusion that before we can even dive in to what all these options mean, we should talk about why JetStream consumers are all so configurable in the first place. Why all the options? I just want to read from my stream! Here's where things get interesting:

[todo: projection drawing]

After diving in to many other messaging technologies, I'm convinced this right here is what make JetStream so special, consumers are able to access the data on their terms - where you can have a single source of truth for your data in the form of a stream, but how that data is accessed (or consumed), can vary so wildly that entire application lifecycles can be built with a single stream.

I like to describe consumers is that they are like **windows** into the stream, but that's just the tip of the iceberg. All of these use cases I'm showing... Order lifecycle, Reporting, Event sourcing, Lookups, Logging/Monitoring, Simulation/Replay, Counting, Grouping. All these things can be expressed and optimized via NATS consumer model.

Most messaging/streaming systems fall short in this regard. Filtering and delivery is typically inflexible and requires a lot of up front design to really get right. The beauty about a model like this is that I can add one of these things in at any time in my development lifecycle, no having to carefully plan topics or partitions, all of this is fluid and very easy to reason about. So I'm going to walk you through create a Consumer for each of these use cases, but before I do. Let's talk about some critical concepts of consumers.

If you're familiar with Event Sourcing or adjacent concepts like CQRS (Command-Query Responsibility Segregation), then a lot of this is going to feel like familiar territory, and in my opinion JetStream consumers are so incredibly supportive in pulling off these types of architectures without having to write a ton of custom stream processing code.

### Consumer Basics

TODO

#### Push vs Pull consumers

TODO

#### Ephemeral vs Durable Consumers

TODO

#### Message Delivery and Acknowledgements

TODO

### Setting up the Stream

- [ ] Use Benthos to generate data
- [ ] Model an order processing workflow
  - Subject: `orders.{location}.{id}.{status}`
  - States:
    - Pending
    - Validation
    - Confirmed
    - Processing
    - Shipped
    - Delivered
    - Completed
    - Cancelled
    - Returned
    - Failed

### Use cases

- [ ] Order fullfilment: durable. standard ack pending, order preservation not important. No duplicates!! R3
- [ ] Reporting: ephemeral, last hour, batch acks
- [ ] Event Sourcing: Max ack pending of 1, slow but can be partitioned (different episode)
- [ ] Lookup: Ephemeral. Filter on a specific order number. Ack flexible
- [ ] Logging: Headers only. Preserve order
- [ ] Simulation/Replay: Replay Policy. No Ack. Ephemeral. Great for viewing replays of workflows in real time
- [ ] Counting: how many orders were delivered in the last X mins. Start time, headers only, filtered by deliver status. Batch ack
- [ ] Group by: headers only, last by subject, filter, batch ack.
