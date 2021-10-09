import React from 'react'
import { useDispatch, useSelector } from 'react-redux'

import { LoginForm } from './LoginForm'

import { login } from '@/store/auth/auth.slice'
import { SignIn } from '@/store/types'
import { getAuth } from '@/store/selectors'
import { PageTitle } from '@/components/helmet/PageTitle'
import { appName } from '@/utils/constants'

const LoginPage: React.FC = () => {
  const { loading } = useSelector(getAuth)
  const dispatch = useDispatch()

  const onLogin = (data: SignIn) => {
    dispatch(login(data))
  }

  return (
    <main className="auth">
      <PageTitle title="Login" />
      <div className="auth__container">
        <div className="auth__header">
          <h1 className="auth__title">{appName}</h1>
        </div>
        <LoginForm className="auth__form" loading={loading} onLogin={onLogin} />
      </div>
    </main>
  )
}

export default LoginPage
