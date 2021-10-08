import { all, takeLatest } from 'redux-saga/effects'

import { login, logout } from '@/store/auth/auth.slice'
import { loginWorker, logoutWorker } from '@/store/auth/auth.saga'

function* rootSaga() {
  yield all([takeLatest(login.type, loginWorker), takeLatest(logout.type, logoutWorker)])
}

export default rootSaga
