<template>
    <div ref="parent">
        <textarea ref="editor"></textarea>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import CodeMirror from 'codemirror/lib/codemirror.js'

import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/material.css'
import 'codemirror/addon/display/autorefresh.js'

import ( /* webpackChunkName: "codemirrorClikeMode" */ 'codemirror/mode/clike/clike.js')
import ( /* webpackChunkName: "codemirrorSqlMode" */ 'codemirror/mode/sql/sql.js')
import ( /* webpackChunkName: "codemirrorShellMode" */ 'codemirror/mode/shell/shell.js')
import ( /* webpackChunkName: "codemirrorPowershellMode" */ 'codemirror/mode/powershell/powershell.js')

export const Props = Vue.extend({
    props: {
        value: {
            type: String,
            default: "",
        },
        lang: {
            type: String,
            default: "",
        },
        readonly: {
            type: Boolean,
            default: false,
        },
        fullHeight: {
            type: Boolean,
            default: false,
        },
        fixedWidth: {
            type: Number,
            default: -1,
        },
    }
})

@Component
export default class GenericCodeEditor extends Props {
    cm : any = null

    $refs! : {
        parent : HTMLElement
        editor : HTMLElement
    }

    get desiredHeight() : string {
        let bb = this.$refs.parent.getBoundingClientRect()
        // Subtract 30px to account for the 30px bottom padding
        // that CodeMirror adds.
        return `calc(100vh -  ${bb.top}px - 30px)`
    }

    mounted() {
        this.cm = CodeMirror.fromTextArea(
            this.$refs.editor,
            {
                mode: this.lang,
                theme: "material",
                indentUnit: 4,
                lineNumbers: true,
                readOnly: this.readonly,
                value: this.value,
                autoRefresh: true
            }
        )

        this.cm!.on('change', (instance : any, obj : any) => {
            this.$emit('input', this.cm!.getValue())
        })

        if (this.fullHeight) {
            this.cm!.setSize(null, this.desiredHeight)
        }

        this.refreshWidth()
        this.cm!.setValue(this.value)
    }

    @Watch('fixedWidth')
    refreshWidth() {
        if (this.fixedWidth != -1) {
            this.cm!.setSize(this.fixedWidth, null)
        }
    }

    @Watch('readonly')
    watchReadonly() {
        this.cm!.setOption('readOnly', this.readonly)
    }

    @Watch('value')
    watchValue() {
        if (this.cm!.getValue() != this.value) {
            this.cm!.setValue(this.value)
        }
    }
}

</script>
