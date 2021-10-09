import React, { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import SplitPane from 'react-split-pane'

import { getAuth } from '../../store/selectors'
import userService from '../../services/User/User.service'

import { login } from '@/store/auth/auth.slice'

const DashboardPage: React.FC = () => {
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
    <SplitPane split="vertical" defaultSize={200} primary="first">
      <div>
        <h1>Spring CMS</h1>
      </div>
      <div>
        <h2>{currentUser?.email}</h2>
        <p>Hello Backend</p>
      </div>
    </SplitPane>
  )
}

export default DashboardPage
