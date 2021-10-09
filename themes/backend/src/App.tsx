import React, { Suspense, lazy } from 'react'
import { Router, Switch, Route } from 'react-router'

import history from '@/utils/history'
import { routes } from '@/utils/constants'
import { PageTitle } from '@/components/helmet/PageTitle'
import { NotificationContainer } from '@/components/notifications/NotificationContainer'
import { Loading } from '@/components/loading/Loading'
import { PublicRoute } from '@/components/routing/PublicRoute'
import { PrivateRoute } from '@/components/routing/PrivateRoute'
import NotFoundPage from '@/views/notFound/NotFoundPage'

const DashboardPage = lazy(() => import(/* webpackChunkName: "dashboard" */ '@/views/dashboard/DashboardPage'))
const LoginPage = lazy(() => import(/* webpackChunkName: "login" */ '@/views/login/LoginPage'))

export const App: React.FC = () => (
  <Router history={history}>
    <PageTitle />
    <NotificationContainer />
    <Suspense fallback={<Loading />}>
      <Switch>
        <PrivateRoute exact path={routes.root}>
          <DashboardPage />
        </PrivateRoute>
        <PublicRoute exact path={routes.auth.login}>
          <LoginPage />
        </PublicRoute>
        <Route path="*">
          <NotFoundPage />
        </Route>
      </Switch>
    </Suspense>
  </Router>
)
