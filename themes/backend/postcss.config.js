const flexboxFixes = require('postcss-flexbugs-fixes')
const presetEnv = require('postcss-preset-env')

module.exports = {
  plugins: [
    flexboxFixes,
    presetEnv({
      browsers: 'last 2 versions',
      autoprefixer: {
        flexbox: 'no-2009',
        // grid: 'autoplace',
      },
      stage: 3,
    }),
  ],
}
