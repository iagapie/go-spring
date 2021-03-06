import { call, put } from 'redux-saga/effects'
import { PayloadAction } from '@reduxjs/toolkit'

import { BackendUser, SignIn, Tokens } from '../types'
import authService from '../../services/Auth/Auth.service'
import userService from '../../services/User/User.service'

import { clearAuth, loginSuccess, setCurrentUser, setTokens } from './auth.slice'

import { routes } from '@/utils/constants'
import history from '@/utils/history'
import { addError } from '@/store/notifications/notifications.slice'

function* fetchTokens(service: any, payload: SignIn) {
  const tokens: Tokens = yield call(service, payload)
  yield put(setTokens(tokens))
}

function* fetchCurrentUser() {
  const currentUser: BackendUser = yield call(userService.me)
  yield put(setCurrentUser(currentUser))
}

function* authWorker(service: any, payload: SignIn) {
  try {
    yield call(fetchTokens, service, payload)
    yield call(fetchCurrentUser)
    yield put(loginSuccess())
  } catch (error) {
    yield put(addError({ message: (error as Error).message }))
    yield put(clearAuth())
  }
}

export function* loginWorker({ payload }: PayloadAction<SignIn>) {
  yield call(authWorker, authService.login, payload)
}

export function* logoutWorker() {
  yield put(clearAuth())
  history.push(routes.auth.login)
}
