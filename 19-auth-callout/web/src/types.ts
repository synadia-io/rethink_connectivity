// Name of a channel. This may become more structured in more advanced use cases
export type Channel = string

// Represents a user ID, which in this example is a base64
// encoded email address (so it can be used inside NATS subjects)
export type UserID = string

// Represents a User
export interface User {
  id: string
  name: string
  email: string
  photoURL: string
}

export interface Message {
  user: User
  text: string
  timestamp: Date
}
