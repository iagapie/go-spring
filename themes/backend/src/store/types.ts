export interface BackendUser {
  [anyProp: string]: any
}

export interface AuthState {
  isAuthenticated: boolean
  currentUser: BackendUser
  accessToken: string
  refreshToken: string
  csrf: string
  loading: boolean
}

export interface RootState {
  authState: AuthState
}
