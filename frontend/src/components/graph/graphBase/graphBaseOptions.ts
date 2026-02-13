export const cytoscapeStylesheet = [
    {
        selector: "node",
        style: {
            backgroundColor: "data(bg)",
            padding: "6px",
            color: "black",
            "border-width": 2,
            "border-color": 'black',
            "min-width": "50px",
            "min-height": "50px",
            "text-wrap": "wrap",
            "text-max-width": "40px",
            "text-events": "yes",
            "text-background-color": "transparent",
            "text-background-opacity": 0,
            "text-background-padding": "0px",
            "text-margin-x": "0px",
            "text-margin-y": "0px",
            "text-valign": "center",
            "text-halign": "center",
            "text-overflow-wrap": "break-word",
        }
    },
    {
        selector: 'node[type="deployment"]',
        style: {
            shape: 'rectangle',
        }
    },
    {
        selector: 'node[type="pod"]',
        style: {
            shape: 'ellipse',
        }
    },
    {
        selector: 'node[type="service"]',
        style: {
            shape: 'octagon',
        }
    },
    {
        selector: 'node[type="internet"]',
        style: {
            shape: "diamond"
        }
    },
    {
        selector: 'node[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        selector: "node[label]",
        style: {
            label: "data(label)",
            "font-size": "7",
            color: "black",
            "text-halign": "center",
            "text-valign": "center",
            "text-margin-x": "0px",
            "text-margin-y": "0px",
            "text-background-color": "transparent",
            "text-background-opacity": 0,
            "text-background-padding": "0px",
            "text-overflow-wrap": "break-word",
        }
    },
    {
        selector: "edge",
        style: {
            width: 1.5,
            "curve-style": "bezier",
            "target-arrow-shape": "triangle",
            "control-point-step-size": 100,
        }
    },
    {
        selector: "edge[label]",
        style: {
            "target-label": "data(label)",
            "font-size": "6",
            "text-background-color": "white",
            "text-background-opacity": 1,
            "text-background-padding": "1px",
            "text-border-color": "black",
            "text-border-style": "solid",
            "text-border-width": 0.5,
            "text-border-opacity": 1,
            "text-rotation": "autorotate",
            'target-text-offset': 35,
        },
    },
    {
        selector: 'edge[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        selector: 'edge[denied="true"]',
        style: {
            "background-color": "#FACD37",
            "text-outline-color": "#FACD37",
            "text-border-color": "red",
            "color": "red",
            "target-arrow-color": "red",
            "line-color": "red"


        }
    }
] as Array<cytoscape.Stylesheet>;

export const cytoscapeStylesheetPrintMode = [
    {
        selector: "node",
        style: {
            backgroundColor: "data(bg)",
            padding: "6px",
            "border-width": 2,
            "border-color": 'black',
            color: "black",
            "min-width": "50px",
            "min-height": "50px",
            "text-wrap": "wrap",
            "text-max-width": "40px",
            "text-events": "yes",
            "text-background-color": "transparent",
            "text-background-opacity": 0,
            "text-background-padding": "0px",
            "text-margin-x": "0px",
            "text-margin-y": "0px",
            "text-valign": "center",
            "text-halign": "center",
            "text-overflow-wrap": "break-word",
        }
    },
    {
        selector: 'node[type="deployment"]',
        style: {
            'shape': 'rectangle',
        }
    },
    {
        selector: 'node[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        selector: "node[label]",
        style: {
            label: "data(label)",
            "font-size": "7",
            color: "black",
            "text-halign": "center",
            "text-valign": "center",
            "text-margin-x": "0px",
            "text-margin-y": "0px",
            "text-background-color": "transparent",
            "text-background-opacity": 0,
            "text-background-padding": "0px",
            "text-overflow-wrap": "break-word",
        }
    },
    {
        selector: "edge",
        style: {
            width: 1.5,
            "curve-style": "bezier",
            "target-arrow-shape": "triangle",
            "control-point-step-size": 100,
        }
    },
    {
        selector: "edge[label]",
        style: {
            "target-label": "data(label)",
            "font-size": "6",
            "text-background-color": "white",
            "text-background-opacity": 1,
            "text-background-padding": "1px",
            "text-border-color": "black",
            "text-border-style": "solid",
            "text-border-width": 0.5,
            "text-border-opacity": 1,
            "text-rotation": "autorotate",
            'target-text-offset': 35,
        },
    },
    {
        selector: 'edge[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        selector: 'edge[denied="true"]',
        style: {
            "background-color": "#FACD37",
            "text-outline-color": "#FACD37",
            "text-border-color": "red",
            "color": "red",
            "target-arrow-color": "red",
            "line-color": "red"


        }
    }
] as Array<cytoscape.Stylesheet>;

