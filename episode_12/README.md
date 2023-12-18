# NATS JetStream: WorkQueue Stream

This example shows how you can build a worker consumer to interact with a work queue stream

## Creating a stream

You'll need a stream to get this working, simply create one with the NATS cli:

```sh
nats s create jobs
[NGS-cli] ? Subjects jobs.*.*
[NGS-cli] ? Storage file
[NGS-cli] ? Replication 1
[NGS-cli] ? Retention Policy Work Queue
[NGS-cli] ? Discard Policy New
[NGS-cli] ? Stream Messages Limit 20000
[NGS-cli] ? Per Subject Messages Limit -1
[NGS-cli] ? Total Stream Size 256MB
[NGS-cli] ? Message TTL 7d
[NGS-cli] ? Max Message Size -1
[NGS-cli] ? Duplicate tracking time window 2m0s
[NGS-cli] ? Allow message Roll-ups No
[NGS-cli] ? Allow message deletion Yes
[NGS-cli] ? Allow purging subjects or the entire stream Yes
Stream jobs was created
```

## Publishing to the stream

Then you can publish high or low priority jobs to the stream:

```sh
nats bench --pub 10 jobs.high.my_job_id --syncpub --msgs 5000
nats bench --pub 10 jobs.low.my_job_id --syncpub --msgs 5000
```

## Running the workers

Once messages are in the stream, you can fire up a worker. This go program takes 2 arguments, one for the priority and one for the id of the worker (for logging purposes):

```sh
# Run 10 high priority workers
seq 10 | xargs -P10 -I {} go run worker.go high {}

# Run 2 low priority workers
seq 2 | xargs -P2 -I {} go run worker.go low {}
```

## Adding a DLQ

Messages that fail to be delivered and reach max attempts will stop being delivered, and NATS will emit an advisory:

`$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.jobs.{consumer_name}`

This advisory can be used to keep track of jobs that have not completed processing for whatever reason. To keep track of this, you can create a stream:

```sh
nats s create jobs_dlq --subjects '$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.jobs.*'
```

The payload information on the events should give you everything you need to process. For tooling purposes, we can also use the new stream metadata feature to mark this stream as a dead letter queue, so tooling can pick up on the semantics:

```sh
nats s edit jobs_dlq --metadata="dead_letter_queue=true"
```
