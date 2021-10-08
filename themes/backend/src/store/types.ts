export interface ErrorResp {
  message: string
  error?: string
}

export interface SignIn {
  email: string
  password: string
}

export interface Tokens {
  accessToken: string
  refreshToken: string
}

export interface BackendUser {
  uuid: string
  name: string
  email: string
  createdAt: string
  updatedAt: string
}

export interface AuthState {
  isAuthenticated: boolean
  currentUser?: BackendUser
  tokens?: Tokens
  loading: boolean
}

export interface RootState {
  authState: AuthState
}
