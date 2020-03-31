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

import ( /* webpackChunkName: "codemirrorClikeMode" */ 'codemirror/mode/clike/clike.js')

const Props = Vue.extend({
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
        }
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
            }
        )

        this.cm!.on('change', (instance : any, obj : any) => {
            this.$emit('input', this.cm!.getValue())
        })

        this.cm!.setSize(null, this.desiredHeight)
    }

    @Watch('readonly')
    watchReadonly() {
        this.cm!.setOption('readOnly', this.readonly)
    }

    @Watch('value')
    watchValue() {
        this.cm!.setValue(this.value)
    }
}

</script>
