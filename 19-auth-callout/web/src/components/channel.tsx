import { SendHorizonal } from "lucide-solid"

interface Message {
  user: User
  message: string
}

interface User {
  id: string
  name: string
  email: string
  photoURL: string
}

interface Props {
  channel: string
  messages: Message[]
}

export default function Channel(props: Props) {
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
              <span class="text-zinc-100 font-semibold">{msg.user.name}</span>
              <span class="text-zinc-300">{msg.message}</span>
            </div>
          </div>
        ))}
      </div>

      <div class="p-4">
        <div class="flex flex-row w-full items-center justify-center p-4 bg-zinc-800 border border-zinc-700 rounded-lg">
          <input
            class="flex-grow bg-transparent text-zinc-100 placeholder-zinc-400 focus:outline-none"
            placeholder={`Message #${props.channel}`}
          />
          <button class="text-zinc-200 rounded-lg border border-zinc-700/0 hover:border-zinc-700 py-1 px-2">
            <SendHorizonal class="w-4" />
          </button>
        </div>
      </div>
    </div>
  )
}
