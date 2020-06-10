<template>
    <div>
        <v-list-item class="px-0">
            <v-list-item-action class="mr-0" style="flex-grow: 1;">
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>
        </v-list-item>
        
        <v-data-table
            :headers="headers"
            :items="data"
            :search="filterText"
        >
        </v-data-table>
    </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'

@Component
export default class JsonArrayViewer extends Vue {
    @Prop({required: true})
    readonly data! : Array<any>

    filterText: string = ""

    get headers() : any[] {
        if (this.data.length == 0) {
            return []
        }

        let keys = Object.keys(this.data[0])
        return keys.map((ele : string) => ({
            text: ele,
            value: ele,
        }))
    }
}

</script>
