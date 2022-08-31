import { GraphEdge, GraphNode, SimpleGraphEdge } from "../entities/graph";
import {
    getColoredData,
    getDeployments, hidePorts,
    mergeEdges, mergeEdgesSimple,
    showAllGraph,
    toClusteredEdge,
    toClusteredNode,
    toCyNode,
    uniqueNodesReducer
} from "./graph";
import { GraphEdgeBuilder, GraphNodeBuilder, SimpleGraphEdgeBuilder } from "../entities/graph.builder";


beforeEach(() => jest.resetAllMocks())
describe('toClusteredNode', () => {
    it('toClusteredNode with deployment', () => {
        const node: GraphNode = {
            //    type: "Pod",
            deployment: "deployment",
            shape: "ellipse",
            group: "group",
            id: "ID",
            label: "label",
            name: "name"
        }
        expect(toClusteredNode(node)).toEqual({
            deployment: "deployment",
            group: "group",
            id: "deployment",
            label: "deployment",
            name: "name",
            shape: "square"
        })
    })
    it('toClusteredNode with no deployment', () => {
        const node: GraphNode = {
            //    type: "Pod",
            deployment: undefined,
            shape: "ellipse",
            group: "group",
            id: "ID",
            label: "label",
            name: "name"
        }
        expect(toClusteredNode(node)).toEqual({
            deployment: undefined,
            group: "group",
            id: "ID",
            label: "label",
            name: "name",
            shape: "ellipse"
        })
    })
})

describe('toClusteredEdge', () => {
    it('toClusteredEdge with deployment', () => {
        const edge: GraphEdge = {
            from: "A", id: "myEdge", ports: [], to: "B",
            fromDeployment: "f-deployment",
            toDeployment: "t-deployment",
            label: "label",
            deniedConnection: false
        }
        expect(toClusteredEdge(edge)).toEqual({
            id: "myEdge",
            from: "f-deployment",
            to: "t-deployment",
            fromDeployment: "f-deployment",
            toDeployment: "t-deployment",
            label: "label",
            ports: [],
            deniedConnection: false
        })
    })
    it('toClusteredEdge with deployment', () => {
        const edge: GraphEdge = {
            from: "A", id: "myEdge", ports: [], to: "B",
            fromDeployment: undefined,
            toDeployment: "t-deployment",
            label: "label",
            deniedConnection: false
        }
        expect(toClusteredEdge(edge)).toEqual({
            id: "myEdge",
            from: "A",
            to: "t-deployment",
            fromDeployment: undefined,
            toDeployment: "t-deployment",
            label: "label",
            ports: [],
            deniedConnection: false
        })
    })
})

describe('uniqueNodesReducer', () => {
    it('should reduce nodes with same ID', function () {
        const nodes: GraphNode[] = [
            {
                id: "deployment1",
                name: "pod1",
                label: "label1",
                deployment: undefined,
                group: undefined
            },
            {
                id: "deployment1",
                name: "pod2",
                label: "label2",
                deployment: undefined,
                group: undefined
            }
        ]
        const result = nodes.reduce(uniqueNodesReducer, [])
        expect(result).toEqual([{
            id: "deployment1",
            name: "pod1",
            label: "label1",
            deployment: undefined,
            group: undefined,
            title: undefined

        }])


    });
    it('should reduce nodes when single title', function () {
        const nodes: GraphNode[] = [
            {
                id: "deployment1",
                name: "pod1",
                label: "label1",
                deployment: undefined,
                group: undefined
            },
            {
                id: "deployment1",
                name: "pod2",
                label: "label2",
                deployment: undefined,
                group: undefined,
                title: "mytitle"
            }
        ]
        const result = nodes.reduce(uniqueNodesReducer, [])
        expect(result).toEqual([{
            id: "deployment1",
            name: "pod1",
            label: "label1",
            deployment: undefined,
            group: undefined,
            title: "mytitle"

        }])


    });
    it('should reduce nodes when multiple', function () {
        const nodes: GraphNode[] = [
            {
                id: "deployment1",
                name: "pod1",
                label: "label1",
                deployment: undefined,
                group: undefined,
                title: "pod1"
            },
            {
                id: "deployment1",
                name: "pod2",
                label: "label2",
                deployment: undefined,
                group: undefined,
                title: "pod2"
            }
        ]
        const result = nodes.reduce(uniqueNodesReducer, [])
        expect(result).toEqual([{
            id: "deployment1",
            name: "pod1",
            label: "label1",
            deployment: undefined,
            group: undefined,
            title: "pod1\npod2"

        }])


    });

})

describe('mergeEdges', () => {
    it('merges edges that are not hidden', () => {

        const edges = [
            {
                ...GraphEdgeBuilder()
                    .id("1")
                    .label("80")
                    .ports(['80']).build()
            },
            {
                ...GraphEdgeBuilder()
                    .id("2")
                    .label("8080")
                    .ports(['8080']).build()
            },
        ]
        expect(mergeEdges(edges)).toEqual([
            {
                ...GraphEdgeBuilder()
                    .id("1")
                    .label("80, 8080")
                    .ports(["80", "8080"]).build()
            }
        ])

    })
    it('does not merge edges that are hidden', () => {

        const edges = [
            {
                ...GraphEdgeBuilder()
                    .hidden(true)
                    .id("1")
                    .label("80")
                    .ports(['80']).build()
            },
            {
                ...GraphEdgeBuilder()
                    .hidden(true)
                    .id("2")
                    .label("8080")
                    .ports(['8080']).build()
            },
        ]
        expect(mergeEdges(edges)).toEqual(edges)

    })
    it('should merge edges that are the same', () => {

        const edges = [
            GraphEdgeBuilder()
                .id("1")
                .label("80")
                .ports(['80']).build(),
            GraphEdgeBuilder()
                .id("2")
                .label("80")
                .ports(['80']).build(),
        ]
        expect(mergeEdges(edges)).toEqual([
            GraphEdgeBuilder()
                .id("1")
                .label("80")
                .ports(["80"]).build(),
        ])

    })

})

describe('getDeployments', () => {
    it('gets the deployments', () => {
        const nodes: GraphNode[] = [
            GraphNodeBuilder()
                .id("1").build(),
            GraphNodeBuilder()
                .id("2")
                .deployment("Deployment1").build(),
            GraphNodeBuilder()
                .id("3")
                .deployment("Deployment1").build(),
            GraphNodeBuilder()
                .id("4")
                .deployment("Deployment2").build()
        ]

        expect(getDeployments(nodes)).toEqual(["Deployment1", "Deployment2"])
    })
})

describe('showAllGraph', () => {
    it('should show all items', function () {
        const nodes: GraphNode[] = [
            GraphNodeBuilder().id("1").hidden(true).build(),
            GraphNodeBuilder().id("2").hidden(false).deployment("Deployment1").build()
        ]
        const edges = [
            GraphEdgeBuilder().id("1").hidden(true).label("80").ports(['80']).build(),
            GraphEdgeBuilder().id("2").hidden(false).label("8080").ports(['8080']).build()
        ]

        expect(showAllGraph({ nodes, edges }))
            .toEqual({
                nodes: nodes.map((node) => ({ ...node, hidden: false })),
                edges: edges.map((node) => ({ ...node, hidden: false })),
            })
    });
})

describe('getColoredData', () => {
    it('works when generic node is provided', () => {
        const data = {
            edges: [{
                from: "A", id: "myEdge", ports: [], to: "B",
                fromDeployment: "f-deployment",
                toDeployment: "t-deployment",
                label: "label",
                deniedConnection: false
            }],
            nodes: [{
                id: "pod1",
                name: "pod1",
                label: "label1",
                deployment: undefined,
                group: undefined
            }]
        }
        expect(getColoredData(data, { 'pod1': 'green' }))
            .toEqual({
                edges: [{
                    from: "A", id: "myEdge", ports: [], to: "B",
                    fromDeployment: "f-deployment",
                    toDeployment: "t-deployment",
                    label: "label",
                    deniedConnection: false
                }],
                nodes: [{
                    id: "pod1",
                    name: "pod1",
                    label: "label1",
                    color: "green",
                    deployment: undefined,
                    group: undefined
                }]
            })
    })

    it('works when node with deployment', () => {
        const data = {
            edges: [{
                from: "A", id: "myEdge", ports: [], to: "B",
                fromDeployment: "f-deployment",
                toDeployment: "t-deployment",
                label: "label",
                deniedConnection: false
            }],
            nodes: [{
                id: "pod1",
                name: "pod1",
                label: "label1",
                deployment: "dep1",
                group: undefined
            }]
        }
        expect(getColoredData(data, { 'dep1': 'green' }))
            .toEqual({
                edges: [{
                    from: "A", id: "myEdge", ports: [], to: "B",
                    fromDeployment: "f-deployment",
                    toDeployment: "t-deployment",
                    label: "label",
                    deniedConnection: false
                }],
                nodes: [{
                    id: "pod1",
                    name: "pod1",
                    label: "label1",
                    color: "green",
                    deployment: 'dep1',
                    group: undefined
                }]
            })
    })

    it('works with test-pod', () => {
        const data = {
            edges: [{
                from: "A", id: "myEdge", ports: [], to: "B",
                fromDeployment: "f-deployment",
                toDeployment: "t-deployment",
                label: "label",
                deniedConnection: false
            }],
            nodes: [{
                id: "test-pod",
                name: "pod1",
                label: "label1",
                deployment: "dep1",
                group: undefined
            }]
        }
        expect(getColoredData(data, { 'dep1': 'green' }))
            .toEqual({
                edges: [{
                    from: "A", id: "myEdge", ports: [], to: "B",
                    fromDeployment: "f-deployment",
                    toDeployment: "t-deployment",
                    label: "label",
                    deniedConnection: false
                }],
                nodes: [{
                    id: "test-pod",
                    name: "pod1",
                    label: "label1",
                    color: "red",
                    deployment: 'dep1',
                    group: undefined
                }]
            })
    })

    it('works with internet', () => {
        const data = {
            edges: [{
                from: "A", id: "myEdge", ports: [], to: "B",
                fromDeployment: "f-deployment",
                toDeployment: "t-deployment",
                label: "label",
                deniedConnection: false
            }],
            nodes: [{
                id: "Internet",
                name: "pod1",
                label: "label1",
                deployment: "dep1",
                group: undefined
            }]
        }
        expect(getColoredData(data, { 'dep1': 'green' }))
            .toEqual({
                edges: [{
                    from: "A", id: "myEdge", ports: [], to: "B",
                    fromDeployment: "f-deployment",
                    toDeployment: "t-deployment",
                    label: "label",
                    deniedConnection: false
                }],
                nodes: [{
                    id: "Internet",
                    name: "pod1",
                    label: "label1",
                    color: "#FFBF00",
                    deployment: 'dep1',
                    group: undefined
                }]
            })
    })
})

describe('hidePorts', () => {
    it('should set port to hidden', function () {
        const edges: SimpleGraphEdge[] = [
            SimpleGraphEdgeBuilder().port('80').build(),
            SimpleGraphEdgeBuilder().port('1234').build()];
        const expectedEdges: SimpleGraphEdge[] = [
            SimpleGraphEdgeBuilder().port('80').build(),
            SimpleGraphEdgeBuilder().port('1234').hidden(true).build()];
        expect(hidePorts(edges, ['1234'])).toEqual(expectedEdges);
    });
})
// When merging hidden edges, they are not displayed.

describe('mergeEdgesSimple', () => {
    it('should merge edges of the same type', function () {
        const edges: SimpleGraphEdge[] = [
            SimpleGraphEdgeBuilder().from('A').to('B').id('1').port('80').label('80').build(),
            SimpleGraphEdgeBuilder().from('A').to('B').id('2').port('70').label('70').build(),
            SimpleGraphEdgeBuilder().from('A').to('C').id('3').port('80').label('80').build(),
        ]
        const expectedOutput: GraphEdge[] = [
            GraphEdgeBuilder().from('A').to('B').ports(['80', '70']).id('1').label("80, 70").build(),
            GraphEdgeBuilder().from('A').to('C').ports(['80']).id('3').label('80').build()
        ]
        expect(mergeEdgesSimple(edges)).toEqual(expectedOutput)
    });
    it('should not return hidden edges', function () {
        const edges: SimpleGraphEdge[] = [
            SimpleGraphEdgeBuilder().from('A').to('B').id('1').port('80').label('80').build(),
            SimpleGraphEdgeBuilder().from('A').to('B').id('2').port('70').label('70').build(),
            SimpleGraphEdgeBuilder().from('A').to('B').id('3').hidden(true).port('8080').label('8080').build(),
        ]
        const expectedOutput: GraphEdge[] = [
            GraphEdgeBuilder().from('A').to('B').ports(['80', '70']).id('1').label("80, 70").build(),
        ]
        expect(mergeEdgesSimple(edges)).toEqual(expectedOutput)
    });
    it('should merge edges with different protocol', function () {
        const edges: SimpleGraphEdge[] = [
            SimpleGraphEdgeBuilder().from('A').to('B').id('1').port('53/UDP').label('53').build(),
            SimpleGraphEdgeBuilder().from('A').to('B').id('2').port('53/TCP').label('53').build(),
        ]
        const expectedOutput: GraphEdge[] = [
            GraphEdgeBuilder().from('A').to('B').ports(['53/UDP', '53/TCP']).id('1').label("53/UDP, 53/TCP").build(),
        ]
        expect(mergeEdgesSimple(edges)).toEqual(expectedOutput)
    });
    it('should return hidden edges belonging to denied rules', function () {
        const edges: SimpleGraphEdge[] = [
            SimpleGraphEdgeBuilder().from('A').to('B').id('1').port('80').label('80').build(),
            SimpleGraphEdgeBuilder().from('A').to('B').id('2').port('70').label('70').build(),
            SimpleGraphEdgeBuilder().from('A').to('B').id('3').hidden(true).deniedConnection(true).port('8080').label('8080').build(),
        ]
        const expectedOutput: GraphEdge[] = [
            GraphEdgeBuilder().from('A').to('B').ports(['80', '70']).id('1').label("80, 70").build(),
            GraphEdgeBuilder().from('A').to('B').ports(['8080']).id('3').label('8080').hidden(true).deniedConnection(true).build()

        ]
        expect(mergeEdgesSimple(edges)).toEqual(expectedOutput)
    });
})

describe('toCyNode', () => {
    it('converts deployment', () => {
        const node = GraphNodeBuilder().id('Custom').deployment('Custom').build()
        expect(toCyNode(node)).toEqual({
            data: {
                id: node.id,
                label: node.label,
                // @ts-ignore
                bg: node.color,
                type: "deployment",
                hidden: "false"
            }
        })
    })
    it('converts node', () => {
        const node = GraphNodeBuilder().deployment('Custom').build()
        expect(toCyNode(node)).toEqual({
            data: {
                id: node.id,
                label: node.label,
                // @ts-ignore
                bg: node.color,
                type: "pod",
                hidden: "false"
            }
        })
    })
})