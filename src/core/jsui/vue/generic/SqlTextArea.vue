<template>
    <div class="hljs" id="editor">
        <pre><code>{{ value }}</code></pre>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'

import Quill from 'quill'

const Props = Vue.extend({
    props: {
        value : {
            type: String,
            default: ""
        },
        readonly: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class SqlTextArea extends Props {
    quill : Quill | null = null

    createQuillEditor() {
        let options = {
            modules: {
                syntax: true,
                toolbar: false,
                keyboard: {}
            },
            placeholder: 'Your SQL query goes here...',
            theme: 'snow',
            readOnly: this.readonly,
        }
        this.quill = new Quill('#editor', options)
        this.quill!.formatLine(0, this.quill!.getLength(), { 'code-block': true });
        this.quill!.on('text-change', () => {
            this.$emit('input', this.quill!.getText().trim())
        })
    }

    mounted() {
        this.createQuillEditor()
    }

    @Watch('readonly')
    onReadOnlyChange() {
        if (!this.quill) {
            return
        }
        this.quill!.enable(!this.readonly)
    }
}

</script>

<style scoped>

code {
    width: 100%;
}

code::before {
    content: "";
}

code::after {
    content: "";
}

</style>
