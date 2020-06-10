<template>
    <div>
        <v-list-item class="px-0">
            <v-spacer></v-spacer>
            <v-list-item-action>
                <v-dialog
                    v-model="showHideRaw"
                    max-width="60%"
                    persistent
                >
                    <template v-slot:activator="{on}">
                        <v-btn
                            color="secondary"
                            v-on="on"
                        >
                            View Raw
                        </v-btn>
                    </template>

                    <v-card>
                        <v-card-title>Raw Data</v-card-title>
                        <v-divider></v-divider>

                        <generic-code-editor
                            :value="raw"
                            :lang="application/json"
                            readonly
                            fixed-height="70vh"
                        >
                        </generic-code-editor>

                        <v-card-actions>
                            <v-btn
                                color="error"
                                @click="showHideRaw = false"
                            >
                                Close
                            </v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

        <v-expansion-panels>
            <v-expansion-panel
                v-for="(key, idx) in keys"
                :key="idx"
            >
                <v-expansion-panel-header class="font-weight-bold">{{key}}</v-expansion-panel-header>
                <v-expansion-panel-content>
                    <json-array-viewer
                        v-if="isDataArray(key)"
                        :data="get(key)"
                    >
                    </json-array-viewer>
                    <generic-json-viewer
                        v-else-if="isDataObject(key)"
                        :data="get(key)"
                    >
                    </generic-json-viewer>
                    <span v-else>{{ get(key) }}</span>
                </v-expansion-panel-content>
            </v-expansion-panel>
        </v-expansion-panels>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import JsonArrayViewer from './JsonArrayViewer.vue'
import GenericCodeEditor from '../code/GenericCodeEditor.vue'

@Component({
    components: {
        JsonArrayViewer,
        GenericCodeEditor,
    }
})
export default class GenericJsonViewer extends Vue {
    @Prop({required: true})
    readonly data! : { [index : string] : any }

    showHideRaw : boolean = false

    get raw() : string {
        return JSON.stringify(this.data, null, 4)
    }

    get keys() : string[] {
        return Object.keys(this.data)
    }

    get(key : string) : any {
        return this.data[key]
    }

    isDataArray(key : string) : boolean {
        return Array.isArray(this.get(key))
    }

    isDataObject(key : string) : boolean {
        return Object.keys(this.get(key)).length > 0
    }
}

</script>
