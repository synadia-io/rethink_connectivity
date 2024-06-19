import { createEffect, createSignal, onCleanup, onMount } from "solid-js";
import Sidebar from "./sidebar";
import Channel from "./channel";
import { StringCodec, connect, type NatsConnection } from "nats.ws";

interface Message {
  userId: string
  message: string
  timestamp: Date
}

type Channel = string

interface ChatStore {
  messages: Record<Channel, Message>
}

const channels = ["general", "random", "dev"]

export default function Chat() {
  const [userId, setUserId] = createSignal("foobar")
  const [selected, setSelected] = createSignal("general")
  const [conn, setConn] = createSignal<NatsConnection>()

  onMount(() => {
    (async () => {
      console.log("connecting...")
      const conn = await connect({
        servers: ["ws://localhost:8222"],
      })
      setConn(conn)

      const js = conn.jetstream()
      const consumer = await js.consumers.get("chat_messages")
      const sub = await consumer.consume()
      for await (const m of sub) {
        console.log(m)
      }
    })()
  })

  onCleanup(() => {
    console.log("closing connection...")
    conn()?.close()
  })

  const sendMessage = (channel: string, message: string) => {
    console.log("sending message", channel, message)
    const sc = StringCodec()
    conn()?.publish(`chat.${channel}.${userId()}`, sc.encode(message))
  }

  return (
    <div class="inset-0 w-full h-lvh absolute flex flex-row">
      <Sidebar channels={channels} selected={selected()} onSelect={setSelected} />
      <Channel channel={selected()} onSend={sendMessage} messages={[
      ]} />
    </div>
  );
}
