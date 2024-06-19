import { createEffect, createSignal, onCleanup, onMount } from "solid-js";
import Sidebar from "./sidebar";
import Channel from "./channel";
import { connect, type NatsConnection } from "nats.ws";

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
  const [selected, setSelected] = createSignal("general")
  const [conn, setConn] = createSignal<NatsConnection>()

  onMount(async () => {
    console.log("connecting...")
    const conn = await connect({
      servers: ["ws://localhost:8222"],
    })
    setConn(conn)

    const jsm = await conn.jetstreamManager()
  })

  onCleanup(() => {
    console.log("closing connection...")
    conn()?.close()
  })

  return (
    <div class="inset-0 w-full h-lvh absolute flex flex-row">
      <Sidebar channels={channels} selected={selected()} onSelect={setSelected} />
      <Channel channel={selected()} messages={[
      ]} />
    </div>
  );
}
