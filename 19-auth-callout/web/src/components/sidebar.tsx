import { For } from "solid-js"
import cn from "../util/styles"

interface Props {
  channels: string[]
  selected?: string

  onSelect: (channel: string) => void
}

export default function Sidebar(props: Props) {
  return (
    <div class="border-r border-zinc-800 h-full w-64 p-4 flex flex-col gap-2">
      <h1 class="text-xl text-white font-semibold">NATS Chat</h1>
      <ul class="w-full">
        <For each={props.channels}>
          {(channel) => (
            <li onClick={() => props.onSelect(channel)} class={cn(
              "rounded-lg px-2 py-1 -mx-2 cursor-pointer",
              props.selected === channel ? "bg-purple-800 hover:bg-purple-800" : "hover:bg-zinc-800"
            )
            }># {channel}</li>
          )}
        </For>
      </ul>
    </div>
  )
}
