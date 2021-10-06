import { combineReducers, Reducer } from 'redux'

import authReducer from './auth/auth.slice'
import { RootState } from './types'

const rootReducer: Reducer<RootState> = combineReducers<RootState>({
  authState: authReducer,
})

export default rootReducer
