import React, { useEffect } from 'react'
import SplitPane from 'react-split-pane'
import { Router } from 'react-router'
import { useDispatch, useSelector } from 'react-redux'

import history from '@/utils/history'
import { login } from '@/store/auth/auth.slice'
import { SignIn } from '@/store/types'
import { getAuth } from '@/store/selectors'
import userService from '@/services/User/User.service'

export const App: React.FC = () => {
  const dispatch = useDispatch()
  const { currentUser } = useSelector(getAuth)

  // const _login = (signIn: SignIn) => dispatch(login(signIn))
  //
  // useEffect(() => {
  //   _login({
  //     email: 'iagapie@gmail.com',
  //     password: 'Admin123',
  //   })
  // }, [_login])

  useEffect(() => {
    userService.me().then((u) => {
      console.log(u)
    })
  }, [userService])

  return (
    <Router history={history}>
      <SplitPane split="vertical" defaultSize={200} primary="first">
        <div>
          <h1>Spring CMS</h1>
        </div>
        <div>
          <h2>{currentUser?.email}</h2>
          <p>Hello Backend</p>
        </div>
      </SplitPane>
    </Router>
  )
}
