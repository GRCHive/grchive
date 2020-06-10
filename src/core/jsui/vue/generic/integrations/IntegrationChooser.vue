<template>
    <div>
        <v-list-item>
            <v-list-item-action id="searchContainer">
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>
        </v-list-item>

        <div class="px-4">
            <v-item-group :value="value">
                <v-row justify="center">
                    <v-col
                        cols="2"
                        v-for="(item, index) in filteredItems"
                        :key="`${index}-${item.text}`"
                    >
                        <v-item v-slot:default="{active, toggle}" :value="item.value">
                            <v-card
                                @click="selectIntegration(item.value, toggle, active)"
                                :class="active ? `integrationItem activeItem` : `integrationItem`"
                            >
                                <div class="py-2">
                                    <v-img
                                        :src="item.logo"
                                        class="mx-2"
                                    >
                                    </v-img>

                                    <div
                                        style="display: flex; justify-content: center;"
                                    >
                                        <span class="subtitle-2">{{ item.text }}</span>
                                    </div>
                                </div>
                            </v-card>
                        </v-item>
                    </v-col>
                </v-row>
            </v-item-group>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Prop } from 'vue-property-decorator'
import { IntegrationType } from '../../../ts/integrations/integration'

@Component
export default class IntegrationChooser extends Vue {
    @Prop()
    readonly value! : IntegrationType | null

    filterText : string = ""

    selectIntegration(val : IntegrationType, tg : any, active: boolean) {
        if (!active) {
            this.$emit('input', val)
        } else {
            this.$emit('input', null)
        }
        tg()
    }

    get integrationItems() : any[] {
        return [
            {
                value: IntegrationType.SapErp,
                text: 'SAP ERP',
                logo: '/static/assets/logos/saperp.png',
            }
        ]
    }

    get filteredItems() : any[] {
        let filter = this.filterText.trim()
        if (filter == "") {
            return this.integrationItems
        }

        return this.integrationItems.filter((ele : any) =>
            ele.text.toLowerCase().includes(filter.toLowerCase())
        )
    }
}

</script>

<style scoped>

#searchContainer {
    flex-grow: 1;
    margin-right: 0px;
}

.integrationItem {
    border-radius: 8px;
    border: 2px solid transparent;
}

.activeItem {
    border: 2px solid #1976d2;
}

</style>
