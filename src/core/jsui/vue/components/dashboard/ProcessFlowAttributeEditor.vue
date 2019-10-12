<template>
    <v-navigation-drawer absolute right :style="clipStyle" ref="attrNavDrawer" :value="showHide">
        <section v-if="enabled" class="ma-1">
            <v-list-item class="pa-0">
                <template v-if="canEdit" v-bind="{saveEdit, cancelEdit}">
                    <v-btn color="error" @click="cancelEdit">
                        Cancel
                    </v-btn>
                    <div class="flex-grow-1"></div>
                    <v-btn color="success" @click="saveEdit">
                        Save
                    </v-btn>
                </template>

                <template v-else v-bind="startEdit">
                    <div class="flex-grow-1"></div>
                    <v-btn @click="startEdit">
                        Edit
                    </v-btn>
                </template>
            </v-list-item>
            <v-form>
                <v-text-field v-model="currentData.Name"
                      label="Name"
                      :disabled="!canEdit"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                ></v-text-field>

                <v-textarea v-model="currentData.Description"
                            label="Description"
                            filled
                            :disabled="!canEdit">
                </v-textarea> 
            </v-form>
        </section>

        <section v-else>
        </section>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import * as rules from "../../../ts/formRules"

export default Vue.extend({
    props: {
        customClipHeight : Number,
        showHide : Boolean
    },
    data : () => ({
        canEdit : false,
        currentData : {} as ProcessFlowNode,
        rules
    }),
    methods : {
        startEdit() {
            this.canEdit = true
        },
        cancelEdit() {
            this.canEdit = false
        },
        saveEdit() {
            this.canEdit = false
        }
    },
    computed: {
        clipStyle() {
            return {
                //"transform": "translateX(0%)",
                //"width": "256px",
                "height":  "100vh !important",
                //@ts-ignore
                "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
                //@ts-ignore
                "top" : this.customClipHeight.toString() + "px"
            }
        },
        enabled() : boolean {
            return VueSetup.store.getters.isNodeSelected
        }
    },
    watch : {
        enabled(val : boolean) {
            if (val) {
                this.currentData = VueSetup.store.getters.currentNodeInfo
            } else {
                this.currentData = {} as ProcessFlowNode
            }
        }
    }
})

</script>
