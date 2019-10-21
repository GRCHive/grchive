<template>
    <v-navigation-drawer absolute right :style="clipStyle" ref="attrNavDrawer" :value="showHide">
        <section v-if="enabled" class="ma-1" style="max-height: calc(100% - 48px);">
            <v-form>
                <v-text-field v-model="currentNode.Name"
                      label="Name"
                      :disabled="!canEditAttr"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      v-on:keydown.stop
                ></v-text-field>

                <v-textarea v-model="currentNode.Description"
                            label="Description"
                            filled
                            :disabled="!canEditAttr"
                            v-on:keydown.stop
                ></v-textarea> 

                <v-select v-model="currentNode.NodeTypeId"
                          :items="nodeTypeItems"
                          filled
                          :disabled="!canEditAttr"
                          label="Type">
                </v-select>
            </v-form>

            <v-list-item class="pb-1">
                <template v-if="canEditAttr" v-bind="{saveEdit, cancelEdit}">
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
                    <v-btn color="primary" @click="startEdit">
                        Edit
                    </v-btn>
                </template>
            </v-list-item>

            <v-divider></v-divider>
            <process-flow-input-output-editor :is-input="true"
                                              :node-id="currentNode.Id">
            </process-flow-input-output-editor>

            <v-divider></v-divider>
            <process-flow-input-output-editor :is-input="false"
                                              :node-id="currentNode.Id">
            </process-flow-input-output-editor>
        </section>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import * as rules from "../../../ts/formRules"
import ProcessFlowInputOutputEditor from './ProcessFlowInputOutputEditor.vue'
import { editProcessFlowNode } from '../../../ts/api/apiProcessFlowNodes'
import { contactUsUrl } from '../../../ts/url'
import MetadataStore from '../../../ts/metadata'

export default Vue.extend({
    data : () => ({
        canEditAttr: false,
        cachedData : {} as ProcessFlowNode,
        rules,
    }),
    props: {
        customClipHeight : Number,
        showHide : Boolean
    },
    components: {
        ProcessFlowInputOutputEditor
    },
    methods : {
        startEdit() {
            this.canEditAttr = true
            this.cachedData = {...this.currentNode}
        },
        cancelEdit() {
            this.canEditAttr = false
            VueSetup.store.commit('updateNodePartial', {
                nodeId: this.currentNode.Id,
                node: this.cachedData
            })
        },
        saveEdit() {
            editProcessFlowNode({
                //@ts-ignore
                csrf: this.$root.csrf,
                nodeId: this.currentNode.Id,
                name: this.currentNode.Name,
                description: this.currentNode.Description,
                type: this.currentNode.NodeTypeId
            }).then((resp : TEditProcessFlowNodeOutput) => {
                VueSetup.store.commit('updateNodePartial', {
                    nodeId: resp.data.Id,
                    node: resp.data
                })
                this.canEditAttr = false
            }).catch((err) => {
                console.log(err)
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please reload the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    },
    computed: {
        clipStyle() : any {
            return {
                "height":  "100vh !important",
                "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
                "top" : this.customClipHeight.toString() + "px"
            }
        },
        enabled() : boolean {
            return VueSetup.store.getters.isNodeSelected
        },
        currentNode() : ProcessFlowNode {
            return VueSetup.store.getters.currentNodeInfo
        },
        nodeTypeItems() : any[] {
            let retItems = [] as any[]
            for (let types of MetadataStore.state.nodeTypes) {
                retItems.push({
                    text: types.Name,
                    value: types.Id
                })
            }
            return retItems
        }
    },
})

</script>
