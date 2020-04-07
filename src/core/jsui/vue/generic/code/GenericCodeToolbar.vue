<template>
    <v-toolbar flat height="30px">
        <v-toolbar-items>
            <v-menu offset-y>
                <template v-slot:activator="{ on }">
                    <v-btn text color="accent" v-on="on">
                        File
                        <v-icon small color="accent">mdi-chevron-down</v-icon>
                    </v-btn>
                </template>

                <v-list dense>
                    <v-list-item dense @click="save">
                        <v-list-item-title>
                            Save
                        </v-list-item-title>
                    </v-list-item>
                    <v-divider></v-divider>
                </v-list>
            </v-menu>
        </v-toolbar-items>

        <v-spacer></v-spacer>

        <span v-if="saveInProgress">
            Saving...
            <v-progress-circular indeterminate size="16"></v-progress-circular>
        </span>

    </v-toolbar>
</template>


<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        codeValue : {
            type : String,
            default: "",
        },
        saveInProgress: { 
            type: Boolean,
            default: false,
        },
    }
})

@Component
export default class GenericCodeToolbar extends Props {
    savedValue : string = ""

    save() {
        this.savedValue = this.codeValue
        this.$emit('save')
    }

    handleHotkeys(e : KeyboardEvent) {
        if (e.ctrlKey) {
            if (e.key == 's') {
                this.save()
                e.preventDefault()
            }
        }
    }

    handleUnload(e : Event) {
        if (this.codeValue != this.savedValue) {
            e.preventDefault()
            e.returnValue = false
        }
    }

    mounted() {
        this.savedValue = this.codeValue
        window.addEventListener('beforeunload', this.handleUnload)
    }
}

</script>
