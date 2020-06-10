<template>
    <v-form :value="internalValid" @input="onChangeFormValid">
        <v-text-field
            label="Client"
            filled
            :rules="[rules.required]"
            :value="internalValue.Client"
            @input="onChangeValue('Client', arguments[0])"
            :readonly="readonly"
        >
        </v-text-field>

        <v-text-field
            label="Instance Number"
            filled
            :rules="[rules.required]"
            :value="internalValue.SysNr"
            @input="onChangeValue('SysNr', arguments[0])"
            :readonly="readonly"
        >
        </v-text-field>

        <v-text-field
            label="Host"
            filled
            :rules="[rules.required]"
            :value="internalValue.Host"
            @input="onChangeValue('Host', arguments[0])"
            :readonly="readonly"
        >
        </v-text-field>

        <v-text-field
            label="Real Hostname"
            filled
            :rules="[rules.required]"
            :value="internalValue.RealHostname"
            @input="onChangeValue('RealHostname', arguments[0])"
            v-if="isIpAddress"
            :readonly="readonly"
        >
        </v-text-field>

        <v-text-field
            label="Username"
            filled
            :rules="[rules.required]"
            :value="internalValue.Username"
            @input="onChangeValue('Username', arguments[0])"
            :readonly="readonly"
        >
        </v-text-field>

        <v-text-field
            label="Password"
            filled
            :rules="[rules.required]"
            :value="internalValue.Password"
            @input="onChangeValue('Password', arguments[0])"
            :type="displayPassword ? '' : 'password'"
            :readonly="readonly"
        >
            <template v-slot:append>
                <v-btn icon @click="displayPassword = !displayPassword">
                    <v-icon v-if="displayPassword">
                        mdi-eye-off
                    </v-icon>

                    <v-icon v-else>
                        mdi-eye
                    </v-icon>
                </v-btn>
            </template>
        </v-text-field>
    </v-form>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop, Watch } from 'vue-property-decorator'
import { SapErpIntegrationSetup, emptySapErpIntegrationSetup } from '../../../../../ts/integrations/sap'
import * as rules from '../../../../../ts/formRules'
import { Address4, Address6} from 'ip-address/ip-address.js'

@Component
export default class SapErpSetup extends Vue {
    @Prop()
    readonly value! : SapErpIntegrationSetup | null

    @Prop({ default: false })
    readonly valid! : boolean

    @Prop({ default: false })
    readonly readonly! : boolean

    rules: any = rules
    internalValid : boolean = false
    internalValue: SapErpIntegrationSetup = emptySapErpIntegrationSetup()
    displayPassword: boolean = false

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
            this.internalValue = emptySapErpIntegrationSetup()
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

    get isIpAddress() : boolean {
        let a4 = new Address4(this.internalValue.Host)
        if (a4.isValid()) {
            return true
        }

        let a6 = new Address6(this.internalValue.Host)
        if (a6.isValid()) {
            return true
        }

        return false
    }

    @Watch('isIpAddress')
    resetRealHostname() {
        if (this.isIpAddress) {
            if (!!this.value!.RealHostname) {
                this.onChangeValue('RealHostname', this.value!.RealHostname)
            } else {
                this.onChangeValue('RealHostname', '')
            }
        } else {
            this.onChangeValue('RealHostname', null)
        }
    }

    mounted() {
        this.syncValid()
        this.syncValue()
    }
}

</script>
