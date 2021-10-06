import React from 'react'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import 'regenerator-runtime/runtime'

import '@/styles/index.scss'

import { sagaMiddleware, store } from '@/store'
import rootSaga from '@/store/rootSaga'
import { App } from '@/App'

sagaMiddleware.run(rootSaga)

render(
  <Provider store={store}>
    <App/>
  </Provider>,
  document.getElementById('root'),
)
