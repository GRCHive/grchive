<template>
    <v-list-item
        class="pa-0"
    >
        <v-list-item-content class="mr-1">
            <v-text-field
                :value="paramName"
                @input="onChangeName"
                filled
                label="Name"
                dense
                :rules="rules"
            >
            </v-text-field>
        </v-list-item-content>

        <v-list-item-content class="ml-1">
            <supported-param-type-select-component
                :value="selectedType"
                @input="onChangeType"
                :rules="rules"
            >
            </supported-param-type-select-component>
        </v-list-item-content>
    </v-list-item>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    CodeParamType,
    SupportedParamType,
} from '../../../../ts/code'
import SupportedParamTypeSelectComponent from './SupportedParamTypeSelectComponent.vue'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => null as CodeParamType | null
        },
        rules: {
            type: Array,
            default: () => []
        }
    }
})

@Component({
    components: {
        SupportedParamTypeSelectComponent,
    }
})
export default class ParamTypeComponent extends Props {
    selectedType: SupportedParamType | null = null

    get paramName() : string {
        if (!this.value) {
            return ""
        }

        return this.value.Name
    }

    get paramTypeId() : number {
        if (!this.value) {
            return -1
        }

        return this.value.ParamId
    }

    onChangeName(nm : string) {
        let newVal = <CodeParamType>{
            Name: nm,
            ParamId: this.paramTypeId,
        }

        this.$emit('input', newVal)
    }

    onChangeType(typ : SupportedParamType) {
        this.selectedType = typ
        let newVal = <CodeParamType>{
            Name: this.paramName,
            ParamId: typ.Id,
        }
        this.$emit('input', newVal)
    }
}

</script>
