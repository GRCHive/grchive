import hljs from 'highlight.js/lib/highlight'
import sql from 'highlight.js/lib/languages/sql'
import 'highlight.js/styles/vs2015.css'

hljs.registerLanguage('sql', sql)
hljs.configure({
    useBR: false,
    languages: ['sql']
})

//@ts-ignore
window.hljs = hljs
