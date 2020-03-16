<template>
    <section>
        <v-form>
            <v-text-field v-model="currentNode.Name"
                  label="Name"
                  :readonly="!canEditAttr"
                  filled
                  :rules="[rules.required, rules.createMaxLength(256)]"
                  v-on:keydown.stop
            ></v-text-field>

            <v-textarea v-model="currentNode.Description"
                        label="Description"
                        filled
                        :readonly="!canEditAttr"
                        v-on:keydown.stop
            ></v-textarea> 

            <v-select v-model="currentNode.NodeTypeId"
                      :items="nodeTypeItems"
                      filled
                      :readonly="!canEditAttr"
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

        <v-divider></v-divider>
        <node-linked-risks-editor>
        </node-linked-risks-editor>

        <v-divider></v-divider>
        <node-linked-controls-editor>
        </node-linked-controls-editor>

        <div v-if="canLinkToSystem && linkedSystems != null">
            <v-divider></v-divider>
            <node-linked-systems-editor>
            </node-linked-systems-editor>
        </div>

        <div v-if="canLinkToGL && linkedGL != null">
            <v-divider></v-divider>
            <node-linked-gl-editor>
            </node-linked-gl-editor>
        </div>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import * as rules from "../../../ts/formRules"
import ProcessFlowInputOutputEditor from './ProcessFlowInputOutputEditor.vue'
import { editProcessFlowNode, TEditProcessFlowNodeInput, TEditProcessFlowNodeOutput } from '../../../ts/api/apiProcessFlowNodes'
import { contactUsUrl } from '../../../ts/url'
import { System } from '../../../ts/systems'
import { GeneralLedgerAccount, GeneralLedger } from '../../../ts/generalLedger'
import { PageParamsStore } from '../../../ts/pageParams'
import MetadataStore from '../../../ts/metadata'
import NodeLinkedControlsEditor from './node/NodeLinkedControlsEditor.vue'
import NodeLinkedRisksEditor from './node/NodeLinkedRisksEditor.vue'
import NodeLinkedSystemsEditor from './node/NodeLinkedSystemsEditor.vue'
import NodeLinkedGlEditor from './node/NodeLinkedGlEditor.vue'

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
        ProcessFlowInputOutputEditor,
        NodeLinkedControlsEditor,
        NodeLinkedRisksEditor,
        NodeLinkedSystemsEditor,
        NodeLinkedGlEditor
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
            editProcessFlowNode(<TEditProcessFlowNodeInput>{
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
        },
    },
    computed: {
        clipStyle() : any {
            return {
                "height":  "100vh !important",
                "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
                "top" : this.customClipHeight.toString() + "px"
            }
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
        },
        canLinkToSystem() : boolean {
            return MetadataStore.state.idToNodeTypes[this.currentNode.NodeTypeId].CanLinkToSystem
        },
        linkedSystems() : System[] | null {
            return VueSetup.store.getters.systemsLinkedToNode(this.currentNode.Id)
        },
        canLinkToGL() : boolean {
            return MetadataStore.state.idToNodeTypes[this.currentNode.NodeTypeId].CanLinkToGL
        },
        linkedGL(): GeneralLedgerAccount[] | null {
            let gl : GeneralLedger | null = VueSetup.store.getters.glLinkedToNode(this.currentNode.Id)
            if (!gl) {
                return null
            }
            return gl.changed && gl.listAccounts
        }

    },
})

</script>
