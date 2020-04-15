<template>
    <div ref="parent" :style="parentContainerStyle">
        <div :style="contentContainerStyle">
            <v-list-item>
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Data Sources
                    </v-list-item-title>
                </v-list-item-content>

                <v-list-item-action>
                    <v-dialog v-model="showHideLinkDataSource"
                              persistent
                              max-width="40%"
                    >
                        <template v-slot:activator="{ on }">
                            <v-btn 
                                color="primary"
                                icon
                                v-on="on"
                            >
                                <v-icon>
                                    mdi-plus
                                </v-icon>
                            </v-btn>
                        </template>

                        <v-card>
                            <client-data-table-with-controls
                                class="ma-4"
                                v-model="stagedClientDataForLink"
                                :exclude="linkedClientData"
                                disable-new
                                disable-delete
                                enable-select
                            >
                            </client-data-table-with-controls>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="showHideLinkDataSource = false"
                                >
                                    Cancel
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-btn
                                    color="success"
                                    @click="doLinkClientData"
                                >
                                    Link
                                </v-btn>
                            </v-card-actions>
                        </v-card>
                    </v-dialog>
                </v-list-item-action>
            </v-list-item>
            <v-divider inset></v-divider>

            <client-data-table
                :resources="linkedClientData"
                use-crud-delete
                @delete="unlinkClientData"
            >
            </client-data-table>

            <v-divider></v-divider>

            <v-list-item>
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Script Parameters
                    </v-list-item-title>
                </v-list-item-content>

                <v-list-item-action>
                    <v-btn 
                        color="primary"
                        icon
                        @click="addEmptyParameterType"
                    >
                        <v-icon>
                            mdi-plus
                        </v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
            <v-divider inset></v-divider>

            <template v-for="(item, index) in scriptParameterTypes">
                <div class="paramContainer">
                    <param-type-component
                        :key="`type-${index}`"
                        v-model="scriptParameterTypes[index]"
                        :rules="[rules.required]"
                    >
                    </param-type-component>

                    <v-btn 
                        color="error"
                        icon
                        @click="removeParameterType(index)"
                    >
                        <v-icon>
                            mdi-delete
                        </v-icon>
                    </v-btn>
                </div>
            </template>

            <v-divider></v-divider>

            <v-list-item>
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Parameter Values
                    </v-list-item-title>
                </v-list-item-content>
            </v-list-item>
            <v-divider inset></v-divider>

            <param-value-component
                v-for="(item, index) in scriptParameterTypes"
                :key="`value-${index}`"
                :param="item"
            >
            </param-value-component>
        </div>

        <v-divider></v-divider>
        <v-list-item>
            <v-list-item-action>
                <v-btn 
                    color="warning"
                >
                    Revert
                </v-btn>
            </v-list-item-action>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn 
                    color="success"
                >
                    Run
                </v-btn>
            </v-list-item-action>
        </v-list-item>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import ClientDataTable from '../ClientDataTable.vue'
import ParamTypeComponent from './parameters/ParamTypeComponent.vue'
import ParamValueComponent from './parameters/ParamValueComponent.vue'
import ClientDataTableWithControls from '../resources/ClientDataTableWithControls.vue'
import { ClientData, FullClientDataWithLink } from '../../../ts/clientData'
import {
    CodeParamType
} from '../../../ts/code'
import * as rules from '../../../ts/formRules'

const Props = Vue.extend({
    props: {
        linkedClientData: {
            type: Array as () => Array<FullClientDataWithLink>,
            default: () => [],
        },
        scriptParameterTypes: {
            type: Array as () => Array<CodeParamType | null>,
            default: () => [],
        },
    }
})

@Component({
    components: {
        ClientDataTable,
        ParamTypeComponent,
        ParamValueComponent,
        ClientDataTableWithControls,
    }
})
export default class ScriptParamsEditor extends Props {
    rules : any = rules
    parentBb : DOMRect | null = null

    showHideLinkDataSource : boolean = false
    stagedClientDataForLink : FullClientDataWithLink[] = []

    $refs! : {
        parent : HTMLElement
    }

    doLinkClientData() {
        this.linkedClientData.unshift(...this.stagedClientDataForLink)
        this.$emit('update:linkedClientData', this.linkedClientData)

        this.stagedClientDataForLink = []
        this.showHideLinkDataSource = false
    }

    unlinkClientData(d : FullClientDataWithLink) {
        let idx = this.linkedClientData.findIndex((ele : FullClientDataWithLink) => ele.Data.Id == d.Data.Id)
        if (idx == -1) {
            return
        }
        this.linkedClientData.splice(idx, 1)
        this.$emit('update:linkedClientData', this.linkedClientData)
    }

    addEmptyParameterType() {
        this.scriptParameterTypes.push(null)
        this.$emit('update:scriptParameterTypes', this.scriptParameterTypes)
    }

    removeParameterType(idx : number) {
        this.scriptParameterTypes.splice(idx, 1)
        this.$emit('update:scriptParameterTypes', this.scriptParameterTypes)
    }

    get parentContainerStyle() : any {
        if (!this.parentBb) {
            Vue.nextTick(() => {
                this.parentBb = <DOMRect>this.$refs.parent.getBoundingClientRect()
            })

            return {}
        }

        let ht = `calc(100vh - ${this.parentBb.top}px)`
        return {
            'height': ht,
            'max-height': ht,
        }
    }

    get contentContainerStyle() : any {
        if (!this.parentBb) {
            return {}
        }

        let ht = `calc(100vh - ${this.parentBb.top}px - 60px)`
        return {
            'height': ht,
            'max-height': ht,
            'overflow': 'auto',
        }
    }
}

</script>

<style scoped>

.paramContainer {
    display: flex;
    align-items: center;
}

</style>
