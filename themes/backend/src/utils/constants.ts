export const isProduction = process.env.NODE_ENV === 'production'
export const apiUrl: string = process.env.APP_API_URL!
export const appName: string = process.env.APP_NAME!

const api = '/backend/api'

export const endpoints = Object.freeze({
  auth: Object.freeze({ refresh: `${api}/refresh`, login: `${api}/sign-in` }),
  user: Object.freeze({ me: `${api}/me` }),
})

export const routes = Object.freeze({
  root: '/backend',
  error: Object.freeze({ notFound: '/backend/404' }),
  auth: Object.freeze({ login: '/backend/login' }),
  cms: Object.freeze({
    pages: '/backend/cms',
    partials: '/backend/cms/partials',
    layouts: '/backend/cms/layouts',
    assets: '/backend/cms/assets',
    components: '/backend/cms/components',
  }),
})

export const validation = Object.freeze({
  email: /\S+@\S+\.\S+/,
  password: /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[\w\s\-^$&*!@#]{8,64}$/,
})
