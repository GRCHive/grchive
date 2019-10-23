<template>

<v-card>
    <v-card-title>
        New Control
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name" label="Name" filled :rules="[rules.required, rules.createMaxLength(256)]">
        </v-text-field>

        <v-textarea v-model="description" label="Description" filled>
        </v-textarea> 

        <v-select
            filled
            label="Control Type"
            v-model="controlType"
            :items="controlTypeItems"
            :rules="[rules.required]"
        ></v-select>

        <user-search-form-component
            label="Control Owner"
            v-bind:user.sync="controlOwner"
        ></user-search-form-component>
        <frequency-form-component
            v-bind:isManual.sync="frequencyData.isManual"
            v-bind:freqInterval.sync="frequencyData.freqInterval"
            v-bind:freqType.sync="frequencyData.freqType"
        ></frequency-form-component>

    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
        >
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="success"
            @click="save"
            :disabled="!canSubmit"
        >
            Save
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
import { newControl } from "../../../ts/api/apiControls"
import { contactUsUrl } from "../../../ts/url"

export default Vue.extend({
    props : {
        nodeId: Number,
        riskId: Number
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
            isManual : false,
            freqInterval : 0,
            freqType: 0
        },
        controlType: Object() as ProcessFlowControlType,
        controlOwner: Object() as User
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
            this.name = ""
            this.description = ""
            this.frequencyData.isManual = false
            this.frequencyData.freqInterval = 0
            this.frequencyData.freqType = 0
            this.controlOwner = Object() as User
        },
        cancel() {
            this.$emit('do-cancel')
            this.clearForm()
        },
        save() {
            //@ts-ignore
            if (!this.canSubmit) {
                return;
            }

            newControl(<TNewControlInput>{
                //@ts-ignore
                csrf: this.$root.csrf,
                name: this.name,
                description: this.description,
                controlType: this.controlType.Id,
                frequencyType : this.frequencyData.freqType,
                frequencyInterval : this.frequencyData.freqInterval,
                ownerId : this.controlOwner.Id,
                nodeId: this.nodeId,
                riskId: this.riskId
            }).then((resp : TNewControlOutput) => {
                this.$emit('do-save', resp.data)
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
        refreshDefaultControlType() {
            if (this.controlTypeItems.length > 0) {
                this.controlType = Metadata.state.controlTypes[0]
            } else {
                this.controlType = Object() as ProcessFlowControlType
            }
        }
    },
    watch : {
        controlTypeItems() {
            this.refreshDefaultControlType()
        }
    },
    mounted() {
        this.refreshDefaultControlType()
    }
})

</script>

