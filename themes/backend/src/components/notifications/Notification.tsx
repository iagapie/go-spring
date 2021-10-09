import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useDispatch } from 'react-redux'
import { XIcon } from '@heroicons/react/outline'

import * as t from '@/store/types'
import { removeNotification } from '@/store/notifications/notifications.slice'

export interface NotificationProps {
  notification: t.Notification
}

export const Notification: React.FC<NotificationProps> = ({ notification }) => {
  const rootElementRef = useRef<HTMLDivElement>(null!)
  const [parentStyle, setParentStyle] = useState({ height: '0px', overflow: 'auto', transition: 'none' })
  // switch (notification.type) {
  //   case NotifyType.success:
  //     return 'notification__container_success'
  //   case NotifyType.warning:
  //     return 'notification__container_warning'
  //   case NotifyType.error:
  //     return 'notification__container_error'
  //   default:
  //     return 'notification__container_info'
  // }
  const classType = useMemo(() => `notification__container_${notification.type}`, [notification.type])

  const dispatch = useDispatch()

  const onClose = useCallback(() => {
    setTimeout(() => dispatch(removeNotification(notification)), 150)
    setParentStyle({
      height: '0px',
      overflow: 'hidden',
      transition: 'height 150ms linear 0ms',
    })
  }, [dispatch, notification])

  useEffect(() => {
    const { scrollHeight } = rootElementRef.current
    setParentStyle({
      overflow: 'auto',
      height: `${scrollHeight}px`,
      transition: 'height 300ms linear 0ms',
    })
  }, [rootElementRef, setParentStyle])

  useEffect(() => {
    if (notification.dismiss! > 0) {
      const timeout = setTimeout(() => onClose(), notification.dismiss)

      return () => {
        clearTimeout(timeout)
      }
    }
  }, [notification.dismiss, onClose])

  return (
    <div ref={rootElementRef} className="notification" style={parentStyle}>
      <div onClick={onClose} className={`notification__container ${classType}`}>
        {notification.title && (
          <div className="notification__header">
            <div className="notification__title">{notification.title}</div>
            <XIcon className="notification__close" />
          </div>
        )}
        <div className="notification__message">{notification.message}</div>
      </div>
    </div>
  )
}
