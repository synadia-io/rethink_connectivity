# Episode 17: Microservices Architectures in NATS

This directory contains the source code for the various examples outlined in this episode.

## Drawing

https://link.excalidraw.com/l/41qIDWgQhAR/2Y95kp8C1Cq

## Example: Ratings Service

Today we'll walk through what it's like to build a simple 1-5 star ratings service, using nothing but NATS.

### Using the API

## Products API

Add a product:

```bash
jo name="Synadia Cloud" | nats req review.products.synadia-cloud.create
# { id: "synadia-cloud" "name": "Synadia Cloud" }
```

List products

```bash
nats req review.products "" # Maybe add some sorting here?
# [{ id: "synadia-cloud" "name": "Synadia Cloud", count: 445, stars: 4.8 }]
```

Remove a product:

```bash
nats req review.products.synadia-cloud.delete ""
# { id: "synadia-cloud" "name": "Synadia Cloud" }
```

Read all reviews for a product:

```
nats req review.products.synadia-cloud.reviews
```

## Users API

Create a user

```bash
jo name="Jeremy" | nats req review.users.codegangsta.create
```

Delete a user

```bash
nats req review.users.codegangsta.delete ""
```

Add a review

```bash
jo stars=5 text="Had a great time!" | nats req review.users.codegangsta.products.synadia-cloud.create
```

Delete a review

```bash
nats req review.users.codegangsta.products.synadia-cloud.delete ""
```

Read a review

```bash
nats req review.users.codegangsta.products.synadia-cloud
```
