<template>

<v-card>
    <v-card-title class="pl-3">
        {{ editMode ? "Edit" : "New" }} SQL Request
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(255)]"
                      :readonly="!canEdit"
        ></v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :readonly="!canEdit"
                    hide-details
        ></v-textarea> 

        <user-search-form-component
            class="mt-4"
            label="Requester"
            :user="requestUser"
            readonly
            v-if="!!referenceRequest"
        >
        </user-search-form-component>

        <v-text-field
            :value="requestTime"
            label="Request Time"
            prepend-icon="mdi-calendar"
            readonly
            v-if="!!referenceRequest"
        >
        </v-text-field>

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
            :disabled="!formValid"
            v-if="canEdit"
        >
            Save
        </v-btn>

        <v-btn
            color="success"
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
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import {
    TNewSqlRequestOutput, newSqlRequest,
    TUpdateSqlRequestOutput, updateSqlRequest,
} from '../../../ts/api/apiSqlRequests'
import {
    DbSqlQueryRequest
} from '../../../ts/sql'
import * as rules from '../../../ts/formRules'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import MetadataStore from '../../../ts/metadata'
import { standardFormatTime } from '../../../ts/time'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'

const Props = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        forceQueryId : {
            type : Number,
            default: -1
        },
        referenceRequest: {
            type: Object,
            default: Object() as () => DbSqlQueryRequest | null
        },
    }
})

@Component({
    components: {
        UserSearchFormComponent,
    }
})
export default class CreateNewSqlRequestForm extends Props {
    canEdit: boolean = false
    formValid: boolean = false

    name : string = ""
    description: string = ""
    rules: any = rules

    get queryId() : number {
        if (this.forceQueryId != -1) {
            return this.forceQueryId
        }
        return -1
    }

    get requestUser() : User | null {
        if (!this.referenceRequest) {
            return null
        }
        return MetadataStore.getters.getUser(this.referenceRequest.UploadUserId)
    }

    get requestTime() : string {
        if (!this.referenceRequest) {
            return ""
        }
        return standardFormatTime(this.referenceRequest.UploadTime)
    }

    cancel() {
        this.$emit('do-cancel')
        this.clearForm()

        if (this.editMode) {
            this.canEdit = false
        }
    }

    save() {
        if (!this.formValid) {
            return
        }

        if (this.editMode) {
            this.doEdit()
        } else {
            this.doSave()
        }
    }

    doEdit() {
        updateSqlRequest({
            requestId: this.referenceRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            description: this.description,
        }).then((resp : TUpdateSqlRequestOutput) => {
            this.$emit('do-save', resp.data)
            this.canEdit = false
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    doSave() {
        newSqlRequest({
            queryId: this.queryId,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            description: this.description,
        }).then((resp : TNewSqlRequestOutput) => {
            this.$emit('do-save', resp.data)
            this.clearForm()
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }

    @Watch('referenceRequest')
    clearForm() {
        if (!!this.referenceRequest) {
            this.name = this.referenceRequest.Name
            this.description = this.referenceRequest.Description
        } else {
            this.name = ""
            this.description = ""
        }
    }
}

</script>
