<template>
    <section class="ma-4">
        <v-dialog v-model="showHideDeleteControl" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="controls"
                :items-to-delete="currentControlsToDelete"
                v-on:do-cancel="showHideDeleteControl = false"
                v-on:do-delete="deleteSelectedControls"
                :use-global-deletion="true"
                :force-global-deletion="true">
            </generic-delete-confirmation-form>
        </v-dialog>

        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Controls
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
            <v-list-item-action>
                <v-dialog v-model="showHideCreateNewControl" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            Create New
                        </v-btn>
                    </template>
                    <create-new-control-form
                        :node-id="-1"
                        :risk-id="-1"
                        @do-save="saveNewControl"
                        @do-cancel="cancelNewControl">
                    </create-new-control-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
        <v-list-item class="headerItem">
            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Control
                </v-list-item-title>
            </v-list-item-content>
            <v-spacer></v-spacer>

            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Type
                </v-list-item-title>
            </v-list-item-content>
            <v-spacer></v-spacer>

            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Owner
                </v-list-item-title>
            </v-list-item-content>
            <v-spacer></v-spacer>

            <v-list-item-content class="font-weight-bold pa-0">
                <v-list-item-title>
                    Frequency
                </v-list-item-title>
            </v-list-item-content>
            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn icon disabled></v-btn>
            </v-list-item-action>
        </v-list-item>
        <v-card
            v-for="(item, index) in filteredRisks"
            :key="index"
            class="my-2"
        >
            <v-list-item two-line @click="goToControl(item.Id)">
                <v-list-item-content>
                    <v-list-item-title v-html="highlightText(item.Name)">
                    </v-list-item-title>
                    <v-list-item-subtitle v-html="highlightText(item.Description)">
                    </v-list-item-subtitle>
                </v-list-item-content>
                <v-spacer></v-spacer>

                <v-list-item-content>
                    {{ getTypeName(item.ControlTypeId) }}
                </v-list-item-content>
                <v-spacer></v-spacer>

                <v-list-item-content>
                    {{ getUserName(item.OwnerId) }}
                </v-list-item-content>
                <v-spacer></v-spacer>

                <v-list-item-content>
                    {{ createFrequencyDisplayString(item.FrequencyType, item.FrequencyInterval) }}
                </v-list-item-content>
                <v-spacer></v-spacer>

                <v-list-item-action>
                    <v-btn icon @click.stop="doDeleteControl(item)" @mousedown.stop @mouseup.stop>
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </v-card>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import { getAllControls, TAllControlInput, TAllControlOutput} from '../../../ts/api/apiControls'
import { deleteControls, TDeleteControlInput, TDeleteControlOutput} from '../../../ts/api/apiControls'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'
import { contactUsUrl, createControlUrl } from '../../../ts/url'
import { createFrequencyDisplayString } from '../../../ts/frequency'
import CreateNewControlForm from './CreateNewControlForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import MetadataStore from '../../../ts/metadata'

export default Vue.extend({
    data : () => ({
        allControls: [] as ProcessFlowControl[],
        filterText : "",
        showHideCreateNewControl : false,
        showHideDeleteControl : false,
        currentDeleteControl: Object() as ProcessFlowControl
    }),
    components: {
        CreateNewControlForm,
        GenericDeleteConfirmationForm
    },
    computed: {
        filter() : (a : ProcessFlowControl) => boolean {
            const filterText = this.filterText.trim()
            return (ele : ProcessFlowControl) : boolean => {
                return ele.Name.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) ||
                    ele.Description.toLocaleLowerCase().includes(filterText.toLocaleLowerCase())
            }
        },
        filteredRisks() : ProcessFlowRisk[] {
            return this.allControls.filter(this.filter)
        },
        currentControlsToDelete() : string[] {
            if (!this.showHideDeleteControl) {
                return []
            }
            return [this.currentDeleteControl.Name]
        }
    },
    methods: {
        highlightText(input : string) : string {
            const safeInput = sanitizeTextForHTML(input)
            const useFilter = this.filterText.trim()
            if (useFilter.length == 0) {
                return safeInput
            }
            return replaceWithMark(
                safeInput,
                sanitizeTextForHTML(useFilter))
        },
        refreshControls() {
            getAllControls(<TAllControlInput>{
                //@ts-ignore
                orgName: this.$root.orgGroupId
            }).then((resp : TAllControlOutput) => {
                this.allControls = resp.data
            }).catch((err : any) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        goToControl(controlId : number) {
            window.location.assign(createControlUrl(
                //@ts-ignore
                this.$root.orgGroupId,
                controlId))
        },
        deleteSelectedControls() {
            deleteControls(<TDeleteControlInput>{
                nodeId: -1,
                riskIds: [-1],
                controlIds: [this.currentDeleteControl.Id],
                global: true
            }).then(() => {
                this.allControls.splice(
                    this.allControls.findIndex((ele : ProcessFlowControl) =>
                        ele.Id == this.currentDeleteControl.Id),
                    1)
                this.showHideDeleteControl = false
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        doDeleteControl(control : ProcessFlowControl) {
            this.showHideDeleteControl = true
            this.currentDeleteControl = control
        },
        saveNewControl(control : ProcessFlowControl) {
            this.allControls.push(control)
            this.showHideCreateNewControl = false
        },
        cancelNewControl() {
            this.showHideCreateNewControl = false
        },
        getTypeName(typeId : number) : string {
            if (!(typeId in MetadataStore.state.idToControlTypes)) {
                return ""
            }
            return MetadataStore.state.idToControlTypes[typeId].Name
        },
        getUserName(userId : number | null) : string {
            if (userId == null) {
                return "No Owner"
            }

            if (!(userId in MetadataStore.state.idToUsers)) {
                return ""
            }

            return `${MetadataStore.state.idToUsers[userId].FirstName} ${MetadataStore.state.idToUsers[userId].LastName} [${MetadataStore.state.idToUsers[userId].Email}]`
        },
        createFrequencyDisplayString: createFrequencyDisplayString
    },
    mounted() {
        this.refreshControls()
    }
})
</script>

<style scoped>

.headerItem {
    max-height: 30px !important;
    min-height: 30px !important;
}

</style>
