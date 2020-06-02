<script lang="ts">

import Vue, {VNode} from 'vue'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import MetadataStore from '../../ts/metadata'
import ShellRunStatus from './shell/ShellRunStatus.vue'
import { createUserString } from '../../ts/users'
import { VChip, VProgressCircular } from 'vuetify/lib'
import { PageParamsStore } from '../../ts/pageParams'
import { standardFormatTime } from '../../ts/time'
import {
    contactUsUrl,
    createSingleShellUrl,
} from '../../ts/url'
import { sortDate } from '../../ts/time'
import { ShellScriptRunPerServer } from '../../ts/shell'
import {
    getShellRunInformation, TGetShellRunOutput,
} from '../../ts/api/apiShellRun'

@Component({
    components: {
        BaseResourceTable,
        ShellRunStatus,
    }
})
export default class ShellRunTable extends ResourceTableProps {
    idProcessed: Set<number> = new Set<number>()

    get tableHeaders() : any[] {
        return [
            {
                text: 'Shell',
                value: 'shell',
            },
            {
                text: 'Version',
                value: 'version',
            },
            {
                text: 'User',
                value: 'user',
            },
            {
                text: 'Create Time',
                value: 'createTime',
                sort: sortDate,
            },
            {
                text: 'Start Time',
                value: 'startTime',
                sort: sortDate,
            },
            {
                text: 'End Time',
                value: 'endTime',
                sort: sortDate,
            },
            {
                text: 'Progress',
                value: 'progress',
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    goToShellScriptRun(item : any) {
        //window.location.assign(createSingleShellUrl(
        //    PageParamsStore.state.organization!.OktaGroupName,
        //    item.value.Id) + `?version=${item.currentVersionNumber}`)
    }

    transformInputResourceToTableItem(inp : any) : any {
        let obj = {
            id: inp.Id,
            shell: null,
            version: null,
            user: createUserString(MetadataStore.getters.getUser(inp.RunUserId)),
            createTime: inp.CreateTime,
            startTime: inp.RunTime,
            endTime: inp.EndTime,
            progress: null,
            value: inp,
        }

        return obj
    }

    loadAuxData(runId : number, obj : any) {
        if (this.idProcessed.has(runId)) {
            return
        }

        this.idProcessed.add(runId)
        getShellRunInformation({
            orgId: PageParamsStore.state.organization!.Id,
            runId: runId,
            includeLogs: false,
        }).then((resp : TGetShellRunOutput) => {
            obj.shell = resp.data.Script
            obj.version = resp.data.VersionNum
            obj.progress = resp.data.ServerRuns
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

    renderShell(props : any) : VNode {
        this.loadAuxData(
            props.item.value.Id,
            props.item
        )

        if (!props.item.shell) {
            return this.$createElement('span', 'Loading...')
        }

        return this.$createElement(
            'a',
            {
                attrs: {
                    href: createSingleShellUrl(PageParamsStore.state.organization!.OktaGroupName, props.item.shell.Id) + `?version=${props.item.version}`
                }
            },
            props.item.shell.Name
        )
    }

    renderVersion(props : any) : VNode {
        let chip = this.$createElement(
            VChip,
            {
                props: {
                    pill: true
                },
            },
            [
                !!props.item.version ? props.item.version : "Loading...",
            ]
        )
        return chip
    }

    renderCreateTime(props : any) : VNode {
        return this.$createElement(
            'span',
            standardFormatTime(props.item.createTime),
        )
    }

    renderStartTime(props : any) : VNode {
        return !!props.item.startTime ?
            this.$createElement('span', standardFormatTime(props.item.startTime)) :
            this.$createElement(VProgressCircular, { props: {size: 16} })
    }

    renderEndTime(props : any) : VNode {
        return !!props.item.endTime ?
            this.$createElement('span', standardFormatTime(props.item.endTime)) :
            this.$createElement(VProgressCircular, { props: {size: 16} })
    }

    renderProgress(props: any) : VNode {
        if (!props.item.progress) {
            return this.$createElement('span', 'Loading...')
        }

        return this.$createElement(
            ShellRunStatus,
            {
                props: {
                    serverRuns: props.item.progress,
                }
            },
        )
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    resourceName: "shell script runs",
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToShellScriptRun,
                },
                scopedSlots: {
                    'item.shell': this.renderShell,
                    'item.version': this.renderVersion,
                    'item.createTime': this.renderCreateTime,
                    'item.startTime': this.renderStartTime,
                    'item.endTime': this.renderEndTime,
                    'item.progress': this.renderProgress,
                }
            }
        )
    }

}

</script>

