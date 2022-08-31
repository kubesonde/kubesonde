import {Column, useTable, useExpanded, Renderer, CellProps,} from 'react-table'
import React, {useEffect, useMemo, useState} from "react";
import {computeMetrics} from "../../../utils/graph";
import "../graphTable/table.css"
import {Graph} from "../../../entities/graph";
interface StatsTableItem {
    stat: string
    value: string[]
}


const ConnectedComponentCell: Renderer<CellProps<string[]>> = (row: CellProps<string[]>) => {
    return(
        <div style={{textAlign:'left'}} >
            {
                (row.value as string[])
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



    const columns: Column<StatsTableItem>[] = useMemo(() => ([
        {
            Header: 'Statistics',
            accessor: 'stat',
        },
        {
            id: 'checkbox-table-column',
            accessor: 'value',
            Cell: ConnectedComponentCell
        },
    ]), [])


    const {
        getTableProps,
        getTableBodyProps,
        headerGroups,
        rows,
        prepareRow,
    } = useTable({
        columns,
        data,
    }, useExpanded)


    // Render Table UI
    return (
        <table {...getTableProps()}>
            <thead>
            {headerGroups.map(headerGroup => (
                <tr {...headerGroup.getHeaderGroupProps()}>
                    {headerGroup.headers.map(column => (
                        <th {...column.getHeaderProps()}>{column.render('Header')}</th>
                    ))}
                </tr>
            ))}
            </thead>
            <tbody {...getTableBodyProps()}>
            {rows.map((row, i) => {
                prepareRow(row)
                return (
                    <tr {...row.getRowProps()}>
                        {row.cells.map(cell => {
                            return <td {...cell.getCellProps()} >{cell.render('Cell')}</td>
                        })}
                    </tr>
                )
            })}
            </tbody>
        </table>
    )
}
