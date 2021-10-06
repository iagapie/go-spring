import { FC } from 'react'
import { Redirect, Route, RouteProps } from 'react-router'
import { useSelector } from 'react-redux'

import { getAuth } from '../../store/selectors'

interface PrivateRouteProps extends RouteProps {
  children: any
}

export const PrivateRoute: FC<PrivateRouteProps> = ({children: Component, ...rest}) => {
  const { } = useSelector(getAuth)

  return <Route render={() => (true ? Component : <Redirect to="/login" />)} {...rest} />
}
