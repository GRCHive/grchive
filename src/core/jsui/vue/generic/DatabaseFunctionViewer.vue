<template>
    <div>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-list-item class="pa-0">
                <v-list-item-content id="querySelector">
                    <v-select
                        label="Functions"
                        filled
                        hide-details
                        dense
                        :items="fnItems"
                        :value="currentFunction"
                        @input="selectFunction"
                    >
                    </v-select>
                </v-list-item-content>
            </v-list-item>

            <sql-text-area
                :value="currentFunction.Src"
                readonly
                :key="fnKey"
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

const Props = Vue.extend({
    props: {
        schema: {
            type: Object,
            default: Object() as () => DbSchema | null
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

    get isLoading() : boolean {
        return this.allFunctions == null
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

        if (!this.schema) {
            return
        }

        getSqlSchema({
            schemaId: this.schema!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            fnMode: true
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
