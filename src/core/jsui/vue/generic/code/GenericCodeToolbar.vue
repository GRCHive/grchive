<template>
    <v-toolbar flat height="38px">
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

            <slot name="custom-menu"></slot>
        </v-toolbar-items>

        <v-spacer></v-spacer>

        <span v-if="saveInProgress">
            Saving...
            <v-progress-circular indeterminate size="16"></v-progress-circular>
        </span>

        <slot name="custom-status"></slot>
    </v-toolbar>
</template>


<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        saveInProgress: { 
            type: Boolean,
            default: false,
        },
    }
})

@Component
export default class GenericCodeToolbar extends Props {
    save() {
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
}

</script>
