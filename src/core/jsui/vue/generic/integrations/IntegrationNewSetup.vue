<template>
    <div>
        <v-stepper
            v-model="step"
            class="integrationStepper"
        >
            <v-stepper-header class="integrationStepperHeader mx-8">
                <v-stepper-step
                    :complete="step > 1"
                    :step="1"
                >
                    Select Integration
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step
                    :complete="step > 2"
                    :step="2"
                >
                    GRCHive Properties
                </v-stepper-step>

                <v-divider></v-divider>

                <v-stepper-step
                    :complete="step > 3"
                    :step="3"
                >
                    Setup Integration
                </v-stepper-step>
            </v-stepper-header>
            <v-divider></v-divider>

            <v-stepper-items>
                <v-stepper-content
                    :step="1"
                    class="integrationContent"
                >
                    <integration-chooser
                        v-model="chosenIntegration"
                    >
                    </integration-chooser>
                </v-stepper-content>

                <v-stepper-content
                    :step="2"
                    class="integrationContent"
                >
                    <generic-integration-form
                        class="mx-4"
                        v-model="integration"
                        :valid.sync="integrationIsValid"
                        :type="chosenIntegration"
                    >
                    </generic-integration-form>
                </v-stepper-content>

                <v-stepper-content
                    :step="3"
                    class="integrationContent"
                >
                    <v-row class="center title" justify="center">
                        <span v-if="isSapErp">SAP ERP</span>
                    </v-row>
                    <sap-erp-setup
                        class="mx-4"
                        v-if="isSapErp"
                        v-model="setupInfo"
                        :valid.sync="setupIsValid"
                    >
                    </sap-erp-setup>
                </v-stepper-content>
            </v-stepper-items>
        </v-stepper>

        <v-list-item>
            <v-list-item-action>
                <v-btn
                    color="error"
                    @click="cancel"
                >
                    Cancel
                </v-btn>
            </v-list-item-action>

            <v-spacer></v-spacer>

            <v-list-item-action v-if="step > 1">
                <v-btn
                    color="secondary"
                    @click="step -= 1"
                >
                    Back
                </v-btn>
            </v-list-item-action>

            <v-list-item-action v-if="step < 3">
                <v-btn
                    color="primary"
                    @click="step += 1"
                    :disabled="!canGoNext"
                >
                    Next
                </v-btn>
            </v-list-item-action>

            <v-list-item-action v-if="step == 3">
                <v-btn
                    color="success"
                    @click="save"
                    :disabled="!setupIsValid"
                    :loading="saving"
                >
                    Save
                </v-btn>
            </v-list-item-action>
        </v-list-item>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import IntegrationChooser from './IntegrationChooser.vue'
import GenericIntegrationForm from './GenericIntegrationForm.vue'
import { IntegrationType, GenericIntegration, UnionIntegrationSetup } from '../../../ts/integrations/integration'
import { SapErpIntegrationSetup } from '../../../ts/integrations/sap'

const SapErpSetup = () => import (/* webpackChunkName: "SapErpSetup" */ './sap/erp/SapErpSetup.vue')

@Component({
    components: {
        IntegrationChooser,
        SapErpSetup,
        GenericIntegrationForm,
    }
})
export default class IntegrationNewSetup extends Vue {
    step : number = 1
    chosenIntegration : IntegrationType | null = null
    integration : GenericIntegration | null = null
    integrationIsValid : boolean = false

    setupInfo : UnionIntegrationSetup = null
    setupIsValid: boolean = false
    saving: boolean = false

    get isSapErp() : boolean {
        return !!this.chosenIntegration && this.chosenIntegration == IntegrationType.SapErp
    }

    get canGoNext() : boolean {
        if (this.step == 1) {
            return !!this.chosenIntegration
        } else if (this.step == 2) {
            return this.integrationIsValid
        } else {
            return this.setupIsValid
        }
        return false
    }

    cleanup() {
        this.step = 1
        this.chosenIntegration = null
        this.setupInfo = null
        this.setupIsValid = false
        this.saving = false
    }

    clearSaving() {
        this.saving = false
    }

    cancel() {
        this.cleanup()
        this.$emit('cancel')
    }

    save() {
        this.saving = true
        this.$emit('save', this.integration!, this.setupInfo!)
    }
}

</script>

<style scoped>

.integrationStepper {
    box-shadow: none;
}

.integrationStepperHeader {
    box-shadow: none;
}

.integrationContent {
    padding: 4px 4px 4px 4px;
}

</style>
