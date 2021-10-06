import { createSlice } from '@reduxjs/toolkit'

import { AuthState } from '../types'
import localStorageService from '@/services/LocalStorage/LocalStorage.service'
import cookiesService from '@/services/Cookies/Cookies.service'

const keyUser = 'auth.user'
const keyTokens = 'tokens'

const currentUser = localStorageService.get(keyUser, {})
const {accessToken, refreshToken} = cookiesService.get(keyTokens, {accessToken: '', refreshToken: ''})
const isAuthenticated = Object.keys(currentUser).length !== 0 && !!accessToken && !!refreshToken

const initialState: AuthState = {
  isAuthenticated,
  currentUser: isAuthenticated ? currentUser : {},
  accessToken: isAuthenticated ? accessToken : '',
  refreshToken: isAuthenticated ? refreshToken : '',
  csrf: '',
  loading: false,
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {},
})

export const {} = authSlice.actions

export default authSlice.reducer
