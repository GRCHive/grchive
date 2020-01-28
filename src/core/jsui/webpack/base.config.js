const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const VueLoaderPlugin = require('vue-loader/lib/plugin')
const VuetifyLoaderPlugin = require('vuetify-loader/lib/plugin')
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')

const cssLoaders = [
    {
        loader: MiniCssExtractPlugin.loader
    },
    'css-loader'
]

const babelLoader = {
    loader: 'babel-loader',
    options: {
        presets: ['@babel/preset-env']
    }
}

const badNodeModulesRegex = /node_modules\/(?!(query-string|split-on-first|strict-uri-encode|vuetify)\/).*/

module.exports = {
    entry: {
        main: ['./ts/main.ts'],
        dashboard: ['./ts/dashboard.ts'],
    },
    output: {
        path: path.resolve(__dirname, '../dist'),
        libraryTarget: 'umd',
        library: 'corejsui',
        libraryExport: 'default',
        publicPath: '/static/corejsui/'
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: [
                    "cache-loader",
                    "thread-loader",
                    babelLoader,
                    {
                        loader: 'ts-loader',
                        options: {
                            appendTsSuffixTo: [/\.vue$/],
                            happyPackMode: true,
                        }
                    }
                ],
                exclude: [ 
                    badNodeModulesRegex
                ],
            },
            {
                test: /\.vue$/,
                use: 'vue-loader'
            },
            {
                test: /\.css$/,
                use: cssLoaders
            },
            {
                test: /\.scss$/,
                use: [
                    ...cssLoaders,
                    {
                        loader: 'sass-loader',
                        options: {
                            implementation: require('sass'),
                            sassOptions: {
                                fiber: require('fibers'),
                                indentedSyntax: false
                            }
                        }
                    }
                ]
            },
            {
                test: /\.sass$/,
                use: [
                    ...cssLoaders,
                    {
                        loader: 'sass-loader',
                        options: {
                            implementation: require('sass'),
                            sassOptions: {
                                fiber: require('fibers'),
                                indentedSyntax: true
                            }
                        }
                    }
                ]
            },
            {
                test: /\.jsx?$/,
                use: [
                    "cache-loader",
                    "thread-loader",
                    babelLoader,
                ],
                exclude: [ 
                    badNodeModulesRegex
                ],
            },
            {
                test: /\.(eot|ttf|woff|woff2|png|gif|svg)$/,
                use: 'url-loader'
            }
        ]
    },
    plugins: [
        new VueLoaderPlugin(),
        new VuetifyLoaderPlugin(),
        new ForkTsCheckerWebpackPlugin({
            checkSyntacticErrors: true,
            vue: true,
        }),
        new HtmlWebpackPlugin({
            chunks: ['main', 'vendors'],
            filename: 'main.tmpl',
            inject: false,
            template: 'loader.ejs',
            minify: false,
        }),
        new HtmlWebpackPlugin({
            chunks: ['dashboard', 'vendors'],
            filename: 'dashboard.tmpl',
            inject: false,
            template: 'loader.ejs',
            minify: false,
        }),
    ],
    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js'
        },
        extensions: ['.ts', '.js', '.vue', '.json',]
    },
};
