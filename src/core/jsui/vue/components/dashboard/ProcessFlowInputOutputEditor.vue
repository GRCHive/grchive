<template>
    <v-list dense class="pa-0">
        <v-list-item class="pa-0">
            <v-list-item-action class="ma-0">
                <v-btn icon @click="minimize = !minimize">
                    <v-icon small>
                        {{ !minimize ? "mdi-window-minimize" : "mdi-arrow-expand-all" }}
                    </v-icon>
                </v-btn>
            </v-list-item-action>

            <v-subheader class="flex-grow-1 pr-0">
                {{ isInput ? `INPUT` : `OUTPUT` }}
            </v-subheader>
            <v-list-item-action class="ma-0">
                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn icon class="ma-0" v-on="on">
                            <v-icon small>mdi-plus</v-icon>
                        </v-btn>
                    </template>
                    <v-list dense>
                        <v-list-item v-for="(item, index) in dataTypes"
                                     :key="index"
                                     @click="addIO(item)"
                                     dense
                        >
                            <v-list-item-title>
                                {{ item.Name }}
                            </v-list-item-title>
                        </v-list-item>

                    </v-list>
                </v-menu>
            </v-list-item-action>
        </v-list-item>

        <div v-if="!minimize">
            <template v-for="item in listedIO" >
                <v-list-item :key="item.Id" class="body-2 super-dense px-1">
                    <v-list-item-action class="ma-0" v-if="!canEdit(item.Id)">
                        <v-btn small icon class="ma-0" @click="deleteIO($event, item.Id)">
                            <v-icon small>mdi-delete</v-icon>
                        </v-btn>
                    </v-list-item-action>

                    <v-list-item-content class="pa-0 mr-1">
                        <input type="text" 
                               v-model="item.Name"
                               required
                               :disabled="!canEdit(item.Id)"
                               :class="canEdit(item.Id)? `name-edit-style name-style` : `name-style`">
                    </v-list-item-content>
                    <v-list-item-content class="pa-0 mr-1">
                        <select v-model="item.TypeId"
                                :disabled="!canEdit(item.Id)"
                                :class="canEdit(item.Id)? `select-edit-style` : ``">
                            <option v-for="typ in dataTypes"
                                    :key="typ.Id"
                                    :value="typ.Id">
                                {{ typ.Name }}
                            </option>
                        </select>
                    </v-list-item-content>

                    <v-list-item-action class="ma-0" v-if="!canEdit(item.Id)">
                        <v-btn small icon class="ma-0" @click="editIO($event, item)">
                            <v-icon small>mdi-pencil</v-icon>
                        </v-btn>
                    </v-list-item-action>

                    <v-list-item-action class="ma-0" v-if="canEdit(item.Id)">
                        <v-btn small icon class="ma-0" @click="cancelIO($event, item)">
                            <v-icon small>mdi-close</v-icon>
                        </v-btn>
                    </v-list-item-action>

                    <v-list-item-action class="ma-0" v-if="canEdit(item.Id)">
                        <v-btn small icon class="ma-0" @click="saveIO($event, item)">
                            <v-icon small>mdi-check</v-icon>
                        </v-btn>
                    </v-list-item-action>

                    <v-list-item-action class="ma-0" v-if="!canEdit(item.Id)">
                        <v-btn
                            small
                            icon
                            class="ma-0"
                            :disabled="item.IoOrder == maxIoOrder(item.TypeId)"
                            @click="raiseLowerIO(item, 1)"
                        >
                            <v-icon small>mdi-chevron-down</v-icon>
                        </v-btn>
                    </v-list-item-action>

                    <v-list-item-action class="ma-0" v-if="!canEdit(item.Id)">
                        <v-btn
                            small
                            icon
                            class="ma-0"
                            :disabled="item.IoOrder == minIoOrder(item.TypeId)"
                            @click="raiseLowerIO(item, -1)"
                        >
                            <v-icon small>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-action>
                </v-list-item>
                <v-divider
                    inset
                    :key="`divider-${item.Id}`"
                    v-if="item.IoOrder == maxIoOrder(item.TypeId)"
                ></v-divider>
            </template>
        </div>
    </v-list>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import MetadataStore from '../../../ts/metadata'
import { contactUsUrl } from '../../../ts/url'
import { TDeleteProcessFlowIOInput, TDeleteProcessFlowIOOutput, deleteProcessFlowIO } from '../../../ts/api/apiProcessFlowIO'
import { TEditProcessFlowIOInput, TEditProcessFlowIOOutput, editProcessFlowIO } from '../../../ts/api/apiProcessFlowIO'
import { TNewProcessFlowIOInput, TNewProcessFlowIOOutput, newProcessFlowIO } from '../../../ts/api/apiProcessFlowIO'
import { TOrderProcessFlowIOOutput, orderProcessFlowIO } from '../../../ts/api/apiProcessFlowIO'

interface IOEditState {
    canEdit: boolean,
    cachedData: ProcessFlowInputOutput
}

export default Vue.extend({
    props: {
        isInput: Boolean,
        nodeId: Number
    },
    data : () => ({
        ioEditState: Object() as Record<number, IOEditState>,
        minimize: false
    }),
    computed: {
        listedIO() : ProcessFlowInputOutput[] {
            if (this.isInput) {
                return VueSetup.store.getters.nodeInfo(this.nodeId).Inputs
            } else {
                return VueSetup.store.getters.nodeInfo(this.nodeId).Outputs
            }
        },
        canEdit() {
            return (ioId : number) => {
                if (!(ioId in this.ioEditState)) {
                    return false
                }
                return this.ioEditState[ioId].canEdit
            }
        },
        dataTypes(): ProcessFlowIOType[] {
            return MetadataStore.state.ioTypes
        },
        minIoOrder() : (a0 : number) => number { 
            return (typeId: number) : number => {
                return Math.min(...this.listedIO.filter((ele : ProcessFlowInputOutput) => ele.TypeId == typeId)
                    .map((ele : ProcessFlowInputOutput) => ele.IoOrder))
            }
        },
        maxIoOrder() : (a0 : number) => number { 
            return (typeId: number) : number => {
                return Math.max(...this.listedIO.filter((ele : ProcessFlowInputOutput) => ele.TypeId == typeId)
                    .map((ele : ProcessFlowInputOutput) => ele.IoOrder))
            }
        }
    },
    methods: {
        addIO(type : ProcessFlowIOType) {
            let name: string  = ""
            if (this.isInput) {
                name = "Input " + this.listedIO.length.toString()
            } else {
                name = "Output " + this.listedIO.length.toString()
            }

            newProcessFlowIO(<TNewProcessFlowIOInput>{
                nodeId: this.nodeId,
                typeId: type.Id,
                isInput: this.isInput,
                name,
            }).then((resp : TNewProcessFlowIOOutput) => {
                if (this.isInput) {
                    VueSetup.store.commit('addNodeInput', {nodeId: this.nodeId, input: resp.data})
                } else {
                    VueSetup.store.commit('addNodeOutput', {nodeId: this.nodeId, output: resp.data})
                }
            }).catch((err : any) => {
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
        deleteIO(e: MouseEvent, ioId : number) {
            // Maybe we want to confirm with the user first?
            deleteProcessFlowIO({
                ioId,
                isInput: this.isInput
            }).then((resp : TDeleteProcessFlowIOOutput) => {
                if (this.isInput) {
                    VueSetup.store.dispatch('deleteBatchNodeInput', {nodeId: this.nodeId, inputs: [ioId]})
                } else {
                    VueSetup.store.dispatch('deleteBatchNodeOutput', {nodeId: this.nodeId, outputs: [ioId]})
                }
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
        raiseLowerIO(io : ProcessFlowInputOutput, dir : number) {
            orderProcessFlowIO({
                ioId: io.Id,
                isInput: this.isInput,
                direction: dir,
            }).then((resp : TOrderProcessFlowIOOutput) => {
                VueSetup.store.commit('updateNodeInputOutput', {
                    nodeId: this.nodeId,
                    io: resp.data.This,
                    isInput: this.isInput
                })
                VueSetup.store.commit('updateNodeInputOutput', {
                    nodeId: this.nodeId,
                    io: resp.data.Other,
                    isInput: this.isInput
                })
            }).catch((err: any) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please reload the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        editIO(e: MouseEvent, io : ProcessFlowInputOutput) {
            this.ioEditState[io.Id].canEdit = true
            this.ioEditState[io.Id].cachedData = {...io}
        },
        saveIO(e: MouseEvent, io : ProcessFlowInputOutput) {
            editProcessFlowIO({
                ioId: io.Id,
                isInput: this.isInput,
                name: io.Name ,
                type: io.TypeId
            }).then((resp : TEditProcessFlowIOOutput) => {
                this.ioEditState[io.Id].canEdit = false
                VueSetup.store.commit('updateNodeInputOutput', {
                    nodeId: this.nodeId,
                    io: resp.data,
                    isInput: this.isInput
                })
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
        cancelIO(e: MouseEvent, io : ProcessFlowInputOutput) {
            this.ioEditState[io.Id].canEdit = false
            VueSetup.store.commit('updateNodeInputOutput', {
                nodeId: this.nodeId,
                io: this.ioEditState[io.Id].cachedData,
                isInput: this.isInput
            })
        },
        updateEditState(val : ProcessFlowInputOutput[]) {
            for (let d of val) {
                if (!(d.Id in this.ioEditState)) {
                    Vue.set(this.ioEditState, d.Id, <IOEditState>{
                        canEdit: false
                    })
                }
            }
        }
    },
    watch : {
        listedIO(val : ProcessFlowInputOutput[]) {
            this.updateEditState(val)
        }
    },
    mounted() {
        this.updateEditState(this.listedIO)
    }
})

</script>

<style scoped>

.super-dense {
    min-height: 30px !important;
}

.name-style {
    flex: initial !important;
    width: 100%;
}

.name-edit-style {
    border-style: solid !important;
    border-color: black;
}

.select-edit-style {
    appearance: menulist !important;
    -moz-appearance: menulist !important;
    -webkit-appearance: menulist !important;
    border-style: unset !important;
    background-color: unset !important;
}

</style>
