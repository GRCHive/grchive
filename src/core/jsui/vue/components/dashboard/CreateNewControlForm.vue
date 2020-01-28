<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Control
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-textarea v-model="description" label="Description" filled
                    :disabled="!canEdit">
        </v-textarea> 

        <v-select
            filled
            label="Control Type"
            v-model="controlType"
            :items="controlTypeItems"
            :rules="[rules.required]"
            :disabled="!canEdit"
        ></v-select>

        <user-search-form-component
            label="Control Owner"
            v-bind:user.sync="controlOwner"
            :disabled="!canEdit"
        ></user-search-form-component>
        <frequency-form-component
            v-bind:freqInterval.sync="frequencyData.freqInterval"
            v-bind:freqType.sync="frequencyData.freqType"
            :disabled="!canEdit"
        ></frequency-form-component>

    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit || dialogMode"
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
            color="primary"
            @click="canEdit = true"
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
import FrequencyFormComponent from "../../generic/FrequencyFormComponent.vue"
import UserSearchFormComponent from "../../generic/UserSearchFormComponent.vue"
import Metadata from "../../../ts/metadata"
import { lazyGetUserFromId, lazyGetControlTypeFromId } from '../../../ts/metadataUtils'
import { newControl, 
         editControl,
         TNewControlInput,
         TNewControlOutput,
         TEditControlInput,
         TEditControlOutput } from "../../../ts/api/apiControls"
import { contactUsUrl } from "../../../ts/url"
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    props : {
        nodeId: Number,
        riskId: Number,
        editMode: {
            type: Boolean,
            default: false
        },
        control: {
            type: Object as () => ProcessFlowControl,
            default: () => Object() as ProcessFlowControl
        },
        stagedEdits: {
            type: Boolean,
            default: false
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
    },
    components: {
        FrequencyFormComponent,
        UserSearchFormComponent
    },
    data: () => ({
        name: "",
        description: "",
        rules,
        formValid: false,
        frequencyData : {
            freqInterval : 0,
            freqType: 0
        },
        controlType: Object() as ProcessFlowControlType,
        controlOwner: Object() as User,
        canEdit: true
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.name.length > 0;
        },
        controlTypeItems() : Object[] {
            let retArr = [] as Object[]
            for (let typ of Metadata.state.controlTypes) {
                retArr.push({
                    text: typ.Name,
                    value: typ
                })
            }
            return retArr
        }
    },
    methods: {
        clearForm() {
            if (this.editMode) {
                let control : ProcessFlowControl = this.control
                this.name = control.Name
                this.description = control.Description
                this.frequencyData.freqType = control.FrequencyType
                this.frequencyData.freqInterval = control.FrequencyInterval
                lazyGetUserFromId(control.OwnerId).then((user : User) => {
                    this.controlOwner = user
                })

                lazyGetControlTypeFromId(control.ControlTypeId).then((typ : ProcessFlowControlType) => {
                    this.controlType = typ
                })
            } else {
                this.name = ""
                this.description = ""
                this.frequencyData.freqInterval = 0
                this.frequencyData.freqType = 0
                this.controlOwner = Object() as User
                this.refreshDefaultControlType()
            }
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
        refreshDefaultControlType() {
            if (this.controlTypeItems.length > 0) {
                this.controlType = Metadata.state.controlTypes[0]
            } else {
                this.controlType = Object() as ProcessFlowControlType
            }
        },
        onSuccess(control : ProcessFlowControl) {
            this.$emit('do-save', control, this.riskId)
        },
        onError(err : any) {
            if (!!err.response && err.response.data.IsDuplicate) {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "A control with this name exists already. Pick another name.",
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
            newControl(<TNewControlInput>{
                name: this.name,
                description: this.description,
                controlType: this.controlType.Id,
                frequencyType : this.frequencyData.freqType,
                frequencyInterval : this.frequencyData.freqInterval,
                ownerId : !!this.controlOwner ? this.controlOwner.Id : undefined,
                nodeId: this.nodeId,
                riskId: this.riskId,
                orgName: PageParamsStore.state.organization!.OktaGroupName
            }).then((resp : TNewControlOutput) => {
                this.onSuccess(resp.data)
            }).catch((err : any) => {
                this.onError(err)
            })
        },
        doEdit() {
            editControl(<TEditControlInput>{
                name: this.name,
                description: this.description,
                controlType: this.controlType.Id,
                frequencyType : this.frequencyData.freqType,
                frequencyInterval : this.frequencyData.freqInterval,
                ownerId : !!this.controlOwner ? this.controlOwner.Id : undefined,
                nodeId: this.nodeId,
                riskId: this.riskId,
                controlId: this.control.Id,
                orgName: PageParamsStore.state.organization!.OktaGroupName
            }).then((resp : TEditControlOutput) => {
                this.onSuccess(resp.data)
            }).catch((err : any) => {
                this.onError(err)
            })
        }
    },
    watch : {
        controlTypeItems() {
            this.refreshDefaultControlType()
        }
    },
    mounted() {
        this.canEdit = (!this.stagedEdits || !this.editMode)
        this.refreshDefaultControlType()
    }
})

</script>

