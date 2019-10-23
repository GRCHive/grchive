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

export default Vue.extend({
    props : {
        nodeId: Number
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
        controlOwner: Object() as User
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.name.length > 0;
        }
    },
    methods: {
        clearForm() {
            this.name = ""
            this.description = ""
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

            this.$emit('do-save')
        }
    }
})

</script>

