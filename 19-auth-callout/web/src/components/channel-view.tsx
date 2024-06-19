import { SendHorizontalIcon } from "lucide-solid"
import { createSignal } from "solid-js"
import type { MessageWithUser } from "../types"

interface Props {
  channel: string
  messages: MessageWithUser[]
  onSend: (channel: string, message: string) => void
}

export default function ChannelView(props: Props) {
  const [message, setMessage] = createSignal("")

  const onSubmit = (e: Event) => {
    e.preventDefault()
    props.onSend(props.channel, message())
    setMessage("")
  }

  const formatDate = (d: Date) => {
    return d.toLocaleTimeString('en-US', { hour: "2-digit", minute: "2-digit" })
  }

  return (
    <div class="w-full h-full flex flex-col">
      <div class="p-4 border-b border-zinc-800">
        <span class="text-xl font-medium"># {props.channel}</span>
      </div>

      <div class="flex-grow flex flex-col justify-end p-6 gap-4">
        {props.messages.map((msg) => (
          <div class="flex flex-row gap-2">
            <img class="w-10 h-10 mt-1 rounded" src={msg.user.photoURL} />
            <div class="flex flex-col">
              <div class="flex flex-row gap-2 items-center">
                <span class="text-zinc-100 font-semibold">{msg.user.name}</span>
                <span class="text-zinc-400 text-xs mt-0.5">{formatDate(msg.timestamp)}</span>
              </div>
              <span class="text-zinc-300">{msg.text}</span>
            </div>
          </div>
        ))}
      </div>

      <div class="p-4">
        <form onSubmit={onSubmit} class="flex flex-row w-full items-center justify-center p-4 bg-zinc-800 border border-zinc-700 rounded-lg">
          <input
            autofocus
            value={message()}
            onInput={(e) => setMessage(e.target.value)}
            class="flex-grow bg-transparent text-zinc-100 placeholder-zinc-400 focus:outline-none"
            placeholder={`Message #${props.channel}`}
          />
          <button onClick={onSubmit} class="text-zinc-200 rounded-lg border border-zinc-700/0 hover:border-zinc-700 py-1 px-2">
            <SendHorizontalIcon class="w-4" />
          </button>
        </form>
      </div>
    </div>
  )
}
