import { createSignal } from "solid-js";
import Sidebar from "./sidebar";

const channels = ["general", "random", "dev"]

export default function Chat() {
  const [selected, setSelected] = createSignal("general")

  return (
    <div class="inset-0 w-full h-lvh absolute flex flex-row">
      <Sidebar channels={channels} selected={selected()} onSelect={setSelected} />
    </div>
  );
}
