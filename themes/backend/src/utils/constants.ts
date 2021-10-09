export const isProduction = process.env.NODE_ENV === 'production'
export const apiUrl = process.env.APP_API_URL
export const appName = process.env.APP_NAME

const api = '/backend/api'

export const endpoints = Object.freeze({
  auth: Object.freeze({ refresh: `${api}/refresh`, login: `${api}/sign-in` }),
  user: Object.freeze({ me: `${api}/me` }),
})

export const routes = Object.freeze({
  root: '/backend',
  error: Object.freeze({ notFound: '/backend/404' }),
  auth: Object.freeze({ login: '/backend/login' }),
})
