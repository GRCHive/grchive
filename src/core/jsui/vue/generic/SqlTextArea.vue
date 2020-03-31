<template>
    <generic-code-editor
        :value="value"
        @input="onInput"
        :readonly="readonly"
        lang="text/x-pgsql"
    >
    </generic-code-editor>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GenericCodeEditor from './code/GenericCodeEditor.vue'

const Props = Vue.extend({
    props: {
        value : {
            type: String,
            default: ""
        },
        readonly: {
            type: Boolean,
            default: false,
        },
    }
})

@Component({
    components: {
        GenericCodeEditor
    }
})
export default class SqlTextArea extends Props {
    onInput(txt : string) {
        this.$emit('input', txt)
    }
}

</script>
