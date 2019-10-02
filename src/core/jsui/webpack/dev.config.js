const base = require('./base.config.js')
const merge = require('webpack-merge')

module.exports = merge(base, {
    mode: 'development'
})
