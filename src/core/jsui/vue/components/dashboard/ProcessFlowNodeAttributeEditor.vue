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

        <div v-if="canLinkToSystem && linkedSystems != null">
            <v-divider></v-divider>
            <v-list dense class="pa-0">
                <v-list-item class="pa-0">
                    <v-subheader class="flex-grow-1 pr-0">
                        LINKED SYSTEMS
                    </v-subheader>

                    <v-list-item-action class="ma-0">
                        <v-dialog persistent max-width="40%" v-model="showLinkSystem">
                            <template v-slot:activator="{ on }">
                                <v-btn
                                    icon
                                    v-on="on"
                                >
                                    <v-icon small>
                                        mdi-plus
                                    </v-icon>
                                </v-btn>
                            </template>

                            <v-card>
                                <v-card-title>
                                    Link System
                                </v-card-title>

                                <system-search-form-component
                                    v-model="systemsToLink"
                                >
                                </system-search-form-component>

                                <v-card-actions>
                                    <v-btn
                                        color="error"
                                        @click="cancelSystemLink"
                                    >
                                        Cancel
                                    </v-btn>
                                    <v-spacer></v-spacer>
                                    <v-btn
                                        color="success"
                                        @click="saveSystemLink"
                                        :disabled="systemsToLink.length == 0"
                                    >
                                        Link
                                    </v-btn>
                                </v-card-actions>
                            </v-card>
                        </v-dialog>
                    </v-list-item-action>
                </v-list-item>

                <systems-table
                    :resources="linkedSystems"
                    use-crud-delete
                    @delete="deleteLinkedSystem"
                >
                </systems-table>
            </v-list>
        </div>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import * as rules from "../../../ts/formRules"
import ProcessFlowInputOutputEditor from './ProcessFlowInputOutputEditor.vue'
import { editProcessFlowNode, TEditProcessFlowNodeInput, TEditProcessFlowNodeOutput } from '../../../ts/api/apiProcessFlowNodes'
import { 
    newNodeSystemLink,
    deleteNodeSystemLink
} from '../../../ts/api/apiNodeSystemLinks'
import { contactUsUrl } from '../../../ts/url'
import { System } from '../../../ts/systems'
import { PageParamsStore } from '../../../ts/pageParams'
import MetadataStore from '../../../ts/metadata'
import SystemSearchFormComponent from '../../generic/SystemSearchFormComponent.vue'
import SystemsTable from '../../generic/SystemsTable.vue'

export default Vue.extend({
    data : () => ({
        canEditAttr: false,
        cachedData : {} as ProcessFlowNode,
        systemsToLink: [] as System[],
        showLinkSystem: false,
        rules,
    }),
    props: {
        customClipHeight : Number,
        showHide : Boolean
    },
    components: {
        ProcessFlowInputOutputEditor,
        SystemSearchFormComponent,
        SystemsTable
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
        cancelSystemLink() {
            this.systemsToLink = []
            this.showLinkSystem = false
        },
        saveSystemLink() {
            if (this.systemsToLink.length == 0) {
                return
            }

            let nodeId : number = this.currentNode.Id
            let system : System = this.systemsToLink[0]
            newNodeSystemLink({
                nodeId: nodeId,
                systemId: system.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then(() => {
                VueSetup.store.commit('addNodeSystemLink', {
                    nodeId: nodeId,
                    system: system,
                })
                this.systemsToLink = []
                this.showLinkSystem = false
            }).catch((err : any) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please reload the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        deleteLinkedSystem(system : System) {
            let id : number = this.currentNode.Id
            deleteNodeSystemLink({
                nodeId: id,
                systemId: system.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then(() => {
                VueSetup.store.commit('deleteNodeSystemLink', {
                    nodeId: id,
                    systemId: system.Id,
                })
            }).catch((err : any) => {
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
    },
})

</script>
