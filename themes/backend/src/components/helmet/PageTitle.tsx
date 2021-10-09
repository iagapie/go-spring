import React, { useMemo } from 'react'
import { Helmet } from 'react-helmet'

import { appName } from '@/utils/constants'

export interface PageTitleProps {
  title?: string
  separator?: string
}

export const PageTitle: React.FC<PageTitleProps> = ({ title, separator }) => {
  const _title = useMemo(
    () => [title, 'Backend', appName].filter((i) => !!i).join(separator ?? ' | '),
    [title, separator]
  )

  return (
    <Helmet>
      <title>{_title}</title>
    </Helmet>
  )
}
