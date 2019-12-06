<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Category: {{ currentCat.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`" >
                        {{ currentCat.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>

            <v-divider></v-divider>

            <v-container fluid>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { ControlDocumentationCategory } from '../../../ts/controls'

@Component
export default class FullEditDocumentationCategoryComponent extends Vue {
    ready: boolean = false
    expandDescription : boolean = false
    currentCat : ControlDocumentationCategory | null = null
}

</script>
