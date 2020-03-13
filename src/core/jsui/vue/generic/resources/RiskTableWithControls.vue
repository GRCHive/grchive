<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Risks
                </v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>
            <v-spacer></v-spacer>
            <v-list-item-action v-if="!disableNew">
                <v-dialog v-model="showHideCreateNewRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            New
                        </v-btn>
                    </template>
                    <create-new-risk-form
                        :node-id="-1"
                        @do-save="saveNewRisk"
                        @do-cancel="cancelNewRisk">
                    </create-new-risk-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

        <advanced-risk-filters
            v-model="riskFilter"
        >
        </advanced-risk-filters>

        <v-divider></v-divider>

        <risk-table
            :value="value"
            :resources="filteredRisks"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteSelectedRisk"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        >
        </risk-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { getAllRisks, TAllRiskOutput } from '../../../ts/api/apiRisks'
import { deleteRisk, TDeleteRiskOutput } from '../../../ts/api/apiRisks'
import { contactUsUrl } from '../../../ts/url'
import CreateNewRiskForm from '../../components/dashboard/CreateNewRiskForm.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import RiskTable from '../RiskTable.vue'
import AdvancedRiskFilters from '../filters/AdvancedRiskFilters.vue'
import { RiskFilterData, NullRiskFilterData } from '../../../ts/risks'

const Props = Vue.extend({
    props: {
        value: {
            type: Array,
            default: () => [],
        },
        exclude: {
            type: Array,
            default: () => [],
        },
        disableNew: {
            type: Boolean,
            default: false,
        },
        disableDelete: {
            type: Boolean,
            default: false,
        },
        enableSelect: {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        CreateNewRiskForm,
        RiskTable,
        AdvancedRiskFilters
    }
})
export default class RiskTableWithControls extends Props {
    allRisks: ProcessFlowRisk[] = []
    filterText : string = ""
    showHideCreateNewRisk: boolean = false
    riskFilter: RiskFilterData = JSON.parse(JSON.stringify(NullRiskFilterData))

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredRisks() : ProcessFlowRisk[] {
        console.log(this.allRisks)
        let test : ProcessFlowRisk[] = this.allRisks.filter((ele : ProcessFlowRisk) => !this.excludeSet.has(ele.Id))
        console.log(test)
        return test
    }

    @Watch('riskFilter', {deep:true})
    refreshRisks() {
        getAllRisks({
            orgName: PageParamsStore.state.organization!.OktaGroupName,
            filter: this.riskFilter,
        }).then((resp : TAllRiskOutput) => {
            this.allRisks = resp.data
        }).catch((err) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    saveNewRisk(risk : ProcessFlowRisk) {
        this.allRisks.unshift(risk)
        this.showHideCreateNewRisk = false
    }

    cancelNewRisk() {
        this.showHideCreateNewRisk = false
    }

    deleteSelectedRisk(risk : ProcessFlowRisk, global : boolean) {
        deleteRisk({
            nodeId: -1,
            riskIds: [risk.Id],
            global: true,
        }).then((resp : TDeleteRiskOutput) => {
            this.allRisks.splice(
                this.allRisks.findIndex((ele : ProcessFlowRisk) =>
                    ele.Id == risk.Id),
                1)
        }).catch((err) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    modifySelected(vals : ProcessFlowRisk[]) {
        this.$emit('input', vals)
    }

    mounted() {
        this.refreshRisks()
    }
}

</script>
