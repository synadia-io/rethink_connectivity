# Episode 7: Leaf Nodes

```sh
# show ngs credentials
nats account info
nsc describe account

# Configure leaf node server
vim leaf.conf

# Run leaf node server
nats-server -c leaf.conf

# Try some communication
nats req ngs.echo "Hello"
nats req ngs.echo "Hello" --context leaf

# Put a responder closer to me
nats reply ngs.echo --echo --context leaf
nats req ngs.echo "Hello" --context leaf

# Let's look at a stream I have on my ngs account
nats s ls
nats s subjects

# Viewing a stream from the leaf node
nats s view --context leaf
nats s view --context leaf --js-domain ngs

# Getting
nats s get orders -S "orders.store_1.*" --context leaf --js-domain ngs

# Sourcing a stream
nats s add --context leaf --source orders

nats s info orders --context leaf

# Getting
nats s get orders -S "orders.store_1.*" --context leaf --js-domain ngs
nats s get orders -S "orders.store_1.*" --context leaf

```
