interface Message {
  user: {
    id: string
    email: string
    image: string
  }
  message: string
}

interface Props {
  channel: string
  messages: string[]
}

export default function Channel(props: Props) {
  return (
    <div></div>
  )
}
