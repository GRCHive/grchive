<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Systems
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
                <v-dialog v-model="showHideNew"
                          persistent
                          max-width="40%"
                >
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>

                    <create-new-system-form
                        @do-cancel="showHideNew = false"
                        @do-save="saveSystem">
                    </create-new-system-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
        <systems-table
            :value="value"
            :resources="filteredSystems"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteSystem"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        ></systems-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import CreateNewSystemForm from '../../components/dashboard/CreateNewSystemForm.vue'
import { System } from '../../../ts/systems'
import { TAllSystemsOutputs, getAllSystems } from '../../../ts/api/apiSystems'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { deleteSystem, TDeleteSystemOutputs } from '../../../ts/api/apiSystems'
import SystemsTable from '../SystemsTable.vue'

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
        CreateNewSystemForm,
        SystemsTable
    }
})
export default class SystemTableWithControls extends Props {
    showHideNew: boolean = false
    systems : System[] = []
    filterText : string = ""

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredSystems() : System[] {
        return this.systems.filter((ele : System) => !this.excludeSet.has(ele.Id))
    }

    saveSystem(newSys : System) {
        this.showHideNew = false
        this.systems.unshift(newSys)
    }

    refreshSystems() {
        getAllSystems({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TAllSystemsOutputs) => {
            this.systems = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    deleteSystem(sys : System) {
        deleteSystem({
            sysId: sys.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TDeleteSystemOutputs) => {
            this.systems.splice(
                this.systems.findIndex((ele : System) =>
                    ele.Id == sys.Id),
                1)
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

    modifySelected(vals : System[]) {
        this.$emit('input', vals)
    }

    mounted() {
        this.refreshSystems()
    }
}

</script>
