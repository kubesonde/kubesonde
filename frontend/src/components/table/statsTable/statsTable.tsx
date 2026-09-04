import {
    ColumnDef,
    useReactTable,
    getCoreRowModel,
    getExpandedRowModel,
    flexRender,
    CellContext,
} from '@tanstack/react-table'
import React, {useEffect, useMemo, useState} from "react";
import {computeMetrics} from "../../../utils/graph";
import "../graphTable/table.css"
import {Graph} from "../../../entities/graph";
interface StatsTableItem {
    stat: string
    value: string[]
}


const ConnectedComponentCell = (info: CellContext<StatsTableItem, string[]>) => {
    return(
        <div style={{textAlign:'left'}} >
            {
                info.getValue()
                    .map((item,index) => (
                        <div
                            key={index}>
                            {item}
                        </div>) )
            }
        </div>)
}

export const StatsTable: React.FC<Graph> = (props: Graph) => {
    const [data,setData] = useState([
        {
            stat: 'Strongly connected components',
            value: [""]
        },
        {
            stat: 'Average out degree',
            value: [""]
        },
        {
            stat: 'Clustering',
            value: [""]
        }])
    useEffect(() => {
        const metrics = computeMetrics(props)


        setData([
            {
                stat: 'Strongly connected components',
                value: metrics.ssc
            },
            {
                stat: 'Average out degree',
                value: [metrics.avgOutDegree.toPrecision(2)]
            },
            {
                stat: 'Clustering',
                value: ["TBA"]
            }])
    },[props])



    const columns: ColumnDef<StatsTableItem, any>[] = useMemo(() => ([
        {
            id: 'stat',
            header: 'Statistics',
            accessorKey: 'stat',
        },
        {
            id: 'checkbox-table-column',
            accessorKey: 'value',
            cell: ConnectedComponentCell
        },
    ]), [])


    const table = useReactTable({
        data,
        columns,
        getCoreRowModel: getCoreRowModel(),
        getExpandedRowModel: getExpandedRowModel(),
    })


    // Render Table UI
    return (
        <table>
            <thead>
            {table.getHeaderGroups().map(headerGroup => (
                <tr key={headerGroup.id}>
                    {headerGroup.headers.map(header => (
                        <th key={header.id}>{header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}</th>
                    ))}
                </tr>
            ))}
            </thead>
            <tbody>
            {table.getRowModel().rows.map((row) => {
                return (
                    <tr key={row.id}>
                        {row.getVisibleCells().map(cell => {
                            return <td key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</td>
                        })}
                    </tr>
                )
            })}
            </tbody>
        </table>
    )
}
