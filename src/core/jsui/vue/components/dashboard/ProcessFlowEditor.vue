<template>
    <section id="content" >
        <v-list-item two-line id="header">
            <v-list-item-content v-if="!editMode" class="mr-2">
                <v-list-item-title>
                    {{ basicData.Name }}
                    <v-btn icon @click="onExpandDescription">
                        <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                        <v-icon small v-else>mdi-chevron-up</v-icon>
                    </v-btn>
                </v-list-item-title>
                <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                    {{ basicData.Description }}
                </v-list-item-subtitle>
            </v-list-item-content>

            <v-list-item-content v-else class="mr-2">
                <v-form ref="editForm" v-model="formValid">
                    <v-text-field v-model="editName" label="Name" filled :rules="[rules.required, rules.createMaxLength(256)]" dense class="title"
                                  v-on:keydown.stop
                    ></v-text-field>

                    <v-textarea v-model="editDescription" label="Description" filled dense class="subtitle-1"
                                v-on:keydown.stop
                    ></v-textarea> 
                </v-form>
            </v-list-item-content>

            <v-list-item-content>
                <v-list-item-subtitle>
                    <span class="font-weight-bold">Created Time: </span>
                    <span>{{ standardFormatTime(basicData.CreationTime) }}</span>
                </v-list-item-subtitle>

                <v-list-item-subtitle>
                    <span class="font-weight-bold">Last Updated Time: </span>
                    <span>{{ standardFormatTime(basicData.LastUpdatedTime) }}</span>
                </v-list-item-subtitle>

            </v-list-item-content>

            <v-list-item-action v-if="!editMode">
                <v-btn color="primary" @click="onEdit">
                    Edit
                </v-btn>
            </v-list-item-action>


            <v-list-item-action v-if="editMode">
                <v-btn color="error" @click="cancelEdit">
                    Cancel
                </v-btn>
            </v-list-item-action>


            <v-list-item-action v-if="editMode">
                <v-btn color="success" @click="saveEdit" :disabled="!canSubmit">
                    Save
                </v-btn>
            </v-list-item-action>

        </v-list-item>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import VueRouter from 'vue-router'
import * as rules from "../../../ts/formRules"
import { standardFormatTime } from '../../../ts/time'
import { postFormUrlEncoded } from "../../../ts/http"
import { contactUsUrl, createUpdateProcessFlowApiUrl } from "../../../ts/url"

import { getCurrentCSRF } from '../../../ts/csrf'
interface ResponseData {
    data : ProcessFlowBasicData
}

export default Vue.extend({
    data : () => ({
        editMode : false,
        editName : "",
        editDescription: "",
        rules,
        formValid: false,
        expandDescription: false
    }),
    computed: {
        basicData() : ProcessFlowBasicData {
            if (VueSetup.store.state.currentProcessFlowIndex >= VueSetup.store.state.allProcessFlowBasicData.length) {
                return {
                    Id: 0,
                    Name : "",
                    Description: "",
                    CreationTime: new Date(0),
                    LastUpdatedTime: new Date(0)
                }
            }

            return VueSetup.store.state.allProcessFlowBasicData[VueSetup.store.state.currentProcessFlowIndex] 
        },
        canSubmit() : boolean {
            return this.formValid && this.editName.length > 0;
        }
    },
    methods: {
        getCurrentCSRF,
        onExpandDescription() {
            this.expandDescription = !this.expandDescription
            this.$emit('on-change')
        },
        standardFormatTime,
        cancelEdit() {
            // Note that this is valid to call when we switch to another process flow
            // since it'll just pull the data from the Vuex store which is accurate.
            this.editMode = false
            this.editName = this.basicData.Name
            this.editDescription = this.basicData.Description
            this.$emit('on-change')
        },
        saveEdit() {
            //@ts-ignore
            if (!this.canSubmit || !this.$refs.editForm.validate()) {
                return
            }

            //@ts-ignore
            postFormUrlEncoded<ResponseData>(createUpdateProcessFlowApiUrl(this.basicData.Id), {
                name: this.editName,
                description: this.editDescription,
                csrf: this.getCurrentCSRF()
            }).then((resp : ResponseData) => {
                VueSetup.store.commit(
                    "setIndividualProcessFlowBasicData", 
                    {
                        index: VueSetup.store.state.currentProcessFlowIndex,
                        data: resp.data
                    })
                this.editMode = false
                this.$emit('on-change')
            }).catch((err) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        onEdit() {
            this.editMode = true
            this.$emit('on-change')
        }
    },
    mounted() {
        VueSetup.store.watch(
            (state, getters) => {
                return state.currentProcessFlowIndex
            },
            () => {
                this.cancelEdit()
            }
        )

        VueSetup.store.watch(
            (state, getters) => {
                return state.allProcessFlowBasicData
            },
            () => {
                this.cancelEdit()
            }
        )
    },
})
</script>

<style scoped>

/* For whatever reason, the 'selectable' prop doesn't work on v-list-item */
#header {
    user-select: text;
}

</style>
