import { createEffect, createSignal, onCleanup, onMount } from "solid-js";
import Sidebar from "./sidebar";
import ChannelView from "./channel-view"
import { JSONCodec, StringCodec, connect, millis, type Consumer, type ConsumerMessages, type JsMsg, type NatsConnection } from "nats.ws";
import { createStore } from "solid-js/store";
import type { Message, Channel, UserID, User } from "../types";

// represents the overall state of our
// chat application, while also allowing for
// efficient updates/appends, which is what
// we need for things like messages
interface ChatStore {
  // NATS related fields
  conn?: NatsConnection
  consumer?: Consumer

  // Represents the ID of the user for publishing messages
  // to various channels. For security, we will want to lock
  // down the subjects that this user is able to publish to
  userID?: string

  // Messages for various channels, for now we will just get
  // all messages from the beginning of time, but with NATS
  // it's quite easy to fetch from a particular time 
  // (7 days back), for instance
  messages: Record<Channel, Message[]>

  // Lookup table of user in this workspace. These user profiles will be supplied
  // by a NATS KV store
  users: Record<UserID, User>
}

const sc = StringCodec()

const channels = ["general", "random", "dev"]

export default function Chat() {
  const [selected, setSelected] = createSignal("general")

  const [store, setStore] = createStore<ChatStore>({
    userID: "amVyZW15QHN5bmFkaWEuY29t",
    messages: {},
    users: {}
  })

  const onMessageReceived = (m: JsMsg) => {
    const [_, channel, userID] = m.subject.split(".")

    const msg: Message = {
      userID: userID,
      text: m.string(),
      timestamp: new Date(millis(m.info.timestampNanos))
    }
    setStore("messages", channel, (prev) => prev ? [...prev, msg] : [msg])
  }


  // Watches information about the workspace, like users
  // returns a promise that is resolved when the workspace
  // Info is caught up, but still runs in the background
  // for updates
  const watchWorkspace = async () => {
    return new Promise(async (res, rej) => {
      const conn = store.conn
      if (!conn) {
        return
      }

      const js = conn.jetstream()
      const workspace = await js.views.kv("chat_workspace")
      const watcher = await workspace.watch({
        initializedFn: () => res(null)
      })

      for await (const entry of watcher) {
        const [resource, ...rest] = entry.key.split(".")

        switch (resource) {
          case "users":
            // Parse and add user to the users lookup table
            const id = rest[0]
            setStore("users", id, entry.json())
            console.log("setting user", entry.json())
            break;
        }
      }
    })

  }

  onMount(() => {
    (async () => {
      const conn = await connect({
        servers: ["ws://localhost:8222"],
      })
      setStore("conn", conn)

      const js = conn.jetstream()
      const consumer = await js.consumers.get("chat_messages")
      setStore("consumer", consumer)

      await watchWorkspace()

      const sub = await consumer.consume()
      for await (const m of sub) {
        onMessageReceived(m)
      }
    })()
  })

  onCleanup(async () => {
    console.log("closing connection...")
    await store.consumer?.delete()
    await store.conn?.close()
  })

  const sendMessage = (channel: string, message: string) => {
    console.log("sending message", channel, message)
    store.conn?.publish(`chat.${channel}.${store.userID}`, sc.encode(message))
  }

  const channelMessages = () => {
    return (store.messages[selected()] || []).map((m) => {
      return {
        ...m,
        user: store.users[m.userID] ?? {
          id: "unknown",
          name: "Unknown",
          email: "Unknown",
        }
      }
    })
  }

  return (
    <div class="inset-0 w-full h-lvh absolute flex flex-row">
      <Sidebar channels={channels} selected={selected()} onSelect={setSelected} />
      <ChannelView channel={selected()} onSend={sendMessage} messages={channelMessages()} />
    </div>
  );
}
