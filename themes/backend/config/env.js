const fs = require('fs')
const dotenv = require('dotenv')
const dotenvExpand = require('dotenv-expand')
const paths = require('./paths')

const { NODE_ENV } = process.env

// if (!NODE_ENV) {
//   throw new Error(
//     'The NODE_ENV environment variable is required but was not specified.',
//   )
// }

// https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use
const dotenvFiles = [
  `${paths.dotenv}.${NODE_ENV}.local`,
  // Don't include `.env.local` for `test` environment
  // since normally you expect tests to produce the same
  // results for everyone
  NODE_ENV !== 'test' && `${paths.dotenv}.local`,
  `${paths.dotenv}.${NODE_ENV}`,
  paths.dotenv,
].filter(Boolean)

// Load environment variables from .env* files. Suppress warnings using silent
// if this file is missing. dotenv will never modify any environment variables
// that have already been set.  Variable expansion is supported in .env files.
// https://github.com/motdotla/dotenv
// https://github.com/motdotla/dotenv-expand
dotenvFiles.forEach((dotenvFile) => {
  if (fs.existsSync(dotenvFile)) {
    dotenvExpand(
      dotenv.config({
        path: dotenvFile,
      })
    )
  }
})

// Grab NODE_ENV and REACT_APP_* environment variables and prepare them to be
// injected into the application via DefinePlugin in webpack configuration.
const APP = /^APP_/i

const nodeEnv = process.env.NODE_ENV || 'development'

const getClientEnvironment = () => {
  const raw = Object.keys(process.env)
    .filter((key) => APP.test(key))
    .reduce((obj, key) => ({ ...obj, [key]: process.env[key] }), {
      BABEL_ENV: process.env.BABEL_ENV || 'development',
      NODE_ENV: nodeEnv,
      DEBUG: process.env.DEBUG || nodeEnv === 'development',
    })

  const stringified = {
    'process.env': Object.keys(raw).reduce((obj, key) => ({ ...obj, [key]: JSON.stringify(raw[key]) }), {}),
  }

  return { raw, stringified }
}

module.exports = getClientEnvironment
