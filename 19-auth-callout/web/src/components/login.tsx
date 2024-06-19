import { createSignal } from "solid-js"

interface Props {
  onSubmit: (email: string) => void
}

export default function Login(props: Props) {
  const [email, setEmail] = createSignal("")

  const onSubmit = (e: Event) => {
    e.preventDefault()
    const emailAddress = email()

    if (emailAddress == "") {
      return
    }

    setEmail("")
    props.onSubmit(emailAddress)
  }

  return (
    <div class="inset-0 w-full h-lvh absolute bg-zinc-900">
      <form onSubmit={onSubmit} class="w-full h-full flex flex-col items-center justify-center">
        <div class="border border-zinc-800 rounded p-12 flex flex-col gap-4">
          <span class="text-2xl font-bold">Enter your email to log in</span>
          <input
            class="border bg-zinc-800 border-zinc-700 rounded px-4 py-3 text-lg focus:outline-none focus:border-zinc-500"
            value={email()}
            onInput={(e) => setEmail(e.target.value)} type="text" placeholder="Email address" />
          <button
            disabled={email() == ""}
            class="bg-zinc-100 text-zinc-800 px-4 py-3 text-lg rounded font-semibold hover:bg-zinc-300 disabled:bg-zinc-100/50" type="submit">Log in</button>
        </div>
      </form>
    </div>
  )
}
