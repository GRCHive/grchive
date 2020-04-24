<template>
    <v-form v-model="formValid">
        <v-text-field
            :value="value.Name"
            @input="handleChange('Name', arguments[0])"
            label="Name"
            filled
            :rules="[rules.required]"
            :readonly="readonly"
        ></v-text-field>

        <v-textarea
            :value="value.Description"
            @input="handleChange('Description', arguments[0])"
            label="Description"
            filled
            :readonly="readonly"
            hide-details
        ></v-textarea> 

        <user-search-form-component
            class="mt-4"
            label="Requester"
            :user="requestUser"
            readonly
        >
        </user-search-form-component>

        <user-search-form-component
            class="mt-4"
            label="Assignee"
            :user.sync="assignee"
            :readonly="readonly"
        >
        </user-search-form-component>

        <date-time-picker-form-component
            :value="value.DueDate"
            @input="handleChange('DueDate', arguments[0])"
            label="Due Date"
            :readonly="readonly"
            clearable
        >
        </date-time-picker-form-component>

        <v-text-field
            :value="uploadTimeStr"
            label="Request Time"
            prepend-icon="mdi-calendar"
            readonly
        >
        </v-text-field>
    </v-form>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { GenericRequest, cleanGenericRequestFromJson } from '../../../ts/requests'
import * as rules from '../../../ts/formRules'
import MetadataStore from '../../../ts/metadata'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'
import DateTimePickerFormComponent from '../../generic/DateTimePickerFormComponent.vue'
import { standardFormatTime } from '../../../ts/time'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => Object() as GenericRequest
        },
        valid: {
            type: Boolean,
            default: false,
        },
        readonly : {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        UserSearchFormComponent,
        DateTimePickerFormComponent,
    }
})
export default class CreateNewGenericRequestForm extends Props {
    refState : GenericRequest = Object() 
    formValid : boolean = false
    rules : any = rules

    assignee : User | null = null

    handleChange(field : string, val : any) {
        let newVal : any = JSON.parse(JSON.stringify(this.value))
        cleanGenericRequestFromJson(newVal)
        newVal[field] = val
        this.$emit('input', newVal)
    }

    get isAssigneeDifferentFromValue() : boolean {
        if (!this.assignee && !this.value.Assignee) {
            return false
        }

        return (!this.assignee && !!this.value.Assignee) ||
            (!!this.assignee && !this.value.Assignee) ||
            (this.assignee!.Id != this.value.Assignee!)
    }

    @Watch('assignee')
    handleAssigneeChange() {
        if (!this.isAssigneeDifferentFromValue) {
            return
        }

        if (!!this.assignee) {
            this.handleChange('Assignee', this.assignee.Id)
        } else {
            this.handleChange('Assignee', null)
        }
    }

    @Watch('value')
    resetAssignee() {
        if (!this.isAssigneeDifferentFromValue) {
            return
        }

        this.assignee = MetadataStore.getters.getUser(this.value.Assignee)
    }

    get requestUser() : User | null {
        return MetadataStore.getters.getUser(this.value.UploadUserId)
    }

    @Watch('formValid')
    updateFormValid() {
        this.$emit('update:valid', this.formValid)
    }

    mounted() {
        this.refState = JSON.parse(JSON.stringify(this.value))
        cleanGenericRequestFromJson(this.refState)
        this.resetAssignee()
    }

    resetToRefState() {
        this.$emit('input', this.refState)
    }

    saveRefState() {
        this.refState = JSON.parse(JSON.stringify(this.value))
        cleanGenericRequestFromJson(this.refState)
    }

    get uploadTimeStr() : string {
        return standardFormatTime(this.value.UploadTime)
    }
}

</script>
