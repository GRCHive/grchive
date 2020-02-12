<template>
    <div>
        <v-list>
            <v-list-item
                class="pa-0"
                v-for="(k, i) in keys"
                :key="i"
            >
                <v-list-item-content class="mr-1">
                    <v-text-field
                        :value="k"
                        @input="modifyKey(k, arguments[0])"
                        filled
                        label="Key"
                        hide-details
                        dense
                    >
                    </v-text-field>
                </v-list-item-content>

                <v-list-item-content class="ml-1">
                    <v-text-field
                        :value="values[i]"
                        @input="modifyValue(k, arguments[0])"
                        filled
                        label="Value"
                        hide-details
                        dense
                    >
                    </v-text-field>

                </v-list-item-content>
            </v-list-item>

            <v-list-item
                class="pa-0"
            >
                <v-btn
                    @click="addEntry"
                    color="primary"
                    outlined
                    block
                >
                    Add Entry
                </v-btn>
            </v-list-item>
        </v-list>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'

const Props = Vue.extend({
    props: {
        value: Object,
        default: Object() as () => Record<string,string>
    }
})

@Component
export default class StringDictFormComponent extends Props {
    get keys() : string[] {
        return Object.keys(this.value)
    }

    get values() : string[] {
        return this.keys.map((ele : string) => this.value[ele])
    }

    modifyKey(k : string, newKey : string) {
        let newObj = this.value
        Vue.set(newObj, newKey, newObj[k])
        delete newObj[k]
        this.$emit('input', newObj)
    }

    modifyValue(k : string, newValue : string) {
        let newObj = this.value
        Vue.set(newObj, k, newValue)
        this.$emit('input', newObj)
    }

    addEntry() {
        let newObj = this.value
        Vue.set(newObj, '', '')
        this.$emit('input', newObj)
    }
}

</script>
