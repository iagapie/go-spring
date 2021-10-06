import React from 'react'
import SplitPane from 'react-split-pane'
import { Router } from 'react-router'

import history from '@/utils/history'

export const App: React.FC = () => {
  return (
    <Router history={history}>
      <SplitPane split="vertical" defaultSize={200} primary="first">
        <div>
          <h1>Spring CMS</h1>
        </div>
        <div>
          <p>Hello Backend</p>
        </div>
      </SplitPane>
    </Router>
  )
}
