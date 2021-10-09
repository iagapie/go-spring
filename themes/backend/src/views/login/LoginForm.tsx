import React, { ComponentPropsWithoutRef, useMemo } from 'react'
import { useForm } from 'react-hook-form'

import { SignIn } from '@/store/types'
import { validation } from '@/utils/constants'

export interface LoginFormProps extends ComponentPropsWithoutRef<'form'> {
  loading: boolean
  onLogin: (data: SignIn) => void
}

export const LoginForm: React.FC<LoginFormProps> = ({ loading, onLogin, ...rest }) => {
  const {
    register,
    handleSubmit,
    formState: { errors, isDirty, isValid },
  } = useForm({ mode: 'onChange' })

  const disabled = useMemo(() => loading || !isDirty || !isValid, [loading, isDirty, isValid])

  return (
    <form onSubmit={handleSubmit(onLogin)} className="form" {...rest}>
      <div className="form__row">
        <input
          type="email"
          placeholder="Enter email"
          className={errors.email ? 'form__input form__input_error' : 'form__input'}
          {...register('email', {
            required: 'required',
            pattern: {
              value: validation.email,
              message: 'email is not valid.',
            },
          })}
        />
        {errors.email && <strong className="form__error">{errors.email.message}</strong>}
      </div>
      <div className="form__row">
        <input
          type="password"
          placeholder="Enter password"
          className={errors.password ? 'form__input form__input_error' : 'form__input'}
          {...register('password', {
            required: 'required',
            pattern: {
              value: validation.password,
              message: 'password is not valid.',
            },
          })}
        />
        {errors.password && <strong className="form__error">{errors.password.message}</strong>}
      </div>
      <div className="form__row">
        <button type="submit" className="form__btn" disabled={disabled}>
          {loading && <span className="rotate-loading form__loading" />}
          Log in
        </button>
      </div>
    </form>
  )
}
