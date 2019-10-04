const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const VueLoaderPlugin = require('vue-loader/lib/plugin');

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

module.exports = {
    devtool: "source-map",
    entry: {
        main: ['./ts/main.ts'],
        dashboard: ['./ts/dashboard.ts'],
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, '../dist'),
        libraryTarget: 'umd',
        library: 'corejsui',
        libraryExport: 'default',
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: [
                    {
                        loader: 'ts-loader',
                        options: {
                            appendTsSuffixTo: [/\.vue$/]
                        }
                    }
                ],
                exclude: [ 
                    /node_modules/,
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
                test: /\.s[ac]ss$/,
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
                test: /\.jsx?$/,
                use: babelLoader,
                exclude: [ 
                    /node_modules/,
                ],
            },
            {
                test: /\.(eot|ttf|woff|woff2)$/,
                use: 'url-loader'
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: '[name].css'
        }),
        new VueLoaderPlugin()
    ],
    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js'
        },
        extensions: ['.ts', '.js', '.vue', '.json',]
    },
};
