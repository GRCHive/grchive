<template>
    <v-form :value="internalValid" @input="onChangeFormValid">
        <v-text-field
            label="Name"
            filled
            :rules="[rules.required, rules.createMaxLength(255)]"
            :value="internalValue.Name"
            @input="onChangeValue('Name', arguments[0])"
            :readonly="readonly"
        >
        </v-text-field>

        <v-textarea
            label="Description"
            filled
            :value="internalValue.Description"
            @input="onChangeValue('Description', arguments[0])"
            :readonly="readonly"
        >
        </v-textarea>
    </v-form>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import {
    GenericIntegration, emptyGenericIntegration,
    IntegrationType
} from '../../../ts/integrations/integration'
import * as rules from '../../../ts/formRules'

@Component
export default class GenericIntegrationForm extends Vue {
    @Prop()
    readonly value! : GenericIntegration | null

    @Prop()
    readonly valid! : boolean

    @Prop()
    readonly type! : IntegrationType | null

    @Prop({default: false})
    readonly readonly!: boolean 

    rules: any = rules
    internalValid : boolean = false
    internalValue: GenericIntegration = emptyGenericIntegration()

    @Watch('type')
    syncType() {
        this.onChangeValue('Type', this.type)
    }

    @Watch('valid')
    syncValid() {
        if (this.valid != this.internalValid) {
            this.internalValid = this.valid
        }
    }

    @Watch('value')
    syncValue() {
        if (!!this.value) {
            this.internalValue = JSON.parse(JSON.stringify(this.value))
        } else {
            this.internalValue = emptyGenericIntegration()
        }

        if (!!this.type) {
            this.internalValue.Type = this.type
        }
    }

    onChangeFormValid(e : boolean) {
        this.internalValid = e
        this.$emit('update:valid', e)
    }

    onChangeValue(field : string, val : any) {
        Vue.set(this.internalValue, field, val)
        this.$emit('input', this.internalValue)
    }

    mounted() {
        this.syncValue()
    }
}

</script>
