<template>
    <div>
        <p class="title">Data Source</p>
        <v-select
            filled
            v-for="d in depthToShow"
            :key="d - 1"
            :items="levelTags(d - 1)"
            :rules="rules"
            @input="changeSelectedTag(d - 1, arguments[0])"
        >
        </v-select>

        <div v-if="!!currentNode && currentNode.isLeaf">
            <div v-if="currentNode._val.Name == 'Root.Database.PostgreSQL'">
                <database-search-form-component
                    label="Database"
                    :value="db"
                    :rules="rules"
                    :type-id="forceDbTypeId"
                    @input="onSelectDb"
                >
                </database-search-form-component>
            </div>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    DataSourceLink,
    DataSourceOption,
    DataSourceOptionNode,
    DataSourceOptionTree
} from '../../ts/clientData'
import { allSupportedDataSources, TAllDataSourceOutput } from '../../ts/api/apiClientData'
import { contactUsUrl } from '../../ts/url'
import { Database } from '../../ts/databases'
import { PageParamsStore } from '../../ts/pageParams'
import DatabaseSearchFormComponent from './DatabaseSearchFormComponent.vue'

const Props = Vue.extend({
    props: {
        value: {
            type: Object,
            default: () => Object() as DataSourceLink | null
        },
        rules: {
            type: Array,
            default: () => []
        },
        dataId: {
            type: Number,
            default: -1,
        }
    }
})

@Component({
    components: {
        DatabaseSearchFormComponent
    }
})
export default class DataSourceFormComponent extends Props {
    options : DataSourceOption[] = []
    selectedTags : string[] = []

    db : Database | null = null

    get tree() : DataSourceOptionTree {
        return new DataSourceOptionTree(this.options)
    }

    get currentNode() : DataSourceOptionNode | null {
        return this.tree.traverse(this.selectedTags)
    }

    get depthToShow() : number {
        if (!this.currentNode) {
            return 0
        }

        if (this.currentNode.isLeaf) {
            return this.selectedTags.length
        } else {
            return this.selectedTags.length + 1
        }
    }

    get levelTags() : (i : number) => string[] {
        return (i : number) : string[] => {
            let node = this.tree.traverse(this.selectedTags.slice(0, i))
            if (!node) {
                return []
            }
            return node.getChildTags()
        }
    }

    changeSelectedTag(d : number, tag : string) {
        if (d > this.selectedTags.length) {
            // This is an error?
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        } else if (d == this.selectedTags.length) {
            this.selectedTags.push(tag)
        } else {
            this.selectedTags[d] = tag
            this.selectedTags = this.selectedTags.slice(0, d + 1)
            this.clearSource()
        }

        this.sync()
    }

    refreshOptions() {
        allSupportedDataSources().then((resp : TAllDataSourceOutput) => {
            this.options = resp.data
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
        this.refreshOptions()
    }

    onSelectDb(db : Database) {
        this.db = db
        this.sync()
    }

    clearSource() {
        this.db = null
    }

    get sourceTarget() : Record<string, any> {
        if (!this.currentNode || !this.currentNode.isLeaf) {
            return {}
        }

        if (this.currentNode._val!.Name == 'Root.Database.PostgreSQL' && !!this.db) {
            return {
                'id': this.db.Id,
            }
        }
        return {}
    }

    get forceDbTypeId() : number {
        if (!this.currentNode || !this.currentNode.isLeaf) {
            return -1
        }

        if (this.currentNode._val!.Name == 'Root.Database.PostgreSQL') {
            return 1
        }
        return -1
    }

    sync() {
        if (!!this.currentNode && this.currentNode.isLeaf) {
            let obj : DataSourceLink = {
                OrgId: PageParamsStore.state.organization!.Id,
                DataId: this.dataId,
                SourceId: this.currentNode._val!.Id,
                SourceTarget: this.sourceTarget,
            }
            this.$emit('input', obj)
        } else {
            this.$emit('input', null)
        }
    }
}

</script>
