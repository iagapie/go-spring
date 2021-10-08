export const isProduction = process.env.NODE_ENV === 'production'
export const apiUrl = process.env.APP_API_URL

export const endpoints = Object.freeze({
  auth: Object.freeze({ refresh: '/backend/refresh', login: '/backend/sign-in' }),
  user: Object.freeze({ me: '/backend/me' }),
})
