const path = require('path')

module.exports = {
    // Source files
    src: path.resolve(__dirname, '../src'),

    // Production build files
    build: path.resolve(__dirname, '../dist'),

    // Assets files that get copied to build folder
    assets: path.resolve(__dirname, '../src/assets'),

    // Env files
    dotenv: path.resolve(__dirname, '../.env'),
}