<template>
    <div>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-list-item class="pa-0">
                <v-list-item-content id="querySelector">
                    <v-autocomplete
                        label="Functions"
                        filled
                        hide-details
                        dense
                        :items="fnItems"
                        :value="currentFunction"
                        @input="selectFunction"
                    >
                    </v-autocomplete>
                </v-list-item-content>
            </v-list-item>

            <sql-text-area
                v-if="!!currentFunction"
                :value="prettySrc"
                readonly
                :key="fnKey"
                full-height
            >
            </sql-text-area>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { getSqlSchema, TGetSqlSchemaOutput } from '../../ts/api/apiSqlSchemas'
import {
    DbSchema,
    DbFunction,
} from '../../ts/sql'
import { Watch } from 'vue-property-decorator'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import SqlTextArea from './SqlTextArea.vue'
import sqlFormatter from 'sql-formatter'

const Props = Vue.extend({
    props: {
        schema: {
            type: Object,
            default: () => null as DbSchema | null
        }
    }
})

@Component({
    components: {
        SqlTextArea,
    }
})
export default class DatabaseFunctionViewer extends Props {
    allFunctions : DbFunction[] | null = null
    currentFunction : DbFunction | null = null

    fnKey : number = 0

    get prettySrc() : string {
        if (!this.currentFunction) {
            return ""
        }

        return sqlFormatter.format(this.currentFunction.Src)
    }

    get isLoading() : boolean {
        return (!this.allFunctions)
    }

    get fnItems() : any[] {
        if (!this.allFunctions) {
            return []
        }
        return this.allFunctions.map((ele : DbFunction) => {
            let fnOrProcedure : string = ""
            if (!!ele.RetType) {
                fnOrProcedure = `FUNCTION RETURNS ${ele.RetType}`
            } else {
                fnOrProcedure = "PROCEDURE"
            }

            return {
                text: `${ele.Name} (${fnOrProcedure})`,
                value: ele,
            }
        })
    }

    selectFunction(fn : DbFunction) {
        this.currentFunction = fn
        this.fnKey += 1
    }

    @Watch('schema')
    refreshData() {
        this.allFunctions = null
        this.currentFunction = null

        if (!this.schema) {
            return
        }

        getSqlSchema({
            schemaId: this.schema!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            fnMode: true,
            start: -1,
            limit: -1,
        }).then((resp : TGetSqlSchemaOutput) => {
            this.allFunctions = resp.data.Functions!
            if (this.allFunctions!.length > 0) {
                this.selectFunction(this.allFunctions![0])
            }
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
        this.refreshData()
    }
}

</script>
