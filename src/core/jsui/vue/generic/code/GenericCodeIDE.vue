<template>
    <div>
        <generic-code-toolbar
            @save="onEvent('save', ...arguments)"
            ref="toolbar"
        >
        </generic-code-toolbar>

        <generic-code-editor
            :value="value"
            :lang="lang"
            :readonly="readonly"
            :full-height="fullHeight"
            @input="onEvent('input', ...arguments)"
            ref="editor"
        >
        </generic-code-editor>
    </div>

</template>


<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GenericCodeEditor, { Props } from './GenericCodeEditor.vue'
import GenericCodeToolbar from './GenericCodeToolbar.vue'

@Component({
    components: {
        GenericCodeToolbar,
        GenericCodeEditor,
    }
})
export default class GenericCodeIDE extends Props {

    $refs!: {
        editor: GenericCodeEditor,
        toolbar: GenericCodeToolbar,
    }

    onEvent(event : string, ...args : any[]) {
        this.$emit(event, ...args)
    }

    mounted() {
        //@ts-ignore
        let ele : HTMLElement = this.$refs.editor.$el
        // Add events here to let toolbar handle input events.
        document.addEventListener('keydown', (e : KeyboardEvent) => {
            if (!document.activeElement) {
                return
            }
            
            // This needs to be here so that the delete doesn't
            // accidentally trigger a hotkey when a dialog is
            // in focus.
            if (!ele.contains(document.activeElement)) {
                return
            }

            this.$refs.toolbar.handleHotkeys(e)
        })

    }
}

</script>
