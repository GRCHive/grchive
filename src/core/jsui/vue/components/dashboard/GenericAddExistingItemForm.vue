<template>

<v-card>
    <v-card-title>
        Add Existing {{ itemName }}
    </v-card-title>
    <v-divider></v-divider>

    <section class="ma-2">
        <v-list two-line>
            <v-list-item-group multiple v-model="selectedItems">
                <v-list-item v-for="(item, index) in selectableItems"
                             :key="index" class="pa-1" :value="item">
                    <template v-slot:default="{active, toggle}">
                        <v-list-item-action class="ma-1">
                            <v-checkbox :input-value="active"
                                        @true-value="item"
                                        @click="toggle">
                            </v-checkbox>
                        </v-list-item-action>

                        <v-list-item-content>
                            <v-list-item-title>
                                {{ item.Name }}
                            </v-list-item-title>

                            <v-list-item-subtitle>
                                {{ item.Description }}
                            </v-list-item-subtitle>
                        </v-list-item-content>
                    </template>
                </v-list-item>
            </v-list-item-group>
        </v-list>
    </section>
    <v-divider></v-divider>
    <v-card-actions>
        <v-btn color="error" @click="onCancel">
            Cancel
        </v-btn>
        <div class="flex-grow-1"></div>
        <v-btn color="success" @click="onSelect">
            Select
        </v-btn>
    </v-card-actions>

</v-card>
    
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from "../../../ts/vueSetup"
import { contactUsUrl } from "../../../ts/url"
import { newRisk } from "../../../ts/api/apiRisks"

export default Vue.extend({
    props: {
        itemName : String,
        selectableItems : Array
    },
    data: () => ({
        selectedItems: [] as any[]
    }),
    methods: {
        onCancel() {
            this.selectedItems = []
            this.$emit('do-cancel')
        },
        onSelect() {
            this.$emit('do-select', this.selectedItems)
            this.selectedItems = []
        }
    }
})

</script>
