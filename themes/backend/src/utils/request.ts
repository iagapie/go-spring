import { apiUrl, endpoints } from './constants'

import { getAuth } from '@/store/selectors'
import { store } from '@/store'
import { clearAuth, setTokens } from '@/store/auth/auth.slice'
import { ErrorResp } from '@/store/types'

export type Options = { token?: boolean; method?: string; headers?: HeadersInit; body?: BodyInit | null }
export type ResponseBody<T> = T | string | null

export interface Resp<T> {
  response: Response
  body: ResponseBody<T>
}

export class ResponseError extends Error {
  public readonly status: number
  public readonly response: Response
  public readonly resp: ResponseBody<ErrorResp>

  constructor(status: number, response: Response, resp: ResponseBody<ErrorResp>) {
    super()
    this.name = 'ResponseError'
    this.status = status
    this.response = response
    this.resp = resp
  }
}

const bearer = (token?: string) => `Bearer ${token}`

function getResponseBody<T>(response: Response): Promise<ResponseBody<T>> {
  const contentType = response.headers.get('Content-Type')
  const text = response.clone().text()

  return contentType && contentType.indexOf('json') >= 0 ? text.then((json: string) => tryParseJSON<T>(json)) : text
}

function tryParseJSON<T>(json: string): T | null {
  if (!json) {
    return null
  }

  try {
    return JSON.parse(json) as T
  } catch {
    throw new Error(`Failed to parse unexpected JSON response: ${json}`)
  }
}

function refresh(): Promise<void> {
  const { tokens } = getAuth(store.getState())

  return requestJSON<any>(endpoints.auth.refresh, {
    method: 'POST',
    body: JSON.stringify({ token: tokens?.refreshToken }),
  })
    .then((resp) => {
      store.dispatch(
        setTokens({
          accessToken: resp.body.access_token,
          refreshToken: resp.body.refresh_token,
        })
      )
    })
    .catch((error) => {
      store.dispatch(clearAuth())
      // TODO
      // history.push(routes.auth.login)
      throw error
    })
}

function req<T>(url: string | URL, options?: Options): Promise<Resp<T>> {
  if (!(url instanceof URL)) {
    url = new URL(url, apiUrl)
  }
  options = { token: false, method: 'GET', ...options }
  options.headers = new Headers(options.headers)

  if (options.token) {
    const { tokens } = getAuth(store.getState())
    options.headers.set('Authorization', bearer(tokens?.accessToken))
  }
  delete options.token

  return fetch(url.href, {
    mode: 'cors',
    credentials: 'include',
    ...options,
  }).then((response: Response) => {
    if (response.ok) {
      return getResponseBody<T>(response).then((body) => ({ response, body }))
    } else {
      return getResponseBody<ErrorResp>(response).then((body) => {
        throw new ResponseError(response.status, response, body)
      })
    }
  })
}

let refreshingTokenPromise: Promise<void> | null = null

export function request<T>(url: string | URL, options?: Options): Promise<Resp<T>> {
  if (refreshingTokenPromise !== null) {
    return refreshingTokenPromise.then(() => req<T>(url, options)).catch(() => req<T>(url, options))
  }

  return req<T>(url, options).catch((error) => {
    const err = error as ResponseError

    if (err.status === 401) {
      if (refreshingTokenPromise === null) {
        refreshingTokenPromise = new Promise<void>((resolve, reject) => {
          refresh()
            .then(() => {
              refreshingTokenPromise = null
              resolve()
            })
            .catch((refreshTokenError) => {
              refreshingTokenPromise = null
              reject(refreshTokenError)
            })
        })
      }

      return refreshingTokenPromise
        .catch(() => {
          throw error
        })
        .then(() => req<T>(url, options))
    }

    throw error
  })
}

export function requestJSON<T>(url: string | URL, options?: Options): Promise<Resp<T>> {
  const headers = new Headers(options?.headers)
  headers.set('Content-Type', 'application/json')
  headers.set('Accept', 'application/json')
  options = { ...options, headers }

  return request(url, options)
}
