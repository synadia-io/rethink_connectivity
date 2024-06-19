import { createSignal } from "solid-js"

interface Props {
  onSubmit: (email: string) => void
}

export default function Login() {
  const [email, setEmail] = createSignal("")

  const onSubmit = (e: Event) => {
    e.preventDefault()

  }

  return (
    <div class="w-full h-full bg-zinc-900 flex flex-row items-center justify-center">
      <form>
        <span>Enter your email to log in</span>
        <input value={email()} onInput={(e) => setEmail(e.target.value)} type="text" placeholder="Email address" />
      </form>
    </div>
  )
}
