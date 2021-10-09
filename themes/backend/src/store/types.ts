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

export enum NotifyType {
  Info = 'info',
  Success = 'success',
  Warning = 'warning',
  Error = 'error',
}

export interface Notification {
  id?: string
  type?: NotifyType
  title?: string
  message: string
  dismiss?: number
}

export interface NotificationsState {
  notifications: Notification[]
}

export interface RootState {
  authState: AuthState
  notificationsState: NotificationsState
}
