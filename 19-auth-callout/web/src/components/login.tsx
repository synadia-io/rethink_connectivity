import { Show, createEffect, createSignal, onMount } from "solid-js"
import { parseJWT } from "../util/jwt"

interface Props {
  onSubmit: (email: string) => void
  useSSO?: boolean
}

export default function Login(props: Props) {
  const [email, setEmail] = createSignal("")
  let ssoButton: HTMLDivElement

  const onSubmit = (e: Event) => {
    e.preventDefault()
    const emailAddress = email()

    if (emailAddress == "") {
      return
    }

    setEmail("")
    props.onSubmit(emailAddress)
  }

  onMount(() => {
    if (props.useSSO) {
      /*@ts-ignore */
      google.accounts.id.initialize({
        client_id: "476732082534-gpp6r1t67lirjfddjckm81vr6c2kikr1.apps.googleusercontent.com",
        callback: (res: any) => {
          if (res.credential) {
            const claims = parseJWT(res.credential)
            props.onSubmit(claims.email)
          }
        }
      })
      /*@ts-ignore */
      google.accounts.id.renderButton(
        ssoButton,
        { theme: "outline", size: "large" }
      );
      /*@ts-ignore */
      google.accounts.id.prompt()
    }
  })


  return (
    <div class="inset-0 w-full h-lvh absolute bg-zinc-900">

      <Show when={props.useSSO}>
        <div class="w-full h-full flex flex-col items-center justify-center">
          <div class="border border-zinc-800 rounded p-12 flex flex-col gap-4 items-center">
            <span class="text-2xl font-bold">Sign in to NATS Chat</span>
            <div ref={ssoButton}></div>
          </div>
        </div>
      </Show>

      <Show when={!props.useSSO}>
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
      </Show >
    </div>
  )
}
