
# Overview

Materializations are processed as cooperative transactions between the Flow Runtime
and a connector Driver, over a long-lived RPC through which the Runtime and Driver
exchange messages.

This RPC workflow maintains a materialized view of a Flow collection
in an external system. It has distinct load, prepare, store, and commit
phases. The Flow runtime and driver cooperatively maintain a fully-reduced
view of each document by loading current states from the store, reducing in
a number of updates, and then storing updated documents and checkpoints.

## Transactions

The beating heart of transactions in materializations is this:
there is a consumption checkpoint, and there is a state of the view.
As the materialization progresses, both the checkpoint and the view state
will change. Updates to the checkpoint and to the view state MUST always
commit together, in the exact same transaction.

Flow transaction tasks have a backing transactional recovery log,
which is capable of durable commits that update both the checkpoint
and also a (reasonably small) driver-defined state. More on driver
states later.

Many interesting systems are also fully transactional in nature.

When implementing a materialization driver, the first question an
implementor must answer is: whose commit is authoritative?
Flow's recovery log, or the materialized system?
This protocol supports either.

# Common Materialization Patterns

There are a few common implementation patterns for materializations.
The choice of pattern depends on the transaction capabilities of the remote endpoint.

## Remote Store is Authoritative

In this pattern, the remote store (for example, a database) persists view states
and the Flow consumption checkpoints which those views reflect.
There are many such checkpoints: one per task split, and in this pattern
the Flow recovery log is effectively ignored.

Typically this workflow runs in the context of a synchronous `BEGIN/COMMIT`
transaction, which updates table states and a Flow checkpoint together.
The transaction need be scoped only to the store phase of this workflow,
as the materialization protocol requires only read-committed isolation
semantics.

Flow is a distributed system, and an important consideration is the effect
of a "zombie" assignment of a materialization task, which can race a
newly-promoted assignment of that same task.

Fencing is a technique which uses the transactional capabilities of a store
to "fence off" an older zombie assignment, such that it's prevented from
committing further transactions. This avoids a failure mode where:
- New assignment N recovers a checkpoint at Ti.
- Zombie assignment Z commits another transaction at Ti+1.
- N beings processing from Ti, inadvertently duplicating the effects of Ti+1.

When a remote store is authoritative, it must implement fencing behavior.
As a sketch, the store can maintain a nonce value alongside the checkpoint
of each task split. The nonce is updated on each open of this RPC,
and each commit transaction then verifies that the nonce has not been
changed.

In the future, if another RPC opens and updates the nonce, it fences off
this instance of the task split and prevents it from committing further
transactions.

## Recovery Log with Non-Transactional Store

In this pattern, the runtime's recovery log persists the Flow checkpoint
and handles fencing semantics. During the Load and Store phases, the driver
directly manipulates a non-transactional store or API, such as a key/value store.

Note that this pattern is at-least-once. A transaction may fail part-way
through and be restarted, causing its effects to be partially or fully
replayed.

Care must be taken if the collection's schema has reduction annotations
such as `sum`, as those reductions may be applied more than once due to
a partially completed, but ultimately failed transaction.

If the collection's schema is last-write-wins, this mode still provides
effectively-once behavior. Collections which aren't last-write-wins
can be turned into last-write-wins through the use of derivations.

## Recovery Log with Idempotent Apply

In this pattern the recovery log is authoritative, but the driver uses
external stable storage to stage the effects of a transaction -- rather
than directly applying them to the store -- such that those effects can be
idempotently applied after the transaction commits.

This allows stores which feature a weaker transactionality guarantee to
still be used in an exactly-once way, so long as they support an idempotent
apply operation.

Driver checkpoints can facilitate this pattern. For example, a driver might
generate a unique filename in S3 and reference it in its prepared
checkpoint, which is committed to the recovery log. During the "store"
phase, it writes to this S3 file. After the transaction commits, it tells
the store of the new file to incorporate. The store must handle
idempotency, by applying the effects of the unique file just once, even if
told of the file multiple times.

A related extension of this pattern is for the driver to embed a Flow
checkpoint into its driver checkpoint. Doing so allows the driver to
express an intention to restart from an older alternative checkpoint, as
compared to the most recent committed checkpoint of the recovery log.

As mentioned above, it's crucial that store states and checkpoints commit
together. While seemingly bending that rule, this pattern is consistent
with it because, on commit, the semantic contents of the store include BOTH
its base state, as well as the staged idempotent update. The store just may
not know it yet, but eventually it must because of the retried idempotent
apply.

Note the driver must therefore ensure that staged updates are fully applied
before returning `Loaded` responses, in order to provide the correct
read-committed semantics required by the Flow runtime.

## Push-only Endpoints & Delta Updates

Some systems, such as APIs, Webhooks, and Pub/Sub, are push-only in nature.
Flow materializations can run in a "delta updates" mode, where loads
are always skipped and Flow does not attempt to store fully-reduced
documents. Instead, during the store phase, the runtime sends delta
updates which reflect the combined roll-up of collection documents
processed only within this transaction.

To illustrate the meaning of a delta update, consider documents which are
simple counters, having a collection schema that uses a `sum` reduction
strategy.

Without delta updates, Flow would reduce documents -1, 3, and 2 by `sum`
to arrive at document 4, which is stored. The next transaction,
document 4 is loaded and reduced with 6, -7, and -1 to arrive at a new
stored document 2. This document, 2, represents the full reduction of the
collection documents materialized thus far.

Compare to delta updates mode: collection documents -1, 3, and 2 are
combined to store a delta-update document of 4. The next transaction starts
anew, and 6, -7, and -1 combine to arrive at a delta-update document of -2.
These delta updates are a windowed combine over documents seen in the
current transaction only, and unlike before are not a full reduction of the
document. If delta updates were written to pub/sub, note that a subscriber
could further reduce over each delta update to recover the fully reduced
document of 2.

Note that many use cases require only `lastWriteWins` reduction behavior,
and for these use cases delta updates does the "right thing" by trivially
re-writing each document with its most recent version. This matches the
behavior of Kafka Connect, for example.

# Transaction Protocol

As a convention and to reduce ambiguity,
message types from the Runtime are named in an imperative fashion (`Load`),
while responses from the driver always have a past-tense name (`Loaded`):

```sequence
Runtime->Driver: Open{MaterializationSpec, driverCP}
Note right of Driver: Connect to endpoint.\nOptionally fetch last-committed\nruntime checkpoint.
Driver->Runtime: Opened{runtimeCP}

Note over Runtime, Driver: One-time initialization ☝️.\n 👇 Repeats for each transaction.

Note left of Runtime: Prior txn commits\nto recovery log.
Note right of Driver: Prior txn commits to DB\n(where applicable).
Runtime->Driver: Acknowledge
Note right of Runtime: Acknowledged MAY be sent\nbefore Acknowledge.
Note right of Driver: MAY perform an idempotent\napply of last txn.
Note left of Runtime: Runtime does NOT await\nAcknowledged before\nproceeding to send Load.
Driver->Runtime: Acknowledged
Note left of Runtime: Runtime may now finalize\na pipelined transaction.

Note over Runtime, Driver: End of Acknowledge phase.

Runtime->Driver: Load<A>
Note left of Runtime: Load keys may\n not exist (yet).
Runtime->Driver: Load<B>
Note right of Driver: MAY evaluate Load immediately,\nor stage for deferred retrieval.
Driver->Runtime: Loaded<A>
Runtime->Driver: Load<C>
Runtime->Driver: Flush
Driver->Runtime: Loaded<C>
Note right of Driver: Omits Loaded for keys\nthat don't exist.
Driver->Runtime: Flushed
Note left of Runtime: All existing keys\nhave been retrieved.

Note over Runtime, Driver: End of Load phase.

Runtime->Driver: Store<X>
Runtime->Driver: Store<Y>
Runtime->Driver: Store<Z>

Runtime->Driver: StartCommit{runtimeCP}
Note right of Driver: * Completes all Store processing.\n* MAY include runtimeCP in DB txn.
Note right of Driver: Commit to DB\nnow underway.
Driver->Runtime: StartedCommit{driverCP}
Note left of Runtime: Begins commit to\n recovery log.

Note over Runtime, Driver: End of Store phase. Loops around\n to Acknowledge <=> Acknowledged.
```

## Acknowledge Phase

`Acknowledge` and `Acknowledged` are always the first messages sent every transaction,
including the very first transaction of an RPC. The Runtime sends `Acknowledge` to
indicate that the last transaction has committed to the recovery log.
The Driver sends `Acknowledged` to indicate that its endpoint transaction has committed.

Acknowledge and Acknowledged _are not ordered_.
Acknowledged may be sent before Acknowledge and vice versa.

The Runtime *does not* wait for `Acknowledged` before sending `Load` messages.
In most cases the Driver should simply not read these `Load` messages
until it has completed its own commit and sent its own `Acknowledged`.

A Driver MAY instead process its commit and acknowledgment in the background while
actively reading `Load` messages. It MUST NOT evaluate `Load`s yet, as this could
otherwise be a violation of read-committed semantics, but it MAY stage them for
deferred evaluation. This is **recommended** for Drivers that have very long
commit and/or acknowledgement operations. While a background commit progresses
the Flow runtime will optimistically pipeline the next transaction,
processing documents and preparing for when the Driver sends `Acknowledged`.

Drivers following the "Recovery Log with Idempotent Apply" pattern must take care
to properly handle the very first acknowledgement phase of an RPC. At startup,
a driver cannot know if the last commit has been acknowledged. For example,
a previous RPC invocation may have failed immediately after commit but prior to
acknowledgement. The Driver must thus idempotent-ly apply or re-apply changes staged
by a prior Driver invocation, and reply with `Acknowledged` only once done.

Drivers with transactional semantics SHOULD send Acknowledged
immediately after a previous, started commit completes.

Drivers with at-least-once semantics SHOULD send Acknowledged
immediately after sending StartedCommit.

**What's Chaning**: currently, Acknowledge and Acknowledged are intermixed among Load and Loaded requests. This requires that the driver handle an Acknowledge / Acknowledged flow while in the middle of a stream of Load & Loaded request. Instead, we will synchronously block for the recovery log commit -- a few milliseconds -- in exchange for a substantially simpler protocol:
the driver is now guaranteed that Acknowledge follows StartCommit, and the runtime has a similar guarantee.
We *do* still allow for pipelined loads, if the driver chooses, but the driver is also free to simply block reads (which was not previously possible, because of the need to read Acknowledge before Acknowledged).
In practice only the Snowflake and BigQuery drivers should do this.
The implementations should be _a lot_ easier because we do not need to send into the RPC
until `Flush` is received, which will _never_ be sent until `Acknowledged` is sent,
so (absent a Runtime bug) the protocol provides natural & easy synchronization points.

## Load phase

Zero or more `Load` messages are sent by the Runtime with documents to fetch.
A given document key will appear at most once in a transaction,
and will not be repeated across `Load` messages.

Drivers may immediately evaluate each `Load` and respond,
or may queue many keys to load and defer their evaluation.
The Runtime does not await any individual `Load` requests.

After the previous transaction has fully completed,
and the driver has sent `Acknowledged` to the Runtime,
the current transaction may begin to close.

The Runtime indicates this by sending a `Flush` message, which is NEVER sent
before `Acknowledged` is received. `Acknowledged` is thus an important signal
as to when the Runtime may begin to finalize an optimistic, pipelined transaction.

On reading `Flush`, Drivers must process all remaining `Load` messages,
including any deferred evaluations, and send all `Loaded` responses prior to
sending its own `Flushed` response.

This signals to the Runtime that all documents which can be loaded *have* been loaded,
and the transaction proceeds to the Store phase.

Materialization bindings which are processing in delta-updates mode
will never receive a `Load` message, but will receive a `Flush`
and must still respond with `Flushed`.

**What's new**: the current `Prepare` / `Prepared` messages are replaced with `Flush`/`Flushed`.
`Flush` and `Flushed` _no longer_ carry Flow and driver checkpoints (respectively).

## Store phase

Zero or more `Store` messages are sent by the Runtime to the Driver,
indicating keys, documents, and extracted fields to store.
No response is required of the Driver for these messages.

Once all documents have been stored, the Runtime sends a
`StartCommit` message which carries its opaque runtime checkpoint.

Drivers implementing the "Remote Store is Authoritative" pattern must
include the runtime checkpoint in its current transaction,
for retrieval in a future Open of a new transactions RPC.
Other driver patterns MAY ignore this checkpoint.

On reading `StartCommit` the driver ensures that all `Store` messages have been processed.
It begins to commit its own transaction (where applicable), and then responds with `StartedCommit`
which contain an update to the driver's checkpoint.

On the Runtime's receipt of `StartedCommit`,
the Runtime now knows that all `Store` messages have been fully processed.
It preserves the updated Driver checkpoint in its recovery log
and begins to commit.

From here, the protocol loops back around to the Acknowledge phase.

**What's new**: the current `Commit` / `DriverCommitted` messages are replaced with `StartCommit`/`StartedCommit`.
Flow and driver checkpoints are added to these messages.
Furthermore, `StartedCommit` is returned immediately by the Driver at the start of
its commit operation, unlike `DriverCommitted` which awaiting the driver's full commit.

## Other Migration Notes

We should remove RunTransactions. It's not serving us, and really
obfuscates how these different implementation patterns actually work.

Instead, we should have pure function helpers (similar to the current lifecycle.go implementations) which make it easy to read expected messages from the RPC stream,
and individual connectors compose these to implement their own top-level transaction loops.

I think this makes it a lot easier to express synchronous vs async semantics,
since synchronization can happen explicitly within the scope of a single loop.
Further, the synchronous and most-common case is a simple loop with no concurrency.

Pipelined cases require a single go goroutine to await the commit and then send
Acknowledged. No other coordination is required (!) since iterating over `Load`
cannot cannot EOF until Acknowledged is sent -- meaning it will not attempt to
transmit into the RPC until that time.