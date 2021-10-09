import React from 'react'
import { Redirect, Route, RouteProps } from 'react-router'
import { useSelector } from 'react-redux'

import { getAuth } from '@/store/selectors'
import { routes } from '@/utils/constants'

export interface PrivateRouteProps extends RouteProps {
  children: any
}

export const PrivateRoute: React.FC<PrivateRouteProps> = ({ children: Component, ...rest }) => {
  const { isAuthenticated } = useSelector(getAuth)

  return <Route render={() => (isAuthenticated ? Component : <Redirect to={routes.auth.login} />)} {...rest} />
}
