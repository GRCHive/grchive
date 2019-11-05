<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Risk
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :disabled="!canEdit">
        </v-textarea> 

    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit"
        >
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="success"
            @click="save"
            :disabled="!canSubmit"
            v-if="canEdit"
        >
            Save
        </v-btn>

        <v-btn
            color="success"
            @click="edit"
            v-if="!canEdit"
        >
            Edit
        </v-btn>
    </v-card-actions>
</v-card>
    
</template>

<script lang="ts">

import Vue from 'vue'
import * as rules from "../../../ts/formRules"
import { contactUsUrl } from "../../../ts/url"
import { editRisk, TEditRiskInput, TEditRiskOutput } from "../../../ts/api/apiRisks"
import { newRisk, TNewRiskInput, TNewRiskOutput } from "../../../ts/api/apiRisks"
import { getCurrentCSRF } from '../../../ts/csrf'

export default Vue.extend({
    props : {
        nodeId: Number,
        editMode: {
            type: Boolean,
            default: false
        },
        defaultName: {
            type: String,
            default: ""
        },
        defaultDescription: {
            type: String,
            default: ""
        },
        riskId: {
            type: Number,
            default: -1
        },
        stagedEdits: {
            type: Boolean,
            default: false
        }
    },
    data: () => ({
        name: "",
        description: "",
        rules,
        formValid: false,
        canEdit: false
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.name.length > 0;
        }
    },
    methods: {
        clearForm() {
            this.name = this.defaultName
            this.description = this.defaultDescription
        },
        edit() {
            this.canEdit = true
        },
        cancel() {
            if (this.stagedEdits) {
                this.canEdit = false
            }
            this.$emit('do-cancel')
            this.clearForm()
        },
        save() {
            //@ts-ignore
            if (!this.canSubmit) {
                return;
            }

            if (this.stagedEdits) {
                this.canEdit = false
            }

            if (this.editMode) {
                this.doEdit()
            } else {
                this.doSave()
            }
        },
        onSuccess(risk : ProcessFlowRisk) {
            this.clearForm()
            this.$emit('do-save', risk)
        },
        onError( err: any) {
            if (!!err.response && err.response.data.IsDuplicate) {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "A risk with this name exists already. Pick another name.",
                    false,
                    "",
                    contactUsUrl,
                    true);
            } else {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            }
        },
        doSave() {
            newRisk(<TNewRiskInput>{
                csrf : getCurrentCSRF(),
                name : this.name,
                description: this.description,
                nodeId: this.nodeId,
                //@ts-ignore
                orgName: this.$root.orgGroupId
            }).then((resp : TNewRiskOutput) => {
                this.onSuccess(resp.data)
            }).catch((err : any) => {
                this.onError(err)
            })
        },
        doEdit() {
            editRisk(<TEditRiskInput>{
                csrf : getCurrentCSRF(),
                name : this.name,
                description: this.description,
                riskId: this.riskId,
                //@ts-ignore
                orgName: this.$root.orgGroupId
            }).then((resp : TEditRiskOutput) => {
                this.onSuccess(resp.data)
            }).catch((err : any) => {
                this.onError(err)
            })
        }
    },
    mounted() {
        this.canEdit = (!this.stagedEdits || !this.editMode)
    }
})

</script>
